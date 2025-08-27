package router

import (
	v1 "github.com/Newbwen/prometheus-webhook-go/api/v1"
	"github.com/gin-gonic/gin"
)

func InitWebhookRouter(Router *gin.RouterGroup) {
	webhookRouter := Router.Group("")
	{
		webhookRouter.POST("/webhook", v1.ApiGroupApp.WebhookApi.HandleWebhook)
		webhookRouter.GET("/health", v1.ApiGroupApp.WebhookApi.HealthzCheck)
	}
}
