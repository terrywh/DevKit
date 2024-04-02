package client

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/quic-go/quic-go"
	"golang.org/x/term"
)

func OpenShell(s quic.Stream, shell string) {
	w, h, _ := term.GetSize(int(os.Stdin.Fd()))
	fmt.Fprintf(s, `shell:{"shell":"%s","width":%d,"height":%d}%s`, shell, w, h, "\n")
	log.Println("accepting input ...")
	state, _ := term.MakeRaw(int(os.Stdin.Fd()))

	go io.Copy(s, os.Stdin)
	io.Copy(os.Stdout, s)

	term.Restore(int(os.Stdin.Fd()), state)
	log.Println("close.")
}
