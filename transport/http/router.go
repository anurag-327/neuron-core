package httpTransport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.POST("/execute", handler.Execute)
		v1.GET("/status", handler.Status)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Healthy")
	})
	return router
}
