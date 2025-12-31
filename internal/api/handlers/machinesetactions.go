package handlers

import (
	"log"
	"net/http"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/gin-gonic/gin"
	omniresources "github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
)

// MachineSetActionsHandler handles machine set action operations
type MachineSetActionsHandler struct {
	state      state.State
	management interface{} // Management service client
}

// NewMachineSetActionsHandler creates a new MachineSetActionsHandler
func NewMachineSetActionsHandler(s state.State, mgmt interface{}) *MachineSetActionsHandler {
	return &MachineSetActionsHandler{
		state:      s,
		management: mgmt,
	}
}

// TriggerDestroy godoc
// @Summary      Trigger machine set destruction
// @Description  Trigger destruction/teardown of a machine set
// @Tags         machinesets
// @Produce      json
// @Param        id   path      string  true  "Machine set ID"
// @Success      202  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machinesets/{id}/actions/destroy [post]
func (h *MachineSetActionsHandler) TriggerDestroy(c *gin.Context) {
	id := c.Param("id")

	// Verify machine set exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineSetType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine set not found"})
		return
	}

	// Use Management service to trigger machine set destruction
	log.Printf("Triggering destruction for machine set %s (Management service integration needed)", id)
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Machine set destruction initiated",
		"machine_set_id": id,
		"note": "Management service integration required for actual destruction",
	})
}
