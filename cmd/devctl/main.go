package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/terrywh/devkit/app"
)

func main() {
	flagGlobal := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagGlobal.Usage = func() {
		fmt.Println(flagGlobal.Name(), "<global-options> <command> <command-options>")
		fmt.Println("global-options:")
		flagGlobal.Usage()
		fmt.Println("command:")
		fmt.Println("  help\t查询各命令帮助")
		fmt.Println("  pull\t从操作客户端拉取文件")
	}

	DefaultConfig.Init(filepath.Join(app.GetBaseDir(), "etc", "devkit.yaml"))
	DefaultConfig.Get().InitFlag(flagGlobal)
	flagGlobal.Parse(os.Args[1:])

	gcmd := flagGlobal.Arg(0)
	scmd := flagGlobal.Arg(0)

	var handler Handler
	if gcmd == "help" {
		scmd = flagGlobal.Arg(1)
	}
	flagCommand := flag.NewFlagSet(scmd, flag.ExitOnError)
	flagCommand.Usage = func() {
		fmt.Println(os.Args[0], "<global-options> ", flagCommand.Name(), " <command-options>")
		fmt.Println("command-options:")
		flagCommand.PrintDefaults()
	}
	switch scmd {
	case "pull":
		handler = &HandlerPull{HandlerBase: HandlerBase{DefaultConfig.Get().Server.Address}}
		handler.InitFlag(flagCommand)
	default:
		fmt.Println("error: unknown command")
		flagGlobal.Usage()
		os.Exit(2)
		return
	}

	if gcmd == "help" {
		fmt.Println(os.Args[0], "<global-options>", scmd, "<command-options>")
		flagCommand.Usage()
		return
	}
	flagCommand.Parse(flagGlobal.Args()[1:])

	sc := app.NewServiceController()
	sc.Start(&HandlerService{scmd, handler})
	sc.WaitForSignal()
	sc.Close()
}
