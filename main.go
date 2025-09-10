package main

import (
	"github.com/Newbwen/prometheus-webhook-go/core"
	"github.com/Newbwen/prometheus-webhook-go/service"
)

func main() {
	//初始化k8s服务
	service.Init()
	//初始化路由并启动服务
	core.RunServer()
}
