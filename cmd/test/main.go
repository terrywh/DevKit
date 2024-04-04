package main

import (
	"context"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	RunWebview(context.Background(), wg)
	wg.Wait()
}
