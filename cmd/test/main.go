package main

import (
	"context"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go RunServer(context.Background(), wg)
	go RunClient(context.Background(), wg)
	wg.Wait()
}
