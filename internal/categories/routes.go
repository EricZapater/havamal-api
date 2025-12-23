package categories

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/categories", handler.Create)
	router.PUT("/categories/:id", handler.Update)
	router.DELETE("/categories/:id", handler.Delete)
}

func RegisterPublicRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/categories", handler.GetAll)
	router.GET("/categories/:id", handler.GetById)
	router.GET("/categories/slug/:slug", handler.GetBySlug)
}

