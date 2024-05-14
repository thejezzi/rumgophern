package main

import (
	"context"
	"fmt"
	"gocontext/lib"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ContextKey string

const CK_TEST ContextKey = "test"

func main() {
	listener, err := net.Listen("tcp", "[::1]:3000")
	if err != nil {
		panic("could not create listener")
	}

	srvAwaiter := lib.NewAwaiter()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go runServer(ctx, listener, srvAwaiter)
	srvAwaiter.Await()
}

func runServer(ctx context.Context, listener net.Listener, awaiter *lib.Awaiter) {
	server := lib.NewServer()
	server.Middleware(lib.LogRequestMiddleware)

	server.GET("/", index)
	fmt.Println("listening on " + listener.Addr().String())
	if err := server.Run(ctx, listener); err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
	awaiter.Done()
}

func index(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res.Write([]byte("Hi there"))
}
