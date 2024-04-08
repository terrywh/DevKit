package main

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
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
	"golang.org/x/term"
	"nhooyr.io/websocket"
)

type ServiceHttpHandlerShell struct {
	app.HttpHandlerBase
	mgr   stream.SessionManager
	mutex *sync.RWMutex
	shell map[entity.ShellID]*entity.RemoteShell
}

func newClientShellHandler(mgr stream.SessionManager, mux *http.ServeMux) *ServiceHttpHandlerShell {
	css := &ServiceHttpHandlerShell{
		mgr:   mgr,
		mutex: &sync.RWMutex{},
		shell: make(map[entity.ShellID]*entity.RemoteShell),
	}

	mux.HandleFunc("/shell/start", css.HandleStart)
	mux.HandleFunc("/shell/{shell_id}/socket", css.HandleSocket)
	mux.HandleFunc("/shell/{shell_id}/resize", css.HandleResize)
	mux.HandleFunc("/shell/run", css.HandleRun)

	// TODO cleanup
	return css
}

func (css *ServiceHttpHandlerShell) put(e *entity.RemoteShell) {
	css.mutex.Lock()
	defer css.mutex.Unlock()
	css.shell[e.ShellId] = e
}

func (css *ServiceHttpHandlerShell) get(shell_id entity.ShellID) *entity.RemoteShell {
	css.mutex.RLock()
	defer css.mutex.RUnlock()

	return css.shell[shell_id]
}

func (css *ServiceHttpHandlerShell) del(e *entity.RemoteShell) {
	css.mutex.Lock()
	defer css.mutex.Unlock()
	delete(css.shell, e.ShellId)
}

func (css *ServiceHttpHandlerShell) HandleStart(rsp http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	d := json.NewDecoder(req.Body)
	e := &entity.RemoteShell{}
	if err := d.Decode(e); err != nil {
		css.Respond(rsp, err)
		return
	}
	if shell_id, err := css.prepareShell(ctx, e); err != nil {
		css.Respond(rsp, err)
	} else {
		css.Respond(rsp, shell_id)
	}
}

func (css *ServiceHttpHandlerShell) HandleSocket(rsp http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	// 确认已注册的会话
	shell_id := entity.ShellID(req.PathValue("shell_id"))
	e := css.get(shell_id)
	if e == nil {
		log.Println("<ServiceHttpShell.HandleSocket> unable to find shell")
		return
	}
	defer css.del(e)
	// 确认对应会话通道
	ss, err := css.mgr.Acquire(ctx, &e.RemotePeer)
	if err != nil {
		log.Println("<ServiceHttpShell.HandleSocket> failed to acquire stream: ", err)
		return
	}
	defer ss.Close()
	// 确认对应前端通道
	c, err := websocket.Accept(rsp, req, &websocket.AcceptOptions{
		Subprotocols: []string{"shell"},
	})
	if err != nil {
		log.Println("<ServiceHttpShell.HandleSocket> failed to accept websocket: ", err)
		return
	}
	defer c.CloseNow()
	r, w := css.splitSocket(ctx, c)
	// 通道双向对转
	if err = css.serveShell(ctx, e, ss, r, w); err != nil {
		c.Close(websocket.StatusNormalClosure, err.Error())
	} else {
		c.Close(websocket.StatusNormalClosure, "")
	}
}

func (css *ServiceHttpHandlerShell) prepareShell(ctx context.Context, server *entity.RemoteShell) (shell_id entity.ShellID, err error) {
	// 确保能够联通（内部可能通过 REGISTRY 进行地址查询和反向发包）
	_, err = css.mgr.EnsureConn(ctx, &server.RemotePeer)
	if err != nil {
		return
	}
	server.ShellId = entity.ShellID(uuid.New().String())
	css.put(server)
	return server.ShellId, nil

}

func (css *ServiceHttpHandlerShell) splitSocket(ctx context.Context, c *websocket.Conn) (r io.Reader, w io.WriteCloser) {
	r = &WebsocketReader{ctx, c, &bytes.Buffer{}}
	w = &WebsocketWriter{ctx, c}
	return
}

func (css *ServiceHttpHandlerShell) serveShell(_ context.Context, e *entity.RemoteShell, ss *stream.SessionStream, r io.Reader, w io.Writer) (err error) {
	if r == os.Stdin { // 对直接透传的 Shell 设定当前 Stdin 状态
		state, _ := term.MakeRaw(int(os.Stdin.Fd()))
		e.Cols, e.Rows, _ = term.GetSize(int(os.Stdin.Fd()))
		defer term.Restore(int(os.Stdin.Fd()), state)
	}

	io.WriteString(ss, "/shell/start:")
	json.NewEncoder(ss).Encode(e)

	go io.Copy(ss, r)
	_, err = io.Copy(w, ss)
	return
}

func (css *ServiceHttpHandlerShell) HandleResize(rsp http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	shell_id := entity.ShellID(req.PathValue("shell_id"))
	e := css.get(shell_id)
	if e == nil {
		log.Println("<ServiceHttpShell.HandleSocket> unable to find shell")
		css.Respond(rsp, entity.ErrSessionNotFound)
		return
	}
	json.NewDecoder(req.Body).Decode(e)
	ss, err := css.mgr.Acquire(ctx, &e.RemotePeer)
	if err != nil {
		log.Println("<ServiceHttpShell.HandleSocket> failed acquire session: ", err)
		css.Respond(rsp, err)
		return
	}
	err = ss.Invoke(ctx, "/shell/resize", e, nil)
	css.Respond(rsp, err)
}

func (css *ServiceHttpHandlerShell) HandleRun(rsp http.ResponseWriter, req *http.Request) {

}
