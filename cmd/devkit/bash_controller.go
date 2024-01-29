package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/terrywh/devkit/k8s"
	"github.com/terrywh/devkit/ssh"
)

type BashController struct {
	mutex sync.Mutex
	cache map[string]BashCache
}

type BashCache struct {
	Bash BashBackend
}

func NewBashController() *BashController {
	return &BashController{
		mutex: sync.Mutex{},
		cache: make(map[string]BashCache),
	}
}

// type BashRequest struct {
// 	SSH ssh.Request
// 	K8S k8s.Request
// }

var ErrShellNotExists error = errors.New("shell does not exist")

func (bc *BashController) FetchShell(ctx context.Context, r *http.Request) (bash BashBackend, err error) {
	key := r.URL.Query().Get("key")
	if bash, err = bc.fetch(ctx, key); err == nil {
		return
	}
	if errors.Is(err, ErrShellNotExists) {
		bash, err = bc.create(ctx, r)
	}

	if err != nil {
		return
	}
	bc.store(ctx, key, bash)
	return
}

func UnmarshalTo(q url.Values, r io.Reader, d interface{}) (err error) {
	decoder := json.NewDecoder(r)
	if err = decoder.Decode(d); err != nil {
		return
	}
	v := reflect.ValueOf(d)

	rows, _ := strconv.Atoi(q.Get("rows"))
	cols, _ := strconv.Atoi(q.Get("cols"))
	if rows == 0 {
		rows = 60
	}
	if cols == 0 {
		cols = 80
	}

	v.Elem().FieldByName("Cols").SetInt(int64(cols))
	v.Elem().FieldByName("Rows").SetInt(int64(rows))
	return
}

func (bc *BashController) create(ctx context.Context, r *http.Request) (bash BashBackend, err error) {
	query := r.URL.Query()
	shell := query.Get("type")

	switch shell {
	case "ssh":
		var req ssh.Request
		if err = UnmarshalTo(query, r.Body, &req); err != nil {
			log.Println("failed to parse shell create: ", err, req)
			return
		}
		log.Println("create ssh shell: ", req)
		bash, err = defaultSSHController.CreateShell(ctx, req)
	case "k8s":
		var req k8s.Request
		if err = UnmarshalTo(query, r.Body, &req); err != nil {
			log.Println("failed to parse shell create: ", err, req)
			return
		}
		log.Println("create k8s shell: ", req)
		bash, err = defaultK8SController.CreateShell(ctx, req)
	default:
		bash = nil
		err = fmt.Errorf("invalid arguments or mising 'key'?")
	}
	return
}

func (bc *BashController) key(form url.Values) string {
	hash := md5.New()
	io.WriteString(hash, form.Encode())
	fmt.Fprint(hash, time.Now())
	return hex.EncodeToString(hash.Sum(nil))
}

func (bc *BashController) fetch(ctx context.Context, key string) (BashBackend, error) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if cache, ok := bc.cache[key]; ok {
		return cache.Bash, nil
	}
	return nil, ErrShellNotExists
}

func (bc *BashController) store(ctx context.Context, key string, bash BashBackend) (err error) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	bc.cache[key] = BashCache{bash}
	return
}
