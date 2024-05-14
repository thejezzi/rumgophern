package main

import (
	"context"
	"fmt"
	"gocontext/config"
	"gocontext/lib"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
)

func main() {
	listener, err := net.Listen("tcp", "[::1]:3000")
	if err != nil {
		panic("could not create listener")
	}
	shutdownComplete := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	srvAwaiter := lib.NewAwaiter()
	ctx, cancel := context.WithCancel(context.Background())
	ctx = lib.InitConfig(ctx)

	go runServer(ctx, listener, srvAwaiter)
	go await(ctx, shutdownComplete, srvAwaiter)

	<-c
	fmt.Println(" Received shutdown signal, shutting down gracefully...")
	cancel()

	select {
	case <-c:
		fmt.Println(" Received second signal, exiting immediately")
		os.Exit(1)
	case <-shutdownComplete:
		fmt.Println("Shutdown complete, exiting")
	}
}

func await(ctx context.Context, done chan struct{}, awaiters ...*lib.Awaiter) {
	<-ctx.Done()
	for _, awaiter := range awaiters {
		<-awaiter.Await()
	}
	fmt.Println("All components shut down successfully")
	close(done)
}

func runServer(ctx context.Context, listener net.Listener, awaiter *lib.Awaiter) {
	server := lib.NewServer()

	server.Middleware(func(next httprouter.Handle) httprouter.Handle {
		return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
			c, ok := ctx.Value(config.MyKey).(*config.MyInnerConfig)
			if !ok {
				fmt.Println("Could not get config from context")
				next(res, req, params)
			}

			nctx := context.WithValue(req.Context(), config.MyKey, c)
			next(res, req.WithContext(nctx), params)
		}
	})

	server.GET("/", index)
	fmt.Println("listening on " + listener.Addr().String())
	if err := server.Run(ctx, listener); err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
	awaiter.Done()
}

func index(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	val, ok := req.Context().Value(config.MyKey).(*config.MyInnerConfig)
	if !ok {
		fmt.Println("Could not get config from context")
		return
	}
	s := val.String()
	res.Write([]byte(s))
}
