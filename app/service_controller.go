package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/terrywh/devkit/infra/log"
)

type Service interface {
	Serve(ctx context.Context)
	Close() error
}
type ServiceWithName interface {
	Service
	Name() string
}

type ServiceRunning struct {
	ctx    context.Context
	cancel context.CancelCauseFunc
	svc    Service
}

type ServiceController struct {
	wg      *sync.WaitGroup
	running []*ServiceRunning
}

func NewServiceController() (sc *ServiceController) {
	sc = &ServiceController{wg: &sync.WaitGroup{}}
	return
}

func (sc *ServiceController) Start(svc Service) {
	sc.wg.Add(1)
	sr := &ServiceRunning{
		svc: svc,
	}
	sr.ctx, sr.cancel = context.WithCancelCause(context.Background())
	if sn, ok := svc.(ServiceWithName); ok {
		sr.ctx = log.WithContextFields(sr.ctx, fmt.Sprintf("name=%s", sn.Name()))
	}
	sc.running = append(sc.running, sr)

	go func(sr *ServiceRunning) {
		defer sc.wg.Done()
		defer sr.cancel(nil)
		log.DebugContext(sr.ctx, "service started ...")
		svc.Serve(sr.ctx)
		log.DebugContext(sr.ctx, "service closed.")
	}(sr)
}

var ErrShutdown = errors.New("shutdown")

func (sc *ServiceController) Close() error {
	log.Trace("<app> service controller shutdown ...")
	for i := len(sc.running) - 1; i >= 0; i-- {
		sr := sc.running[i]
		sr.cancel(ErrShutdown)
	}
	// 自主结束超时后无人读取，估使用 BUFFERED CHANNEL 允许忽略
	done := make(chan bool, 1)
	go func() {
		timeout := time.NewTimer(10 * time.Second)
		timeout.Stop()
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

func (sc *ServiceController) Wait() {
	sc.wg.Wait()
}

func (sc *ServiceController) WaitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // 等待上述信号
}
