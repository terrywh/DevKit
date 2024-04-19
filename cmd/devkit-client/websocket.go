package main

import (
	"bufio"
	"context"
	"io"
	"time"

	"nhooyr.io/websocket"
)

type Websocket struct {
	ctx  context.Context
	conn *websocket.Conn

	rr *io.PipeReader
	rw *io.PipeWriter
	rb *bufio.Writer
	rc chan []byte
}

func NewWebsocket(ctx context.Context, conn *websocket.Conn) *Websocket {
	wsr := &Websocket{ctx: ctx, conn: conn}
	wsr.rr, wsr.rw = io.Pipe()
	wsr.rb = bufio.NewWriterSize(wsr.rw, 64)
	wsr.rc = make(chan []byte)

	ctx, cancel := context.WithCancel(ctx)
	go wsr.write(ctx)
	go wsr.read(ctx, cancel)
	return wsr
}

func (wsr *Websocket) write(ctx context.Context) {
	ticker := time.NewTicker(160 * time.Millisecond)
	defer ticker.Stop()
	var err error
	var data []byte
SERVING:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case data = <-wsr.rc:
			_, err = wsr.rb.Write(data)
		case <-ticker.C:
			err = wsr.rb.Flush()
		}
		if err != nil {
			break SERVING
		}
	}
	// log.Debug("websocket flush error: ", err)
	wsr.rw.Close()
}

func (wsr *Websocket) read(ctx context.Context, cancel context.CancelFunc) {
	var err error
	var data []byte
	for {
		if _, data, err = wsr.conn.Read(ctx); err != nil {
			break
		}
		wsr.rc <- data
	}
	// log.Debug("<devkit-client> websocket read error:", err)
	cancel() // 通知 flush 停止
}

// Read io.Reader
func (wsr *Websocket) Read(data []byte) (n int, err error) {
	return wsr.rr.Read(data)
}

func (wsr *Websocket) Write(data []byte) (n int, err error) {
	n = len(data)
	err = wsr.conn.Write(wsr.ctx, websocket.MessageBinary, data)
	return
}

func (wsr *Websocket) CloseWrite() error {
	return wsr.conn.Close(websocket.StatusNormalClosure, "")
}
