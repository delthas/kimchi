package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	Listeners []*Listener

	http1         *http.Server
	http1Listener *pipeListener
}

func NewServer() *Server {
	srv := &Server{}
	srv.http1Listener = newPipeListener()
	srv.http1 = &http.Server{
		Handler: srv,
	}
	return srv
}

func (srv *Server) Start() error {
	for _, ln := range srv.Listeners {
		if err := ln.Start(); err != nil {
			return err
		}
	}

	go func() {
		if err := srv.http1.Serve(srv.http1Listener); err != nil {
			log.Fatalf("HTTP/1 server: %v", err)
		}
	}()

	return nil
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "yo", 200)
}

type Listener struct {
	Network string
	Address string
	Server  *Server
}

func (ln *Listener) Start() error {
	netLn, err := net.Listen(ln.Network, ln.Address)
	if err != nil {
		return err
	}
	log.Printf("HTTP server listening on %q", ln.Address)

	go func() {
		if err := ln.serve(netLn); err != nil {
			log.Fatalf("failed to serve listener %q: %v", ln.Address, err)
		}
	}()

	return nil
}

func (ln *Listener) serve(netLn net.Listener) error {
	for {
		conn, err := netLn.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %v", err)
		}

		go func() {
			if err := ln.serveConn(conn); err != nil {
				log.Printf("listener %q: %v", ln.Address, err)
			}
		}()
	}
}

func (ln *Listener) serveConn(conn net.Conn) error {
	return ln.Server.http1Listener.ServeConn(conn)
}

var errPipeListenerClosed = fmt.Errorf("pipe listener closed")

// pipeListener is a hack to workaround the lack of http.Server.ServeConn.
// See: https://github.com/golang/go/issues/36673
type pipeListener struct {
	ch     chan net.Conn
	closed bool
	mu     sync.Mutex
}

func newPipeListener() *pipeListener {
	return &pipeListener{
		ch: make(chan net.Conn, 64),
	}
}

func (ln *pipeListener) Accept() (net.Conn, error) {
	conn, ok := <-ln.ch
	if !ok {
		return nil, errPipeListenerClosed
	}
	return conn, nil
}

func (ln *pipeListener) Close() error {
	ln.mu.Lock()
	defer ln.mu.Unlock()

	if ln.closed {
		return errPipeListenerClosed
	}
	ln.closed = true
	close(ln.ch)
	return nil
}

func (ln *pipeListener) ServeConn(conn net.Conn) error {
	ln.mu.Lock()
	defer ln.mu.Unlock()

	if ln.closed {
		return errPipeListenerClosed
	}
	ln.ch <- conn
	return nil
}

func (ln *pipeListener) Addr() net.Addr {
	return pipeAddr{}
}

type pipeAddr struct{}

func (pipeAddr) Network() string {
	return "pipe"
}

func (pipeAddr) String() string {
	return "pipe"
}
