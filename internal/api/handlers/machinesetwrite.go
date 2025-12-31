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

// MachineSetCreateRequest represents a request to create a machine set
type MachineSetCreateRequest struct {
	ID              string `json:"id" binding:"required"`
	Cluster         string `json:"cluster" binding:"required"`
	MachineClass    string `json:"machine_class" binding:"required"`
	BootstrapSpec   string `json:"bootstrap_spec,omitempty"`
	UpdateStrategy  string `json:"update_strategy,omitempty"`
	DeleteStrategy  string `json:"delete_strategy,omitempty"`
	MachineCount    uint32 `json:"machine_count,omitempty"`
}

// MachineSetUpdateRequest represents a request to update a machine set
type MachineSetUpdateRequest struct {
	MachineClass    string `json:"machine_class,omitempty"`
	UpdateStrategy  string `json:"update_strategy,omitempty"`
	DeleteStrategy  string `json:"delete_strategy,omitempty"`
	MachineCount    uint32 `json:"machine_count,omitempty"`
}

// MachineSetWriteHandler handles machine set write operations
type MachineSetWriteHandler struct {
	state      state.State
	management interface{} // Management service client
}

// NewMachineSetWriteHandler creates a new MachineSetWriteHandler
func NewMachineSetWriteHandler(s state.State, mgmt interface{}) *MachineSetWriteHandler {
	return &MachineSetWriteHandler{
		state:      s,
		management: mgmt,
	}
}

// CreateMachineSet godoc
// @Summary      Create a new machine set
// @Description  Create a new machine set in Omni
// @Tags         machinesets
// @Accept       json
// @Produce      json
// @Param        machineset  body      MachineSetCreateRequest  true  "Machine set creation request"
// @Success      201         {object}  map[string]string
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /machinesets [post]
func (h *MachineSetWriteHandler) CreateMachineSet(c *gin.Context) {
	var req MachineSetCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use Management service to create machine set
	// The Management service API needs to be integrated here
	// Typically: managementClient.CreateMachineSet(ctx, req.ID, req.Cluster, req.MachineClass, ...)
	log.Printf("Creating machine set %s for cluster %s (Management service integration needed)", req.ID, req.Cluster)
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Machine set creation initiated",
		"id": req.ID,
		"cluster": req.Cluster,
		"note": "Management service integration required for actual creation",
	})
}

// UpdateMachineSet godoc
// @Summary      Update a machine set
// @Description  Update an existing machine set in Omni
// @Tags         machinesets
// @Accept       json
// @Produce      json
// @Param        id          path      string                  true  "Machine set ID"
// @Param        machineset  body      MachineSetUpdateRequest  true  "Machine set update request"
// @Success      200         {object}  map[string]string
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /machinesets/{id} [put]
func (h *MachineSetWriteHandler) UpdateMachineSet(c *gin.Context) {
	id := c.Param("id")
	var req MachineSetUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify machine set exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineSetType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine set not found"})
		return
	}

	// Use Management service to update machine set
	log.Printf("Updating machine set %s (Management service integration needed)", id)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Machine set update initiated",
		"id": id,
		"note": "Management service integration required for actual update",
	})
}

// DeleteMachineSet godoc
// @Summary      Delete a machine set
// @Description  Delete a machine set from Omni
// @Tags         machinesets
// @Produce      json
// @Param        id   path      string  true  "Machine set ID"
// @Success      202  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machinesets/{id} [delete]
func (h *MachineSetWriteHandler) DeleteMachineSet(c *gin.Context) {
	id := c.Param("id")

	// Verify machine set exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineSetType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine set not found"})
		return
	}

	// Use Management service to delete/teardown machine set
	log.Printf("Deleting machine set %s (Management service integration needed)", id)
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Machine set deletion initiated",
		"id": id,
		"note": "Management service integration required for actual deletion",
	})
}
