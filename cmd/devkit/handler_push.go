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

	"github.com/schollz/progressbar/v3"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
)

type HandlerPush struct {
	HandlerBase
	file     string
	override bool
}

func (handler *HandlerPush) InitFlag(fs *flag.FlagSet) {
	fs.StringVar(&handler.file, "f", "", "")
	fs.StringVar(&handler.file, "file", "", "待下载推送的文件")
}

func (handler *HandlerPush) Do(ctx context.Context) (err error) {
	bashpid, err := GetBashPid()
	if err != nil {
		return err
	}

	sf := entity.StreamFile{
		Source: entity.File{
			Path: handler.file,
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

	log.InfoContext(ctx, "<devkit> stream file:", sf.Source.Path)

	prog := progressbar.DefaultBytes(sf.Source.Size)
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

	// file, _ := os.Create("./pull.rst")
	// defer file.Close()
	// app.Debug(io.Copy(file, rsp.Body))
	// return
	r := bufio.NewReader(rsp.Body)
	err = app.Read(r, &sf)
	return
}

func (handler *HandlerPush) Close() error {
	return nil
}
