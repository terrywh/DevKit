package util

import (
	"bytes"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/fatih/color"
	"golang.org/x/term"
)

var colorInfo = color.New(color.FgCyan)
var colorWarn = color.New(color.FgYellow)
var colorFail = color.New(color.FgRed)

var colorTitle = color.New(color.BgHiMagenta, color.FgHiWhite)

func Info(arg ...interface{}) {
	colorInfo.Fprint(os.Stderr, arg...)
	os.Stderr.Write([]byte{'\n'})
	os.Stderr.Sync()
}

func Warn(arg ...interface{}) {
	colorWarn.Fprint(os.Stderr, arg...)
	os.Stderr.Write([]byte{'\n'})
	os.Stderr.Sync()
}

func Fail(arg ...interface{}) {
	colorFail.Fprint(os.Stderr, arg...)
	os.Stderr.Write([]byte{'\n'})
	os.Stderr.Sync()
}

func Title(arg ...interface{}) {
	w, _, _ := term.GetSize(int(os.Stderr.Fd()))
	var buffer bytes.Buffer
	fmt.Fprint(&buffer, arg...)
	l := 0
	for _, rune := range buffer.String() {
		if utf8.RuneLen(rune) > 1 {
			l += 2
		} else {
			l += 1
		}
	}
	for i := 0; i < (w-l)/2; i++ {
		colorTitle.Fprint(os.Stderr, " ")
	}
	colorTitle.Fprint(os.Stderr, buffer.String())
	for i := 0; i < (w-l)/2; i++ {
		colorTitle.Fprint(os.Stderr, " ")
	}
	os.Stderr.Write([]byte{'\n'})
	os.Stderr.Sync()
}
