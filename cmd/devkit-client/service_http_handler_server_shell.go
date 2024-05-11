package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
	"github.com/terrywh/devkit/stream"
	"nhooyr.io/websocket"
)

type ServerShellHandler struct {
	app.HttpHandlerBase
	mgr   stream.SessionManager
	mutex *sync.RWMutex
	shell map[entity.ShellID]*entity.ServerShell
}

func initHttpServerShellHandler(mgr stream.SessionManager, mux *http.ServeMux) *ServerShellHandler {
	ssh := &ServerShellHandler{
		mgr:   mgr,
		mutex: &sync.RWMutex{},
		shell: make(map[entity.ShellID]*entity.ServerShell),
	}

	mux.HandleFunc("/server/shell/prepare", ssh.HandlePrepare)
	mux.HandleFunc("/server/shell/{shell_id}/socket", ssh.HandleSocket)
	mux.HandleFunc("/server/shell/{shell_id}/resize", ssh.HandleResize)
	mux.HandleFunc("/server/shell/run", ssh.HandleRun)

	// TODO cleanup
	return ssh
}

func (ssh *ServerShellHandler) put(e *entity.ServerShell) {
	ssh.mutex.Lock()
	defer ssh.mutex.Unlock()
	ssh.shell[e.Shell.ID] = e
}

func (ssh *ServerShellHandler) get(shell_id entity.ShellID) *entity.ServerShell {
	ssh.mutex.RLock()
	defer ssh.mutex.RUnlock()

	return ssh.shell[shell_id]
}

func (ssh *ServerShellHandler) del(e *entity.ServerShell) {
	ssh.mutex.Lock()
	defer ssh.mutex.Unlock()
	delete(ssh.shell, e.Shell.ID)
}

func (ssh *ServerShellHandler) HandlePrepare(rsp http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	d := json.NewDecoder(req.Body)
	e := entity.ServerShell{}
	if err := d.Decode(&e); err != nil {
		ssh.Respond(rsp, err)
		return
	}
	if err := ssh.prepareShell(ctx, &e); err != nil {
		ssh.Respond(rsp, err)
	} else {
		ssh.Respond(rsp, e)
	}
}

func (ssh *ServerShellHandler) HandleSocket(rsp http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	// 确认已注册的会话
	shell_id := entity.ShellID(req.PathValue("shell_id"))
	e := ssh.get(shell_id)
	if e == nil {
		log.Warn("<devkit-client> unable to find shell:", shell_id)
		return
	}
	defer ssh.del(e)
	// 确认对应会话通道
	dst, err := ssh.mgr.Acquire(ctx, &e.Server)
	if err != nil {
		log.Warn("<devkit-client> failed to acquire stream:", err)
		return
	}
	defer dst.CloseWrite()
	// 确认对应前端通道
	c, err := websocket.Accept(rsp, req, &websocket.AcceptOptions{
		Subprotocols: []string{"shell"},
	})
	if err != nil {
		log.Warn("<devkit-client> failed to accept websocket:", err)
		return
	}
	c.SetReadLimit(8 * 1024 * 1024)
	// 通道双向对转
	ssh.serveShell(ctx, e, dst, NewWebsocket(ctx, c))
	c.Close(websocket.StatusNormalClosure, "")
	// c.CloseNow()
}

func (ssh *ServerShellHandler) prepareShell(ctx context.Context, shell *entity.ServerShell) (err error) {
	// 确保能够联通（内部可能通过 REGISTRY 进行地址查询和反向发包）
	var dst *stream.SessionStream
	dst, err = ssh.mgr.Acquire(ctx, &shell.Server)
	if err != nil {
		return
	}
	defer dst.CloseWrite()
	if err = app.Invoke(ctx, dst, "/server/query", &shell.Server, &shell.Server); err != nil {
		return
	}
	shell.Shell.ID = entity.ShellID(uuid.New().String())
	shell.Server.Address = dst.Peer.Address
	ssh.put(shell)
	return nil
}

func (ssh *ServerShellHandler) serveShell(ctx context.Context, e *entity.ServerShell, dst *stream.SessionStream, src *Websocket) (err error) {
	// if r == os.Stdin { // 对直接透传的 Shell 设定当前 Stdin 状态
	// 	state, _ := term.MakeRaw(int(os.Stdin.Fd()))
	// 	e.Cols, e.Rows, _ = term.GetSize(int(os.Stdin.Fd()))
	// 	defer term.Restore(int(os.Stdin.Fd()), state)
	// }
	log.DebugContext(ctx, "<devkit-client> shell started: ", &dst)
	io.WriteString(dst, "/shell/start:")
	json.NewEncoder(dst).Encode(e)

	go func(ctx context.Context) {
		_, err = io.Copy(dst, src)
		dst.CloseWrite()
		// dst.CloseRead()
	}(ctx)
	_, err = io.Copy(src, dst)
	src.CloseWrite()
	log.DebugContext(ctx, "<devkit-client> shell closed.")
	return
}

func (ssh *ServerShellHandler) HandleResize(rsp http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	shell_id := entity.ShellID(req.PathValue("shell_id"))
	e := ssh.get(shell_id)
	if e == nil {
		log.Warn("<devkit-client> unable to find shell:", shell_id)
		ssh.Respond(rsp, entity.ErrSessionNotFound)
		return
	}
	json.NewDecoder(req.Body).Decode(e)
	dst, err := ssh.mgr.Acquire(ctx, &e.Server)
	if err != nil {
		log.Warn("<devkit-client> failed acquire session:", err)
		ssh.Respond(rsp, err)
		return
	}
	defer dst.CloseWrite()
	err = app.Invoke(ctx, dst, "/shell/resize", e, nil)
	ssh.Respond(rsp, err)
}

func (ssh *ServerShellHandler) HandleRun(rsp http.ResponseWriter, req *http.Request) {

}
