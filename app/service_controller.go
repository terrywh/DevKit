package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Service interface {
	Serve(ctx context.Context)
	Close() error
}

type ServiceRunning struct {
	ctx    context.Context
	cancel context.CancelCauseFunc
	svc    Service
}

type ServiceController struct {
	wg      *sync.WaitGroup
	running []ServiceRunning
}

func NewServiceController() (sc *ServiceController) {
	sc = &ServiceController{wg: &sync.WaitGroup{}}
	return
}

func (sc *ServiceController) Start(svc Service) {
	sc.wg.Add(1)
	sr := ServiceRunning{
		svc: svc,
	}
	sr.ctx, sr.cancel = context.WithCancelCause(context.Background())
	sc.running = append(sc.running, sr)

	go func() {
		defer sc.wg.Done()
		defer sr.cancel(nil)
		svc.Serve(context.Background())
	}()
}

var ErrShutdown = errors.New("shutdown")

func (sc *ServiceController) Close() error {
	for i := len(sc.running) - 1; i >= 0; i-- {
		sr := sc.running[i]
		sr.cancel(ErrShutdown)
	}
	done := make(chan bool)
	go func() {
		timeout := time.NewTimer(10 * time.Second)
		select { // 等待服务自然停止或超时
		case <-done:
		case <-timeout.C:
		}
		for i := len(sc.running) - 1; i >= 0; i-- {
			sr := sc.running[i]
			sr.svc.Close()
		}
	}()
	sc.wg.Wait()
	done <- true
	return nil
}

func (sc *ServiceController) WaitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // 等待上述信号
}
