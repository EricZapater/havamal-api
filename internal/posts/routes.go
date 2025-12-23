package posts

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/posts", handler.Create)
	router.GET("/posts/:id", handler.GetPost)
	router.GET("/posts", handler.GetPosts)
	router.PUT("/posts/:id", handler.UpdatePost)
	router.DELETE("/posts/:id", handler.DeletePost)
	router.POST("/posts/category", handler.AddCategory)
	router.DELETE("/posts/category", handler.DeleteCategory)
	router.POST("/posts/version", handler.AddVersion)
	router.DELETE("/posts/version", handler.DeleteVersion)
}

func RegisterPublicRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/author/:author_id", handler.GetPostsByAuthor)
	router.GET("/slug/:slug", handler.GetPostBySlug)
	router.GET("/category/:category", handler.GetSummariesByCategory)
	router.GET("/posts/published", handler.GetPublishedPosts)
}