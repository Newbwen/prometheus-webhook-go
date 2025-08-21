package router

import "github.com/gin-gonic/gin"

type RouterGroup struct {
}

func (r *RouterGroup) Init(Router *gin.Engine) {
	apiGroup := Router.Group("")
	{
		InitWebhookRouter(apiGroup)
	}
}

var RouterGroupApp = new(RouterGroup)
