package versions

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/versions", handler.Create)
	router.PUT("/versions/:id", handler.Update)
	router.DELETE("/versions/:id", handler.Delete)
}

func RegisterPublicRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/versions", handler.GetAll)
	router.GET("/versions/:id", handler.GetById)
}