package util

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
)

func FileExists(file string) (err error) {
	_, err = os.Stat(file)
	return
}

func ReadString(r io.Reader) (string, error) {
	data, err := ioutil.ReadAll(r)
	return string(data), err
}

func FileExecute(ctx context.Context, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	buf, err := cmd.Output()
	return string(buf), err
}

func FileExecuteJSON(ctx context.Context, dst interface{}, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	buf, err := cmd.Output()
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, dst)
}

func FileAddPerm(path string, mode fs.FileMode) (err error) {
	var stat fs.FileInfo
	if stat, err = os.Stat(path); err != nil {
		return err
	}

	return os.Chmod(path, stat.Mode()|mode)
}
