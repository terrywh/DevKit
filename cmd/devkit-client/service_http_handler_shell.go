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
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/stream"
	"nhooyr.io/websocket"
)

type ShellHandler struct {
	app.HttpHandlerBase
	mgr   stream.SessionManager
	mutex *sync.RWMutex
	shell map[entity.ShellID]*entity.ServerShell
}

func initHttpShellHandler(mgr stream.SessionManager, mux *http.ServeMux) *ShellHandler {
	css := &ShellHandler{
		mgr:   mgr,
		mutex: &sync.RWMutex{},
		shell: make(map[entity.ShellID]*entity.ServerShell),
	}

	mux.HandleFunc("/shell/prepare", css.HandlePrepare)
	mux.HandleFunc("/shell/{shell_id}/socket", css.HandleSocket)
	mux.HandleFunc("/shell/{shell_id}/resize", css.HandleResize)
	mux.HandleFunc("/shell/run", css.HandleRun)

	// TODO cleanup
	return css
}

func (css *ShellHandler) put(e *entity.ServerShell) {
	css.mutex.Lock()
	defer css.mutex.Unlock()
	css.shell[e.ShellId] = e
}

func (css *ShellHandler) get(shell_id entity.ShellID) *entity.ServerShell {
	css.mutex.RLock()
	defer css.mutex.RUnlock()

	return css.shell[shell_id]
}

func (css *ShellHandler) del(e *entity.ServerShell) {
	css.mutex.Lock()
	defer css.mutex.Unlock()
	delete(css.shell, e.ShellId)
}

func (css *ShellHandler) HandlePrepare(rsp http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	d := json.NewDecoder(req.Body)
	e := &entity.ServerShell{}
	if err := d.Decode(e); err != nil {
		css.Respond(rsp, err)
		return
	}
	if err := css.prepareShell(ctx, e); err != nil {
		css.Respond(rsp, err)
	} else {
		css.Respond(rsp, e)
	}
}

func (css *ShellHandler) HandleSocket(rsp http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	// 确认已注册的会话
	shell_id := entity.ShellID(req.PathValue("shell_id"))
	e := css.get(shell_id)
	if e == nil {
		infra.Warn("<devkit-client> unable to find shell:", shell_id)
		return
	}
	defer css.del(e)
	// 确认对应会话通道
	dst, err := css.mgr.Acquire(ctx, &e.Server)
	if err != nil {
		infra.Warn("<devkit-client> failed to acquire stream:", err)
		return
	}
	defer dst.CloseRead()
	defer dst.CloseWrite()
	// 确认对应前端通道
	c, err := websocket.Accept(rsp, req, &websocket.AcceptOptions{
		Subprotocols: []string{"shell"},
	})
	if err != nil {
		infra.Warn("<devkit-client> failed to accept websocket:", err)
		return
	}
	// 通道双向对转
	if err = css.serveShell(ctx, e, dst, NewWebSocketReader(ctx, c), NewWebSocketWriter(ctx, c)); err != nil {
		c.Close(websocket.StatusAbnormalClosure, err.Error())
	} else {
		c.Close(websocket.StatusNormalClosure, "")
	}
	// c.CloseNow()
}

func (css *ShellHandler) prepareShell(ctx context.Context, shell *entity.ServerShell) (err error) {
	// 确保能够联通（内部可能通过 REGISTRY 进行地址查询和反向发包）
	var ss *stream.SessionStream
	ss, err = css.mgr.Acquire(ctx, &shell.Server)
	if err != nil {
		return
	}
	if err = app.Invoke(ctx, ss, "/device/query", &shell.Server, &shell.Server); err != nil {
		return
	}
	shell.ShellId = entity.ShellID(uuid.New().String())
	css.put(shell)
	return nil

}

func (css *ShellHandler) serveShell(_ context.Context, e *entity.ServerShell, src *stream.SessionStream, r io.Reader, w io.WriteCloser) (err error) {
	// if r == os.Stdin { // 对直接透传的 Shell 设定当前 Stdin 状态
	// 	state, _ := term.MakeRaw(int(os.Stdin.Fd()))
	// 	e.Cols, e.Rows, _ = term.GetSize(int(os.Stdin.Fd()))
	// 	defer term.Restore(int(os.Stdin.Fd()), state)
	// }

	io.WriteString(src, "/shell/start:")
	json.NewEncoder(src).Encode(e)

	go func() {
		io.Copy(src, r)
		src.CloseWrite()
	}()
	_, err = io.Copy(w, src)
	w.Close()
	return
}

func (css *ShellHandler) HandleResize(rsp http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	shell_id := entity.ShellID(req.PathValue("shell_id"))
	e := css.get(shell_id)
	if e == nil {
		infra.Warn("<devkit-client> unable to find shell:", shell_id)
		css.Respond(rsp, entity.ErrSessionNotFound)
		return
	}
	json.NewDecoder(req.Body).Decode(e)
	ss, err := css.mgr.Acquire(ctx, &e.Server)
	if err != nil {
		infra.Warn("<devkit-client> failed acquire session:", err)
		css.Respond(rsp, err)
		return
	}
	err = app.Invoke(ctx, ss, "/shell/resize", e, nil)
	css.Respond(rsp, err)
}

func (css *ShellHandler) HandleRun(rsp http.ResponseWriter, req *http.Request) {

}
