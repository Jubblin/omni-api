package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	omniresources "github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
)

// ConfigPatchCreateRequest represents a request to create a config patch
type ConfigPatchCreateRequest struct {
	ID      string `json:"id" binding:"required"`
	Cluster string `json:"cluster" binding:"required"`
	Data    string `json:"data" binding:"required"`
}

// ConfigPatchUpdateRequest represents a request to update a config patch
type ConfigPatchUpdateRequest struct {
	Data string `json:"data" binding:"required"`
}

// ConfigPatchWriteHandler handles config patch write operations
type ConfigPatchWriteHandler struct {
	state      state.State
	management interface{} // Management service client
}

// NewConfigPatchWriteHandler creates a new ConfigPatchWriteHandler
func NewConfigPatchWriteHandler(s state.State, mgmt interface{}) *ConfigPatchWriteHandler {
	return &ConfigPatchWriteHandler{
		state:      s,
		management: mgmt,
	}
}

// CreateConfigPatch godoc
// @Summary      Create a new config patch
// @Description  Create a new config patch for a cluster
// @Tags         configpatches
// @Accept       json
// @Produce      json
// @Param        patch  body      ConfigPatchCreateRequest  true  "Config patch creation request"
// @Success      201    {object}  ConfigPatchResponse
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /configpatches [post]
func (h *ConfigPatchWriteHandler) CreateConfigPatch(c *gin.Context) {
	var req ConfigPatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use Management service to create config patch
	log.Printf("Creating config patch %s for cluster %s (Management service integration needed)", req.ID, req.Cluster)
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Config patch creation initiated",
		"id": req.ID,
		"cluster": req.Cluster,
		"note": "Management service integration required for actual creation",
	})
}

// UpdateConfigPatch godoc
// @Summary      Update a config patch
// @Description  Update an existing config patch
// @Tags         configpatches
// @Accept       json
// @Produce      json
// @Param        id     path      string                  true  "Config patch ID"
// @Param        patch  body      ConfigPatchUpdateRequest  true  "Config patch update request"
// @Success      200    {object}  ConfigPatchResponse
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /configpatches/{id} [put]
func (h *ConfigPatchWriteHandler) UpdateConfigPatch(c *gin.Context) {
	id := c.Param("id")
	var req ConfigPatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify config patch exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ConfigPatchType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "config patch not found"})
		return
	}

	// Use Management service to update config patch
	log.Printf("Updating config patch %s (Management service integration needed)", id)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Config patch update initiated",
		"id": id,
		"note": "Management service integration required for actual update",
	})
}

// DeleteConfigPatch godoc
// @Summary      Delete a config patch
// @Description  Delete a config patch from Omni
// @Tags         configpatches
// @Produce      json
// @Param        id   path      string  true  "Config patch ID"
// @Success      202  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /configpatches/{id} [delete]
func (h *ConfigPatchWriteHandler) DeleteConfigPatch(c *gin.Context) {
	id := c.Param("id")

	// Verify config patch exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ConfigPatchType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "config patch not found"})
		return
	}

	// Use Management service to delete config patch
	log.Printf("Deleting config patch %s (Management service integration needed)", id)
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Config patch deletion initiated",
		"id": id,
		"note": "Management service integration required for actual deletion",
	})
}
