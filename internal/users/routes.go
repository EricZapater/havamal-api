package users

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/users", handler.Create)
	router.GET("/users", handler.FindAll)
	router.GET("/users/:id", handler.FindByID)
	router.PUT("/users/:id", handler.Update)
	router.DELETE("/users/:id", handler.Delete)
}