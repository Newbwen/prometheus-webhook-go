package core

import (
	"github.com/Newbwen/prometheus-webhook-go/initialize"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type server interface {
	ListenAndServe() error
}

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Minute
	s.WriteTimeout = 10 * time.Minute
	s.MaxHeaderBytes = 1 << 20
	return s
}

func RunServer() {
	Router := initialize.Routers()
	s := initServer(":8080", Router)
	log.Fatal(s.ListenAndServe())
}
