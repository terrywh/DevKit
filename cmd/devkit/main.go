package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/infra/log"
)

func main() {
	flagGlobal := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagGlobal.Usage = func() {
		fmt.Println(flagGlobal.Name(), "<全局选项> <命令> <命令选项>")
		fmt.Println("全局选项:")
		flagGlobal.PrintDefaults()
		fmt.Println("可用命令:")
		fmt.Println("  help\t查询各命令帮助")
		fmt.Println("  pull\t从操作客户端拉取文件")
		fmt.Println("  push\t向操作客户端推送文件")
	}

	DefaultConfig.Init(filepath.Join(app.GetBaseDir(), "etc", "devkit.yaml"))
	DefaultConfig.Get().InitFlag(flagGlobal)
	flagGlobal.Parse(os.Args[1:])

	log.DefaultLogger.SetLevel(log.LevelFromString(DefaultConfig.Get().Log.Level))

	gcmd := flagGlobal.Arg(0)
	scmd := flagGlobal.Arg(0)

	var handler Handler
	if gcmd == "help" {
		scmd = flagGlobal.Arg(1)
	}
	flagCommand := flag.NewFlagSet(scmd, flag.ExitOnError)
	flagCommand.Usage = func() {
		fmt.Println(os.Args[0], "<全局选项>", flagCommand.Name(), "<命令选项>")
		fmt.Println("全局选项:")
		flagGlobal.PrintDefaults()
		fmt.Println("命令选项:")
		flagCommand.PrintDefaults()
	}
	switch scmd {
	case "pull":
		handler = &HandlerPull{}
		handler.InitFlag(flagCommand, flagGlobal)
	case "push":
		handler = &HandlerPush{}
		handler.InitFlag(flagCommand, flagGlobal)
	default:
		fmt.Println("错误: 未知命令")
		flagGlobal.Usage()
		os.Exit(2)
		return
	}

	if gcmd == "help" {
		flagCommand.Usage()
		return
	}
	flagCommand.Parse(flagGlobal.Args()[1:])

	sc := app.NewServiceController()
	sc.Start(&HandlerService{scmd, handler})
	sc.Wait()
	sc.Close()
}
