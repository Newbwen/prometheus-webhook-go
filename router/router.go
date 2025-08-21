package router

import (
	"github.com/Newbwen/prometheus-webhook-go/api"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Prometheus Webhook API
	prometheus := r.Group("")
	{
		prometheus.POST("/webhook", api.ApiGroupApp.PrometheusApiGroup.AlertApi.WebhookHandler)
	}

	return r
}
