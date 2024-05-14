package lib

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type MiddlewareFunc func(next httprouter.Handle) httprouter.Handle

type Server struct {
	httpServer  *http.Server
	router      *httprouter.Router
	middlewares []MiddlewareFunc
}

func NewServer() *Server {
	router := httprouter.New()
	return &Server{
		router: router,
		httpServer: &http.Server{
			Handler: router,
		},
	}
}

func (s *Server) applyMiddlewares(handler httprouter.Handle) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		for _, middleware := range s.middlewares {
			middleware(handler)(res, req, params)
		}
	}
}

func (s *Server) GET(path string, handler httprouter.Handle) {
	s.router.GET(path, s.applyMiddlewares(handler))
}

func (s *Server) POST(path string, handler httprouter.Handle) {
	s.router.POST(path, s.applyMiddlewares(handler))
}

func (s *Server) PUT(path string, handler httprouter.Handle) {
	s.router.PUT(path, s.applyMiddlewares(handler))
}

func (s *Server) DELETE(path string, handler httprouter.Handle) {
	s.router.DELETE(path, s.applyMiddlewares(handler))
}

func (s *Server) PATCH(path string, handler httprouter.Handle) {
	s.router.PATCH(path, s.applyMiddlewares(handler))
}

func (s *Server) OPTIONS(path string, handler httprouter.Handle) {
	s.router.OPTIONS(path, s.applyMiddlewares(handler))
}

func (s *Server) HEAD(path string, handler httprouter.Handle) {
	s.router.HEAD(path, s.applyMiddlewares(handler))
}

func (s *Server) ServeHTTP(ctx context.Context, listener net.Listener) error {
	s.httpServer.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	return s.httpServer.Serve(listener)
}

func (s *Server) Run(ctx context.Context, listener net.Listener) error {
	// Start the server in a new goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- s.ServeHTTP(ctx, listener)
	}()

	// Listen for the context to be done and shutdown the server gracefully
	select {
	case <-ctx.Done():
		fmt.Println("Shutting down server...")
		time.Sleep(5 * time.Second)
		return s.Shutdown(context.Background())
	case err := <-serverErrCh:
		return err
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Middleware(middleware MiddlewareFunc) {
	s.middlewares = append(s.middlewares, middleware)
}

func LogRequestMiddleware(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	b := strings.Builder{}
	b.WriteString(req.Method)
	b.WriteString(" ")
	b.WriteString(req.URL.Path)
	b.WriteString(" ")
	b.WriteString(req.Proto)
	b.WriteString("\n")
	b.WriteString("Host: ")
	b.WriteString(req.Host)
	b.WriteString("\n")
	for k, v := range req.Header {
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(strings.Join(v, ", "))
		b.WriteString("\n")
	}
	fmt.Println(b.String())
}
