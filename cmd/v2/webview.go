package main

import (
	"context"
	"sync"

	webview "github.com/webview/webview_go"
)

func RunWebview(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Basic Example")
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml("Thanks for using webview!")
	w.Run()
}
