package initialize

import (
	"github.com/Newbwen/prometheus-webhook-go/api/v1"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Logger())
	RouterGroup := Router.Group("")
	ApiRouter := &ApiRouter{}
	ApiRouter.InitRouter(RouterGroup)
	return Router
}

type ApiRouter struct{}

func (ar *ApiRouter) InitRouter(Router *gin.RouterGroup) {
	router := Router.Group("/")
	{
		router.GET("health", v1.ApiGroupApp.HealthzCheck)
		router.POST("webhook", v1.ApiGroupApp.WebhookApi.HandleWebhook)
		router.POST("remsg", v1.ApiGroupApp.ReMessage)
	}
}
