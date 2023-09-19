package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v3"
)

var ConstBaseDir string
var ConstBinaryDir string
var ConstConfigDir string
var ConstCachesDir string

type Config interface {
	Get(string, string) string
}

var DefaultConfig Config

func OnInit(conf interface{}) {
	if ptr, ok := conf.(Config); !ok || reflect.ValueOf(conf).Kind() != reflect.Pointer {
		log.Fatal("failed to initialize config, must be Pointer type that satisfy: 'Config'")
		return
	} else {
		DefaultConfig = ptr
	}

	bin, _ := os.Executable()
	ConstBaseDir, _ = filepath.Abs(filepath.Dir(bin))
	ConstBinaryDir = filepath.Join(ConstBaseDir, "bin")
	ConstConfigDir = filepath.Join(ConstBaseDir, "etc")
	ConstCachesDir = filepath.Join(ConstBaseDir, "var")
	os.MkdirAll(ConstBinaryDir, 0o766)
	os.MkdirAll(ConstConfigDir, 0o700)
	os.MkdirAll(ConstCachesDir, 0o700)

	path := filepath.Join(ConstConfigDir, fmt.Sprintf("%s.yaml", filepath.Base(bin)))
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		d := yaml.NewDecoder(file)
		err = d.Decode(DefaultConfig)
	}
	if err != nil {
		log.Println("<warning> failed to initialize config: ", path, ", due to: ", err)
		return
	}
}
