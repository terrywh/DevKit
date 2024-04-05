package application

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Service interface {
	Serve(ctx context.Context)
	Close() error
}

type ServiceController struct {
	wg      *sync.WaitGroup
	running []Service
}

func NewServiceController() (sc *ServiceController) {
	sc = &ServiceController{wg: &sync.WaitGroup{}}
	return
}

func (sc *ServiceController) Start(svc Service) {
	sc.wg.Add(1)
	sc.running = append(sc.running, svc)

	go func() {
		defer sc.wg.Done()
		svc.Serve(context.Background())
	}()
}

func (sc *ServiceController) Shutdown() {
	for _, svc := range sc.running {
		svc.Close()
	}
	sc.wg.Wait()
}

func (sc *ServiceController) WaitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // 等待上述信号
}
