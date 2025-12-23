package navigation

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	navigation, err := h.service.Create(&request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, navigation)
}

func (h *Handler) GetAll(c *gin.Context) {
	navigations, err := h.service.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, navigations)
}

func (h *Handler) GetById(c *gin.Context) {
	id := c.Param("id")
	navigation, err := h.service.GetById(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, navigation)
}

func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	navigation, err := h.service.GetBySlug(slug)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, navigation)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	navigation, err := h.service.Update(id, &request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, navigation)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Navigation deleted successfully"})
}
