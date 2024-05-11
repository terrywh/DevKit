package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
)

type HandlerPush struct {
	HandlerBase
	flagCommand *flag.FlagSet
}

func (handler *HandlerPush) InitFlag(flagCommand, flagGlobal *flag.FlagSet) {
	handler.flagCommand = flagCommand
	flagCommand.Usage = func() {
		fmt.Println(os.Args[0], "<全局选项>", flagCommand.Name(), "<命令选项> <目标文件>")
		fmt.Println("全局选项:")
		flagGlobal.PrintDefaults()
		fmt.Println("命令选项:")
		flagCommand.PrintDefaults()
	}
}

func (handler *HandlerPush) Do(ctx context.Context) (err error) {
	if handler.flagCommand.NArg() < 1 {
		err = entity.ErrInvalidArguments
		return
	}

	bashpid, err := GetBashPid()
	if err != nil {
		return err
	}
	err = handler.doFile(ctx, bashpid, handler.flagCommand.Arg(0))
	return
}

func (handler *HandlerPush) doFile(ctx context.Context, bashpid int, file string) (err error) {
	sf := entity.StreamFile{
		Source: entity.File{
			Path: file,
		},
	}
	info, err := os.Stat(sf.Source.Path)
	if err != nil {
		return err
	}
	sf.Source.Size = info.Size()
	sf.Source.Perm = uint32(info.Mode().Perm())

	body, err := os.Open(sf.Source.Path)
	if err != nil {
		return err
	}
	// HTTP POST 会自行关闭 Body 但未能将 file 作为 io.Closer 传递
	defer body.Close()

	log.DebugContext(ctx, "<devkit> stream file:", sf.Source.Path)

	prog := progressbar.DefaultBytes(sf.Source.Size, filepath.Base(sf.Source.Path))
	defer prog.Close()

	rbody := io.TeeReader(body, prog)
	path := fmt.Sprintf("/file/push?bash_pid=%d&path=%s&size=%d&perm=%d",
		bashpid,
		url.QueryEscape(sf.Source.Path),
		sf.Source.Size,
		sf.Source.Perm,
	)
	var rsp *http.Response
	if rsp, err = handler.Post(path, rbody); err != nil {
		return err
	}
	defer rsp.Body.Close()

	r := bufio.NewReader(rsp.Body)
	err = app.Read(r, &sf)
	return
}

func (handler *HandlerPush) Close() error {
	return nil
}
