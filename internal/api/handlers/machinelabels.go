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

// MachineLabelsResponse represents the machine labels information returned by the API
type MachineLabelsResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// MachineLabelsHandler handles machine labels requests
type MachineLabelsHandler struct {
	state state.State
}

// NewMachineLabelsHandler creates a new MachineLabelsHandler
func NewMachineLabelsHandler(s state.State) *MachineLabelsHandler {
	return &MachineLabelsHandler{state: s}
}

// GetMachineLabels godoc
// @Summary      Get machine labels
// @Description  Get user-defined labels for a specific machine
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Machine ID"
// @Success      200  {object}  MachineLabelsResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machines/{id}/labels [get]
func (h *MachineLabelsHandler) GetMachineLabels(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineLabelsType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine labels %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ml, ok := res.(*omni.MachineLabels)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	resp := MachineLabelsResponse{
		ID:        ml.Metadata().ID(),
		Namespace: ml.Metadata().Namespace(),
		Labels:    make(map[string]string),
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/machines/"+id+"/labels"),
			"machine": buildURL(c, "/api/v1/machines/"+id),
		},
	}

	// Collect all labels from metadata (MachineLabelsSpec is empty, labels are in metadata)
	for key, value := range ml.Metadata().Labels().Raw() {
		resp.Labels[key] = value
	}

	// Try to find cluster ID from labels
	if clusterID, ok := ml.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
