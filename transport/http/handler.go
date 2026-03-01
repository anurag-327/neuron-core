package httpTransport

import (
	"net/http"

	"github.com/anurag-327/neuron-core/engine"
	"github.com/anurag-327/neuron-core/pkg/api"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	executionService *engine.ExecutionService
}

func NewHandler(executionService *engine.ExecutionService) *Handler {
	return &Handler{executionService: executionService}
}

func (h *Handler) Execute(c *gin.Context) {
	ctx := c.Request.Context()
	var body api.ExecuteRequest

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateExecuteParams(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.executionService.Execute(ctx, engine.ExecuteSpec{
		Code:     body.Code,
		Language: body.Language,
		Input:    body.Input,
		Limit: engine.Limit{
			MemoryKB: body.Limit.MemoryKB,
			TimeMs:   body.Limit.TimeMs,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Status(c *gin.Context) {
	err := h.executionService.Health(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"runner": "down", "runnerError": err.Error(), "runnerVersion": "1.0.0"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"runner": "up", "runnerError": "", "runnerVersion": "1.0.0"})
}
