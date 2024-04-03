package main

import (
	"context"
	"flag"
	"sync"

	"github.com/terrywh/devkit/cmd/v2/client"
	"github.com/terrywh/devkit/cmd/v2/server"
)

func main() {
	svr := server.New()
	cli := client.New()

	svr.ParseFlags()
	cli.ParseFlags()
	flag.Parse()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go svr.Serve(context.Background(), wg)
	go cli.Serve(context.Background(), wg)
	wg.Wait()
}

// func main() {
// 	wg := &sync.WaitGroup{}
// 	wg.Add(1)
// 	RunWebview(context.Background(), wg)
// 	wg.Wait()
// }
