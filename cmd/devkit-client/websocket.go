package main

import (
	"context"
	"io"

	"nhooyr.io/websocket"
)

type WebsocketReader struct {
	ctx  context.Context
	conn *websocket.Conn
	r    *io.PipeReader
}

func NewWebSocketReader(ctx context.Context, conn *websocket.Conn) io.Reader {
	wsr := &WebsocketReader{ctx: ctx, conn: conn}
	var w *io.PipeWriter
	wsr.r, w = io.Pipe()
	go func(w *io.PipeWriter) {
		var err error
		var data []byte
		for {
			if _, data, err = conn.Read(ctx); err != nil {
				break
			}
			if _, err = w.Write(data); err != nil {
				break
			}
		}
		w.CloseWithError(err)
	}(w)
	return wsr
}

// Read io.Reader
func (wsr *WebsocketReader) Read(data []byte) (n int, err error) {
	return wsr.r.Read(data)
}

type WebsocketWriter struct {
	ctx  context.Context
	conn *websocket.Conn
}

func NewWebSocketWriter(ctx context.Context, conn *websocket.Conn) *WebsocketWriter {
	return &WebsocketWriter{ctx, conn}
}

func (wsr *WebsocketWriter) Write(data []byte) (n int, err error) {
	err = wsr.conn.Write(wsr.ctx, websocket.MessageBinary, data)
	if err == nil {
		n = len(data)
	}
	return
}

func (wsr *WebsocketWriter) Close() error {
	return wsr.conn.Close(websocket.StatusNormalClosure, "")
}
