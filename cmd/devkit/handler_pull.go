package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	var rsp *http.Response
	if rsp, err = handler.Post(fmt.Sprintf("/file/pull?bash_pid=%d", os.Getppid()), nil); err != nil {
		return err
	}
	defer rsp.Body.Close()

	wd, _ := os.Getwd()
	sf := entity.StreamFile{
		// Target: entity.File{} // 获取到文件流，不指定目标
	}

	r := bufio.NewReader(rsp.Body)
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
