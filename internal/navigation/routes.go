package navigation

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/navigation", handler.Create)			
	router.PUT("/navigation/:id", handler.Update)
	router.DELETE("/navigation/:id", handler.Delete)
}

func RegisterPublicRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/navigation", handler.GetAll)
	router.GET("/navigation/:id", handler.GetById)
	router.GET("/navigation/slug/:slug", handler.GetBySlug)
}
