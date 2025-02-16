package router

import (
	"atomwoz.com/remitly_task/src/controller"
	"github.com/gin-gonic/gin"
)

func CreateRouter(prefix string) *gin.Engine {
	router := gin.Default()
	api := router.Group(prefix)
	{
		api.GET("swift-codes/:code", controller.SwiftCode)
	}
	return router
}
