package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/gin-gonic/gin"
	omniresources "github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
)

// HealthResponse represents the health status of the API
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version"`
	Omni      OmniHealthStatus  `json:"omni,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// OmniHealthStatus represents the health status of the Omni connection
type OmniHealthStatus struct {
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
}

// HealthHandler handles health check requests
type HealthHandler struct {
	state state.State
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(s state.State) *HealthHandler {
	return &HealthHandler{state: s}
}

// GetHealth godoc
// @Summary      Get API health status
// @Description  Get the health status of the API server and Omni connection
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func (h *HealthHandler) GetHealth(c *gin.Context) {
	status := "healthy"
	omniHealth := OmniHealthStatus{Connected: false}

	// Check Omni connectivity by attempting to list clusters (lightweight operation)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterType, "", resource.VersionUndefined)
	_, err := h.state.List(ctx, md)
	if err != nil {
		status = "degraded"
		omniHealth.Error = err.Error()
	} else {
		omniHealth.Connected = true
	}

	resp := HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "0.0.1",
		Omni:      omniHealth,
		Links: map[string]string{
			"self":    buildURL(c, "/health"),
			"metrics": buildURL(c, "/metrics"),
		},
	}

	c.JSON(http.StatusOK, resp)
}
