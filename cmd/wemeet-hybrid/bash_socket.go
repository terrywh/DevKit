package main

import "golang.org/x/net/websocket"

type BashSocket struct {
	conn *websocket.Conn
}

func (bs BashSocket) Read(data []byte) (int, error) {
	return bs.conn.Read(data)
}

func (bs BashSocket) Write(data []byte) (n int, err error) {
	n = len(data)
	err = websocket.Message.Send(bs.conn, data)
	return
}

func (bs BashSocket) Close() error {
	return bs.conn.Close()
}