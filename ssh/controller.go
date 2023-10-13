package ssh

import (
	"context"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

type Controller struct {
	controllerClientManager
}

func NewController() *Controller {
	var c Controller
	prepareControllerClientManager(&c.controllerClientManager)
	return &c
}


func (c *Controller) CreateShell(ctx context.Context, req Request) (session *Session, err error) {
	session = &Session{ Req: req }
	var sc *ssh.Client
	var ss *ssh.Session
	for retry := 0; retry < 3; retry ++ {
		if sc, err = c.FetchClient(req); err != nil {
			time.Sleep(1200 * time.Millisecond)
			continue
		}
		if ss, err = sc.NewSession(); err != nil {
			c.CloseClient(req, sc)
			time.Sleep(700 * time.Millisecond)
			continue
		}
		if err = ss.RequestPty("xterm-256color", req.Rows, req.Cols, ssh.TerminalModes{}); err != nil {
			sc.Close()
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		log.Println("<ssh> failed to create session: ", err)
		return
	}
	session.session = ss
	return
}

