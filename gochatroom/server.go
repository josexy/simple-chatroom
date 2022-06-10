package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	address string
	app     *gin.Engine
}

func NewServer(port int, app *gin.Engine) *Server {
	return &Server{
		address: fmt.Sprintf(":%d", port),
		app:     app,
	}
}

func (svr *Server) Run() {
	server := &http.Server{
		Addr:    svr.address,
		Handler: svr.app,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// Timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
