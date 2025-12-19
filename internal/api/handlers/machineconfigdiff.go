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

// MachineConfigDiffResponse represents the machine config diff information
type MachineConfigDiffResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Diff      string            `json:"diff,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// MachineConfigDiffHandler handles machine config diff requests
type MachineConfigDiffHandler struct {
	state state.State
}

// NewMachineConfigDiffHandler creates a new MachineConfigDiffHandler
func NewMachineConfigDiffHandler(s state.State) *MachineConfigDiffHandler {
	return &MachineConfigDiffHandler{state: s}
}

// GetMachineConfigDiff godoc
// @Summary      Get machine config diff
// @Description  Get the configuration difference for a machine
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Machine ID"
// @Success      200  {object}  MachineConfigDiffResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machines/{id}/config-diff [get]
func (h *MachineConfigDiffHandler) GetMachineConfigDiff(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineConfigDiffType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine config diff %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "machine config diff not found"})
		return
	}

	mcd, ok := res.(*omni.MachineConfigDiff)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := mcd.TypedSpec().Value
	resp := MachineConfigDiffResponse{
		ID:        mcd.Metadata().ID(),
		Namespace: mcd.Metadata().Namespace(),
		Diff:      spec.Diff,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/machines/"+id+"/config-diff"),
			"machine": buildURL(c, "/api/v1/machines/"+id),
		},
	}

	// Try to find cluster ID from labels
	if clusterID, ok := mcd.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
