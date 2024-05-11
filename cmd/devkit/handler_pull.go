package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
)

type HandlerPull struct {
	HandlerBase
	override bool
}

func (handler *HandlerPull) InitFlag(flagCommand, flagGlobal *flag.FlagSet) {
	flagCommand.BoolVar(&handler.override, "o", false, "")
	flagCommand.BoolVar(&handler.override, "override", false, "覆盖本地已存在的文件")
}

func (handler *HandlerPull) Do(ctx context.Context) (err error) {
	var files []entity.SelectFile
	files, err = handler.doListFile(ctx)
	if err != nil {
		return
	}
	for _, file := range files {
		if err = handler.doPullFile(ctx, file); err != nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	return
}

func (handler *HandlerPull) doListFile(ctx context.Context) (sf []entity.SelectFile, err error) {
	var rsp *http.Response
	if rsp, err = handler.Post(fmt.Sprintf("/file/list?bash_pid=%d", os.Getppid()), nil); err != nil {
		return nil, err
	}
	r := bufio.NewReader(rsp.Body)
	err = app.Read(r, &sf)
	return
}

func (handler *HandlerPull) doPullFile(ctx context.Context, file entity.SelectFile) (err error) {
	var rsp *http.Response
	if rsp, err = handler.Post(fmt.Sprintf("/file/pull?bash_pid=%d", os.Getppid()), file); err != nil {
		return err
	}
	defer rsp.Body.Close()

	wd, _ := os.Getwd()
	r := bufio.NewReader(rsp.Body)
	sf := entity.StreamFile{}
	if err = app.Read(r, &sf); err != nil {
		return err
	}
	log.DebugContext(ctx, "<devkit> stream file:", sf.Source.Path)

	// 填写目标文件，从流接收写入
	sf.Target.Path = filepath.Join(wd, filepath.Base(sf.Source.Path))
	sf.Options.Override = handler.override

	proc := app.NewStreamFile(&sf, true)
	err = proc.Do(ctx, r)
	return
}

func (handler *HandlerPull) Close() error {
	return nil
}
