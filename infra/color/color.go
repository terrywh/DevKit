package color

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"golang.org/x/term"
)

var debug = color.New()
var info = color.New(color.FgCyan)
var warn = color.New(color.FgYellow)
var fatal = color.New(color.FgRed)

func Output(fp io.Writer, c *color.Color, args ...any) {
	if file, ok := fp.(*os.File); ok {
		if term.IsTerminal(int(file.Fd())) {
			c.Fprint(fp, args...)
			return
		}
	}
	fmt.Fprint(fp, args...)
}

func Debug(args ...any) {
	Output(os.Stdout, debug, args...)
}

func Info(args ...any) {
	Output(os.Stdout, info, args...)
}

func Warn(args ...any) {
	Output(os.Stdout, warn, args...)
}

func Error(args ...any) {
	Output(os.Stdout, fatal, args...)
}
