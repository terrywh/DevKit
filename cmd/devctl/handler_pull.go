package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"io/fs"
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

func (handler *HandlerPull) Do(ctx context.Context) (err error) {
	wd, _ := os.Getwd()
	sf := entity.ServerStreamFilePull{
		StreamFilePull: entity.StreamFilePull{
			StreamFile: entity.StreamFile{
				Path: wd,
			},
		},
		Pid: os.Getppid(),
	}

	var rsp *http.Response
	if rsp, err = handler.HTTPPost("/file/pull", sf); err != nil {
		return err
	}

	// file, _ := os.Create("./pull.rst")
	// defer file.Close()
	// log.Println(io.Copy(file, rsp.Body))
	// return

	x := entity.HttpResponse{Data: &sf.StreamFile}
	r := bufio.NewReader(rsp.Body)
	if err = app.ReadJSON(r, &x); err != nil {
		return err
	}
	var path string
	if path, err = handler.streaming(ctx, &sf.StreamFile, r); err != nil {
		return err
	}
	if !handler.override && handler.exists(sf.Path) {
		return entity.ErrFileExisted
	}
	err = os.Rename(path, sf.Path)
	return
}

func (handler *HandlerPull) streaming(ctx context.Context, sf *entity.StreamFile, src io.Reader) (path string, err error) {
	file, err := os.CreateTemp(filepath.Dir(sf.Path), filepath.Base(sf.Path)+".devkit_tmp_")
	path = file.Name()
	if err != nil {
		return
	}
	defer file.Close()

	log.Println("<HandlerPull.streaming> streaming file:", sf.Path, sf.Size, sf.Perm)
	pr := progressbar.DefaultBytes(
		sf.Size,
		"传输",
	)
	defer pr.Close()

	size, err := io.Copy(io.MultiWriter(file, pr, app.ContextDiscardWriter{ctx}), src)
	if err != nil {
		return
	}
	if size != sf.Size {
		err = entity.ErrFileCorrupted
		return
	}
	file.Close()
	err = os.Chmod(file.Name(), fs.FileMode(sf.Perm))
	return
}

func (handler *HandlerPull) exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (handler *HandlerPull) Close() error {
	return nil
}
