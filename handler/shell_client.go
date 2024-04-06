package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
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
	mutex *sync.RWMutex
	shell map[entity.ShellID]*entity.StartShell
}

func NewClientShellHandler(mux *http.ServeMux) *ClientShellHandler {
	css := &ClientShellHandler{
		mutex: &sync.RWMutex{},
		shell: make(map[entity.ShellID]*entity.StartShell),
	}

	mux.HandleFunc("/shell/start", css.HandleStart)
	mux.HandleFunc("/shell/{shell_id}/socket", css.HandleSocket)
	mux.HandleFunc("/shell/{shell_id}/resize", css.HandleResize)
	mux.HandleFunc("/shell/run", css.HandleRun)

	// TODO cleanup
	return css
}

func (css *ClientShellHandler) put(e *entity.StartShell) {
	css.mutex.Lock()
	defer css.mutex.Unlock()
	css.shell[e.ShellId] = e
}

func (css *ClientShellHandler) get(shell_id entity.ShellID) *entity.StartShell {
	css.mutex.RLock()
	defer css.mutex.RUnlock()

	return css.shell[shell_id]
}

func (css *ClientShellHandler) del(e *entity.StartShell) {
	css.mutex.Lock()
	defer css.mutex.Unlock()
	delete(css.shell, e.ShellId)
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

	shell_id := entity.ShellID(req.PathValue("shell_id"))
	defer func() {
		delete(css.shell, shell_id)
	}()

	e := css.get(shell_id)
	if e == nil {
		log.Println("<ServiceHttpShell.HandleSocket> failed to find shell: ", shell_id)
		return
	}
	defer css.del(e)

	s, err := stream.DefaultSessionManager.AcquireStream(ctx, e.DeviceId)
	if err != nil {
		log.Println("<ServiceHttpShell.HandleSocket> failed to acquire stream: ", err)
		return
	}

	c, err := websocket.Accept(rsp, req, &websocket.AcceptOptions{
		Subprotocols: []string{"shell"},
	})
	if err != nil {
		log.Println("<ServiceHttpShell.HandleSocket> failed to accept websocket: ", err)
		return
	}
	defer c.CloseNow()
	r, w := css.splitSocket(ctx, c)
	if err = css.serveShell(ctx, e, s, r, w); err != nil {
		c.Close(websocket.StatusNormalClosure, err.Error())
	} else {
		c.Close(websocket.StatusNormalClosure, "")
	}
}

func (css *ClientShellHandler) prepareShell(ctx context.Context, e *entity.StartShell) (shell_id entity.ShellID, err error) {
	// 确保能够联通
	_, err = stream.DefaultSessionManager.Acquire(ctx, e.DeviceId)
	if err != nil {
		return
	}
	e.ShellId = entity.ShellID(uuid.New().String())
	css.put(e)
	return e.ShellId, nil
}

func (css *ClientShellHandler) splitSocket(ctx context.Context, c *websocket.Conn) (r io.Reader, w io.WriteCloser) {
	r = &WebSocketReader{ctx, c, &bytes.Buffer{}}
	w = &WebSocketWriteCloser{ctx, c}
	return
}

func (css *ClientShellHandler) serveShell(_ context.Context, e *entity.StartShell, s quic.Stream, r io.Reader, w io.Writer) (err error) {
	if r == os.Stdin { // 对直接透传的 Shell 设定当前 Stdin 状态
		state, _ := term.MakeRaw(int(os.Stdin.Fd()))
		e.Cols, e.Rows, _ = term.GetSize(int(os.Stdin.Fd()))
		defer term.Restore(int(os.Stdin.Fd()), state)
	}

	io.WriteString(s, "/shell/start:")
	json.NewEncoder(s).Encode(e)

	go io.Copy(s, r)
	_, err = io.Copy(w, s)
	return
}

func (css *ClientShellHandler) HandleResize(rsp http.ResponseWriter, req *http.Request) {

}

func (css *ClientShellHandler) HandleRun(rsp http.ResponseWriter, req *http.Request) {

}
