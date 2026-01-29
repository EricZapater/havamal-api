package images

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the protected image routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	images := router.Group("/images")
	{
		images.POST("/upload", handler.UploadImage)
	}
}
