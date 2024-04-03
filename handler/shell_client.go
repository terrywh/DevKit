package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
	"golang.org/x/term"
	"nhooyr.io/websocket"
)

type ClientShellHandler struct {
	HttpHandlerBase
	start map[entity.ShellID]*entity.StartShell
}

func NewClientShellHandler(mux *http.ServeMux) *ClientShellHandler {
	css := &ClientShellHandler{}
	css.start = make(map[entity.ShellID]*entity.StartShell)

	mux.HandleFunc("/shell/start", css.HandleStart)
	mux.HandleFunc("/shell/{shell_id}/socket", css.HandleSocket)
	mux.HandleFunc("/shell/{shell_id}/resize", css.HandleResize)
	mux.HandleFunc("/shell/run", css.HandleRun)

	// TODO cleanup
	return css
}

func (css *ClientShellHandler) HandleStart(rsp http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d := json.NewDecoder(req.Body)
	e := &entity.StartShell{}
	if err := d.Decode(e); err != nil {
		css.Failure(rsp, err)
		return
	}
	if shell_id, err := css.prepareShell(ctx, e); err != nil {
		css.Failure(rsp, err)
	} else {
		css.Success(rsp, shell_id)
	}
}

func (css *ClientShellHandler) HandleSocket(rsp http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	c, err := websocket.Accept(rsp, req, &websocket.AcceptOptions{
		Subprotocols: []string{"shell"},
	})
	if err != nil {
		log.Println("<ServiceHttpShell.handleSocket> failed to accept websocket: ", err)
		return
	}
	defer c.CloseNow()

	shell_id := entity.ShellID(req.PathValue("shell_id"))
	defer func() {
		delete(css.start, shell_id)
	}()

	e, ok := css.start[shell_id]
	if !ok {
		log.Println("<ServiceHttpShell.handleSocket> failed to find shell: ", shell_id)
		return
	}

	s, err := stream.DefaultSessionManager.AcquireStream(ctx, e.DeviceId)
	if err != nil {
		log.Println("<ServiceHttpShell.handleSocket> failed to acquire stream: ", err)
		return
	}

	r, w, err := css.splitSocket(ctx, c)
	if err != nil {
		log.Println("<ServiceHttpShell.handleSocket> failed to split socket: ", err)
		return
	}
	css.serveShell(ctx, e, s, r, w)
	c.Close(websocket.StatusNormalClosure, "")
}

func (css *ClientShellHandler) prepareShell(ctx context.Context, e *entity.StartShell) (shell_id entity.ShellID, err error) {
	// 确保能够联通
	_, err = stream.DefaultSessionManager.Acquire(ctx, e.DeviceId)
	if err != nil {
		return
	}
	shell_id = entity.ShellID(uuid.New().String())
	css.start[shell_id] = e
	return
}

func (css *ClientShellHandler) splitSocket(ctx context.Context, c *websocket.Conn) (r io.Reader, w io.WriteCloser, err error) {
	var typ websocket.MessageType
	typ, r, err = c.Reader(ctx)
	if err != nil {
		return
	}
	w, err = c.Writer(ctx, typ)
	if err != nil {
		return
	}
	return
}

func (css *ClientShellHandler) serveShell(_ context.Context, e *entity.StartShell, s quic.Stream, r io.Reader, w io.Writer) {
	e.ApplyDefaults()

	if r == os.Stdin { // 对直接透传的 Shell 设定当前 Stdin 状态
		state, _ := term.MakeRaw(int(os.Stdin.Fd()))
		e.Cols, e.Rows, _ = term.GetSize(int(os.Stdin.Fd()))
		defer term.Restore(int(os.Stdin.Fd()), state)
	}

	io.WriteString(s, "StartShell:")
	json.NewEncoder(s).Encode(e)

	go io.Copy(s, r)
	io.Copy(w, s)
}

func (css *ClientShellHandler) HandleResize(rsp http.ResponseWriter, req *http.Request) {

}

func (css *ClientShellHandler) HandleRun(rsp http.ResponseWriter, req *http.Request) {

}
