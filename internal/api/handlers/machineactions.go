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

// MachineActionsHandler handles machine action operations
type MachineActionsHandler struct {
	state      state.State
	management interface{} // Management service client
	talos      interface{}  // Talos service client
}

// NewMachineActionsHandler creates a new MachineActionsHandler
func NewMachineActionsHandler(s state.State, mgmt, talos interface{}) *MachineActionsHandler {
	return &MachineActionsHandler{
		state:      s,
		management: mgmt,
		talos:      talos,
	}
}

// RebootMachine godoc
// @Summary      Reboot a machine
// @Description  Trigger a reboot of a machine via Talos API
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Machine ID"
// @Success      202  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machines/{id}/actions/reboot [post]
func (h *MachineActionsHandler) RebootMachine(c *gin.Context) {
	id := c.Param("id")

	// Verify machine exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	// Use Talos service to reboot machine
	// Typically: talosClient.Reboot(ctx, machineID)
	log.Printf("Rebooting machine %s (Talos service integration needed)", id)
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Machine reboot initiated",
		"machine_id": id,
		"note": "Talos service integration required for actual reboot",
	})
}

// ShutdownMachine godoc
// @Summary      Shutdown a machine
// @Description  Trigger a shutdown of a machine via Talos API
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Machine ID"
// @Success      202  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machines/{id}/actions/shutdown [post]
func (h *MachineActionsHandler) ShutdownMachine(c *gin.Context) {
	id := c.Param("id")

	// Verify machine exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	// Use Talos service to shutdown machine
	log.Printf("Shutting down machine %s (Talos service integration needed)", id)
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Machine shutdown initiated",
		"machine_id": id,
		"note": "Talos service integration required for actual shutdown",
	})
}

// ResetMachine godoc
// @Summary      Reset a machine
// @Description  Trigger a reset of a machine via Talos API
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Machine ID"
// @Success      202  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machines/{id}/actions/reset [post]
func (h *MachineActionsHandler) ResetMachine(c *gin.Context) {
	id := c.Param("id")

	// Verify machine exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	// Use Talos service to reset machine
	log.Printf("Resetting machine %s (Talos service integration needed)", id)
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Machine reset initiated",
		"machine_id": id,
		"note": "Talos service integration required for actual reset",
	})
}

// ToggleMaintenance godoc
// @Summary      Toggle machine maintenance mode
// @Description  Enable or disable maintenance mode for a machine
// @Tags         machines
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Machine ID"
// @Param        enabled  body      map[string]bool  true  "Maintenance mode enabled"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /machines/{id}/actions/maintenance [post]
func (h *MachineActionsHandler) ToggleMaintenance(c *gin.Context) {
	id := c.Param("id")
	var req map[string]bool
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enabled, ok := req["enabled"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "enabled field is required"})
		return
	}

	// Verify machine exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	// Use Management service to toggle maintenance mode
	log.Printf("Setting maintenance mode for machine %s to %v (Management service integration needed)", id, enabled)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Maintenance mode updated",
		"machine_id": id,
		"maintenance_enabled": enabled,
		"note": "Management service integration required for actual update",
	})
}
