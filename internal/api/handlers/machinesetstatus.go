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

// MachineSetStatusResponse represents the machine set status information returned by the API
type MachineSetStatusResponse struct {
	ID             string            `json:"id"`
	Namespace      string            `json:"namespace"`
	Phase          string            `json:"phase"`
	Ready          bool              `json:"ready"`
	Error          string            `json:"error,omitempty"`
	ConfigHash     string            `json:"config_hash,omitempty"`
	LockedUpdates  uint32            `json:"locked_updates,omitempty"`
	Machines       struct {
		Total     uint32 `json:"total,omitempty"`
		Healthy   uint32 `json:"healthy,omitempty"`
		Connected uint32 `json:"connected,omitempty"`
		Requested uint32 `json:"requested,omitempty"`
	} `json:"machines,omitempty"`
	Links map[string]string `json:"_links,omitempty"`
}

// MachineSetStatusHandler handles machine set status requests
type MachineSetStatusHandler struct {
	state state.State
}

// NewMachineSetStatusHandler creates a new MachineSetStatusHandler
func NewMachineSetStatusHandler(s state.State) *MachineSetStatusHandler {
	return &MachineSetStatusHandler{state: s}
}

// GetMachineSetStatus godoc
// @Summary      Get machine set status
// @Description  Get status information for a specific machine set
// @Tags         machinesets
// @Produce      json
// @Param        id   path      string  true  "Machine Set ID"
// @Success      200  {object}  MachineSetStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machinesets/{id}/status [get]
func (h *MachineSetStatusHandler) GetMachineSetStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineSetStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine set status %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mss, ok := res.(*omni.MachineSetStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := mss.TypedSpec().Value
	resp := MachineSetStatusResponse{
		ID:            mss.Metadata().ID(),
		Namespace:     mss.Metadata().Namespace(),
		Phase:         spec.Phase.String(),
		Ready:         spec.Ready,
		Error:         spec.Error,
		ConfigHash:    spec.ConfigHash,
		LockedUpdates: spec.LockedUpdates,
		Links: map[string]string{
			"self":      buildURL(c, "/api/v1/machinesets/"+id+"/status"),
			"machineset": buildURL(c, "/api/v1/machinesets/"+id),
		},
	}

	if spec.Machines != nil {
		resp.Machines.Total = spec.Machines.Total
		resp.Machines.Healthy = spec.Machines.Healthy
		resp.Machines.Connected = spec.Machines.Connected
		resp.Machines.Requested = spec.Machines.Requested
	}

	// Try to find cluster ID from labels
	if clusterID, ok := mss.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
