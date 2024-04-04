package stream

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/quic-go/quic-go"
)

type ServerStreamHandler interface {
	ServeStream(ctx context.Context, r *bufio.Reader, w io.WriteCloser)
}

type ServerStreamHandlerFn struct {
	fn func(ctx context.Context, r *bufio.Reader, w io.WriteCloser)
}

func (shf ServerStreamHandlerFn) ServeStream(ctx context.Context, r *bufio.Reader, w io.WriteCloser) {
	shf.fn(ctx, r, w)
}

type ServeMux struct {
	handler map[string]ServerStreamHandler
}

func NewServeMux() (mux *ServeMux) {
	mux = &ServeMux{
		handler: make(map[string]ServerStreamHandler),
	}
	return mux
}

func (mux ServeMux) Handle(path string, handler ServerStreamHandler) {
	mux.handler[path] = handler
}

func (mux ServeMux) HandleFunc(path string, fn func(ctx context.Context, r *bufio.Reader, w io.WriteCloser)) {
	mux.handler[path] = ServerStreamHandlerFn{fn}
}

func (mux ServeMux) ServeStream(ctx context.Context, s quic.Stream) {
	defer s.Close()
	r := bufio.NewReader(s)
	path, err := r.ReadString(':')
	if err != nil {
		log.Print(err)
		return
	}
	path = path[:len(path)-1]

	if handler, found := mux.handler[path]; found {
		handler.ServeStream(ctx, r, s)
	} else {
		log.Println("<ServeMux.ServeStream> handle not found for path: ", path)
		fmt.Fprintln(s, "invalid path")
		s.CancelRead(quic.StreamErrorCode(10000))
	}
}

type Server struct {
	listener *quic.Listener
	handler  ServerHandler
	mutex    *sync.Mutex
	conn     map[uint64]quic.Connection
	connid   atomic.Uint64
}

func newServer(listener *quic.Listener, handler ServerHandler) (qs *Server, err error) {
	qs = &Server{
		listener: listener,
		handler:  handler,
		mutex:    &sync.Mutex{},
		conn:     make(map[uint64]quic.Connection),
	}
	return
}

func (svr *Server) Serve(ctx context.Context) {
	defer svr.listener.Close()
	log.Println("<QuicService.Serve> accepting... ")
SERVING:
	for {
		conn, err := svr.listener.Accept(ctx)
		if err != nil {
			log.Print(err)
			break SERVING
		}
		go svr.ServeConn(context.Background(), conn)
	}
	log.Println("<QuicService.Serve> closed.")
}

func (svr *Server) put(conn quic.Connection) (connid uint64) {
	svr.mutex.Lock()
	defer svr.mutex.Unlock()

	connid = svr.connid.Add(1)
	svr.conn[connid] = conn
	return
}

func (svr *Server) del(connid uint64) {
	svr.mutex.Lock()
	defer svr.mutex.Unlock()

	delete(svr.conn, connid)
}

func (svr *Server) ServeConn(ctx context.Context, conn quic.Connection) {
	log.Println("<QuicService.ServeConn> connection started: ", &conn)
	connid := svr.put(conn)
SERVING:
	for {
		s, err := conn.AcceptStream(ctx)
		if ae, ok := err.(*quic.ApplicationError); ok && ae.ErrorCode == quic.ApplicationErrorCode(0) {
			break SERVING
		} else if err != nil {
			if _, ok := err.(*quic.ApplicationError); !ok {
				log.Print("failed to accept stream: ", err)
			}
			return
		}
		go svr.handler.ServeStream(context.Background(), s)
	}
	svr.del(connid)
	log.Println("<QuicService.ServeConn> connection closed: ", &conn)
}

func (svr *Server) Close() error {
	svr.mutex.Lock()
	defer svr.mutex.Unlock()
	for _, conn := range svr.conn {
		conn.CloseWithError(quic.ApplicationErrorCode(1), "application shutdown")
	}
	return svr.listener.Close()
}
