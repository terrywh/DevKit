package app

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var base string = initBase()

func initBase() string {
	bin, _ := os.Executable()
	base, _ := filepath.Abs(bin)
	for i := 0; i < 5; i++ {
		base = filepath.Dir(base)
		if _, err := os.Stat(filepath.Join(base, "var")); os.IsNotExist(err) {
			continue
		}
		break
	}
	return base
}

func GetBaseDir() string {
	return base
}

var ErrUnsupportedFileType error = errors.New("unsupported file type")

func UnmarshalConfig(path string, v interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	switch filepath.Ext(path) {
	case ".yaml":
		err = yaml.NewDecoder(file).Decode(v)
	default:
		err = ErrUnsupportedFileType
	}
	return
}
