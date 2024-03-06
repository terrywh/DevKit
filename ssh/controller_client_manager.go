package ssh

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type controllerClientManager struct {
	client      map[string]*ssh.Client
	clientMutex *sync.RWMutex
	clientClose bool
}

func prepareControllerClientManager(c *controllerClientManager) {
	c.client = make(map[string]*ssh.Client)
	c.clientMutex = &sync.RWMutex{}
	c.clientClose = false
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ticker := time.NewTicker(65 * time.Second)
		defer ticker.Stop()
	SERVING:
		for {
			select {
			case <-ctx.Done():
				break SERVING
			case <-ticker.C:
				c.CleanupClient(ctx)
			}
		}
	}()
}

func (c *controllerClientManager) prepareClient(req Request) (cli *ssh.Client, err error) {
	var clients []*ssh.Client
	for _, route := range req.Route {
		var conn net.Conn
		var addr string
		addr = fmt.Sprintf("%s:%d", route.Host, route.Port)
		if cli == nil {
			conn, err = net.Dial("tcp", addr)
		} else {
			conn, err = cli.Dial("tcp", addr)
		}
		if err != nil {
			break
		}
		var cconn ssh.Conn
		var cchan <-chan ssh.NewChannel
		var creqs <-chan *ssh.Request
		cconf := &ssh.ClientConfig{
			User:            route.User,
			Auth:            c.prepareAuth(route.Pass),
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		cconn, cchan, creqs, err = ssh.NewClientConn(conn, addr, cconf)
		if err != nil {
			break
		}
		cli = ssh.NewClient(cconn, cchan, creqs)
		clients = append(clients, cli)
	}
	if err != nil {
		for _, client := range clients {
			client.Close()
		}
	}
	return
}

func (c *controllerClientManager) prepareAuth(pass string) (auth []ssh.AuthMethod) {
	if method, err := c.makePublicKeyAuth(); err == nil {
		auth = append(auth, method)
	}
	if len(pass) > 0 {
		// 密码认证
		auth = append(auth, ssh.Password(pass))
		// 交互认证（动态密钥？）
		auth = append(auth, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
			if len(questions) == 0 {
				return []string{}, nil
			}
			return []string{pass}, nil
		}))
	}
	return
}

func (c *controllerClientManager) makePublicKeyAuth() (ssh.AuthMethod, error) {
	// 公钥认证
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	pkey, err := os.ReadFile(filepath.Join(home, ".ssh", "id_rsa"))
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(pkey)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func (c *controllerClientManager) FetchClient(req Request) (client *ssh.Client, err error) {
	if client = c.fetchClient1(req); client != nil {
		return
	}
	client, err = c.fetchClient2(req)
	return
}

func (c *controllerClientManager) fetchClient1(req Request) (client *ssh.Client) {
	c.clientMutex.RLock()
	defer c.clientMutex.RUnlock()
	client = c.client[req.Key()]
	return
}

func (c *controllerClientManager) CloseClient(req Request, client *ssh.Client) {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()
	client.Close()
	delete(c.client, req.Key())
}

func (c *controllerClientManager) fetchClient2(req Request) (client *ssh.Client, err error) {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()
	client, err = c.prepareClient(req)
	if err != nil {
		return
	}
	c.client[req.Key()] = client
	return
}

func (c *controllerClientManager) CleanupClient(ctx context.Context) {

}
