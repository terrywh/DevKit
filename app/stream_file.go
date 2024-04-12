package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra"
)

type StreamFile struct {
	Desc *entity.StreamFile

	Prog io.Writer
}

type cancelableWriter struct {
	ctx context.Context
}

func (w cancelableWriter) Write(payload []byte) (size int, err error) {
	size = len(payload)
	err = w.ctx.Err()
	return
}

type progressWriter struct{}

func (w progressWriter) Write(payload []byte) (size int, err error) {
	size = len(payload)
	return
}

func (s *StreamFile) Do(ctx context.Context, src io.Reader) (err error) {
	var target string
	if infra.IsDirectory(s.Desc.Target.Path) {
		target = filepath.Join(s.Desc.Target.Path, filepath.Base(s.Desc.Source.Path))
	} else {
		target = s.Desc.Target.Path
	}
	if !s.Desc.Options.Override && infra.Exists(target) {
		err = entity.ErrFileExisted
		return
	}

	var dst *os.File
	if dst, err = os.CreateTemp(filepath.Dir(s.Desc.Target.Path), filepath.Base(s.Desc.Target.Path)+".devkit_tmp_"); err != nil {
		err = fmt.Errorf("stream file (temp): %w", err)
		return
	}
	defer dst.Close()

	if s.Prog == nil {
		s.Prog = progressWriter{}
	}

	size, err := io.Copy(io.MultiWriter(dst, s.Prog, &cancelableWriter{ctx}), src)
	if err != nil {
		err = fmt.Errorf("stream file (copy): %w", err)
		return
	}
	if size != s.Desc.Source.Size {
		err = fmt.Errorf("stream file (size): %w", entity.ErrFileCorrupted)
		return
	}
	if err = os.Chmod(dst.Name(), os.FileMode(s.Desc.Source.Perm)); err != nil {
		infra.Warn("<app> failed to stream file (perm): ", err)
		err = nil
	}

	if err = os.Rename(dst.Name(), target); err != nil {
		err = fmt.Errorf("stream file (rename): %w", err)
	}
	return
}
