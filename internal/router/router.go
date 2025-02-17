package router

import (
	"atomwoz.com/remitly_task/internal/controller"
	"github.com/gin-gonic/gin"
)

// CreateRouter creates a new router with the specified prefix.
func CreateRouter(prefix string) *gin.Engine {
	router := gin.Default()
	api := router.Group(prefix)
	{
		api.GET("swift-codes/:code", controller.SwiftCode)
		api.GET("swift-codes/country/:country_code", controller.GetByCountry)
		api.POST("swift-codes/", controller.PostNewSwiftRow)
	}
	return router
}
