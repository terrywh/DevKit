package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
)

type HandlerPull struct {
	HandlerBase
	override bool
}

func (handler *HandlerPull) InitFlag(fs *flag.FlagSet) {
	fs.BoolVar(&handler.override, "o", false, "")
	fs.BoolVar(&handler.override, "override", false, "覆盖本地已存在的文件")
}

type BashQuery struct {
	BashPid int `json:"bash_pid"`
}

func (handler *HandlerPull) Do(ctx context.Context) (err error) {
	wd, _ := os.Getwd()
	sf := entity.StreamFile{
		// Target: entity.File{} // 获取到文件流，不指定目标
	}

	var rsp *http.Response
	if rsp, err = handler.HTTPPost(fmt.Sprintf("/file/pull?bash_id=%d", os.Getppid()), sf); err != nil {
		return err
	}

	// file, _ := os.Create("./pull.rst")
	// defer file.Close()
	// log.Println(io.Copy(file, rsp.Body))
	// return

	x := entity.Response{Error: &entity.DefaultErrorCode{}, Data: &sf}
	r := bufio.NewReader(rsp.Body)
	if err = app.ReadJSON(r, &x); err != nil {
		return err
	}
	prog := progressbar.DefaultBytes(sf.Source.Size)
	defer prog.Close()
	// 填写目标文件，从流接收写入
	sf.Target.Path = filepath.Join(wd, filepath.Base(sf.Source.Path))
	sf.Options.Override = handler.override

	log.Printf("<HandlePull.Do> StreamFile: %+v", sf)

	proc := &app.StreamFile{Desc: &sf, Prog: prog}
	err = proc.Do(ctx, r)
	return
}

func (handler *HandlerPull) Close() error {
	return nil
}
