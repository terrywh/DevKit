package app

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/terrywh/devkit/infra/log"
	"gopkg.in/yaml.v3"
)

var base string = initBase()

func initBase() string {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") || arg == "test" {
			_, filename, _, _ := runtime.Caller(0)
			return filepath.Dir(filepath.Dir(filename))
		}
	}
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

type ConfigWithFlag interface {
	InitFlag()
}

type Config[T any] struct {
	path    string
	payload atomic.Pointer[T]
}

func (c *Config[T]) Init(path string) {
	c.path = path
	var v interface{}
	v = new(T)
	c.payload.Store(v.(*T))

	c.Reload()           // 需要一个前置存储的数据基础进行复制、覆盖
	v = c.payload.Load() // 从配置文件加载的配置
	if cf, ok := v.(ConfigWithFlag); ok {
		cf.InitFlag()
	}
}

func (c *Config[T]) Path() string {
	return c.path
}

func (c *Config[T]) OnChange() {
	c.Reload()
}

func (c *Config[T]) Get() *T {
	return c.payload.Load()
}

func (c *Config[T]) Reload() {
	log.Trace("<app> config reload: ", c.path)
	cp := *c.payload.Load()      // 保留当前值
	UnmarshalConfig(c.path, &cp) // 覆盖
	c.payload.Store(&cp)
}
