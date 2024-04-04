package handler

import (
	"bytes"
	"context"

	"nhooyr.io/websocket"
)

type WebSocketReader struct {
	ctx    context.Context
	conn   *websocket.Conn
	buffer *bytes.Buffer
}

// Read io.Reader
func (wsr *WebSocketReader) Read(data []byte) (n int, err error) {
	if wsr.buffer.Len() > 0 {
		return wsr.buffer.Read(data)
	}
	var payload []byte
	if _, payload, err = wsr.conn.Read(wsr.ctx); err != nil {
		return
	}
	n, err = wsr.buffer.Write(payload)
	if err != nil {
		return
	}
	return wsr.buffer.Read(data)
}

type WebSocketWriteCloser struct {
	ctx  context.Context
	conn *websocket.Conn
}

func (wsr *WebSocketWriteCloser) Write(data []byte) (n int, err error) {
	err = wsr.conn.Write(wsr.ctx, websocket.MessageBinary, data)
	if err == nil {
		n = len(data)
	}
	return
}

func (wsr *WebSocketWriteCloser) Close() error {
	return wsr.conn.Close(websocket.StatusNormalClosure, "")
}
