package main

import (
	"fmt"
	"log"

	"github.com/Newbwen/prometheus-webhook-go/router"
	"github.com/Newbwen/prometheus-webhook-go/service"
	"github.com/gin-gonic/gin"
)

func main() {
	service.Init()
	r := gin.Default()
	router.RouterGroupApp.Init(r)

	fmt.Println("Gin Webhook server listening on :8080")
	log.Fatal(r.Run(":8080"))
}
