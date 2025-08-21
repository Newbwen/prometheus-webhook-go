package main

import (
	"fmt"
	"github.com/Newbwen/prometheus-webhook-go/router"
	"log"
)

func main() {
	r := router.InitRouter()
	fmt.Println("Gin Webhook server listening on :8080")
	log.Fatal(r.Run(":8080"))
}
