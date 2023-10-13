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

RECONNECT:
	for i:=0;i<3;i++ {
		if sc, err = c.FetchClient(req); err != nil {
			log.Println("<ssh> failed to fetch client: ", err)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		return
	}
	for i:=0;i<3;i++ {
		if ss, err = sc.NewSession(); err != nil {
			c.CloseClient(req, sc)
			goto RECONNECT
		}
		break
	}
	if err != nil {
		log.Println("<ssh> failed to create session: ", err)
		return nil, err
	}
	if err = ss.RequestPty("xterm-256color", req.Rows, req.Cols, ssh.TerminalModes{}); err != nil {
		sc.Close()
		return nil, err
	}
	session.session = ss
	return
}

