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

// MachineSetDestroyStatusResponse represents the machine set destroy status information
type MachineSetDestroyStatusResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Phase     string            `json:"phase,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// MachineSetDestroyStatusHandler handles machine set destroy status requests
type MachineSetDestroyStatusHandler struct {
	state state.State
}

// NewMachineSetDestroyStatusHandler creates a new MachineSetDestroyStatusHandler
func NewMachineSetDestroyStatusHandler(s state.State) *MachineSetDestroyStatusHandler {
	return &MachineSetDestroyStatusHandler{state: s}
}

// GetMachineSetDestroyStatus godoc
// @Summary      Get machine set destroy status
// @Description  Get the status of machine set destruction operation
// @Tags         machinesets
// @Produce      json
// @Param        id   path      string  true  "Machine Set ID"
// @Success      200  {object}  MachineSetDestroyStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machinesets/{id}/destroy-status [get]
func (h *MachineSetDestroyStatusHandler) GetMachineSetDestroyStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineSetDestroyStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine set destroy status %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "machine set destroy status not found"})
		return
	}

	msds, ok := res.(*omni.MachineSetDestroyStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := msds.TypedSpec().Value
	resp := MachineSetDestroyStatusResponse{
		ID:        msds.Metadata().ID(),
		Namespace: msds.Metadata().Namespace(),
		Phase:     spec.Phase,
		Links: map[string]string{
			"self":      buildURL(c, "/api/v1/machinesets/"+id+"/destroy-status"),
			"machineset": buildURL(c, "/api/v1/machinesets/"+id),
		},
	}

	// Try to find cluster ID from labels
	if clusterID, ok := msds.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
