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

// MachineUpdateRequest represents a request to update machine labels or extensions
type MachineUpdateRequest struct {
	Labels     map[string]string `json:"labels,omitempty"`
	Extensions []string           `json:"extensions,omitempty"`
	Maintenance *bool             `json:"maintenance,omitempty"`
}

// MachineWriteHandler handles machine write operations
type MachineWriteHandler struct {
	state      state.State
	management interface{} // Management service client
}

// NewMachineWriteHandler creates a new MachineWriteHandler
func NewMachineWriteHandler(s state.State, mgmt interface{}) *MachineWriteHandler {
	return &MachineWriteHandler{
		state:      s,
		management: mgmt,
	}
}

// UpdateMachine godoc
// @Summary      Update a machine
// @Description  Update machine labels, extensions, or maintenance mode
// @Tags         machines
// @Accept       json
// @Produce      json
// @Param        id       path      string              true  "Machine ID"
// @Param        machine  body      MachineUpdateRequest  true  "Machine update request"
// @Success      200      {object}  MachineResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /machines/{id} [patch]
func (h *MachineWriteHandler) UpdateMachine(c *gin.Context) {
	id := c.Param("id")
	var req MachineUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify machine exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	// Update machine using Management service
	// For labels: managementClient.UpdateMachineLabels(ctx, machineID, labels)
	// For extensions: managementClient.UpdateMachineExtensions(ctx, machineID, extensions)
	// For maintenance: managementClient.SetMachineMaintenance(ctx, machineID, enabled)
	log.Printf("Updating machine %s (Management service integration needed)", id)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Machine update initiated",
		"id": id,
		"note": "Management service integration required for actual update",
	})
}
