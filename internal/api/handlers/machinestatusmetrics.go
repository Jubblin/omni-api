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

// MachineStatusMetricsResponse represents the machine status metrics information
type MachineStatusMetricsResponse struct {
	ID                      string            `json:"id"`
	Namespace               string            `json:"namespace"`
	RegisteredMachinesCount uint32            `json:"registered_machines_count,omitempty"`
	ConnectedMachinesCount  uint32            `json:"connected_machines_count,omitempty"`
	AllocatedMachinesCount  uint32            `json:"allocated_machines_count,omitempty"`
	PendingMachinesCount    uint32            `json:"pending_machines_count,omitempty"`
	Platforms               map[string]uint32  `json:"platforms,omitempty"`
	SecureBootStatus        map[string]uint32  `json:"secure_boot_status,omitempty"`
	UkiStatus               map[string]uint32  `json:"uki_status,omitempty"`
	Links                   map[string]string  `json:"_links,omitempty"`
}

// MachineStatusMetricsHandler handles machine status metrics requests
type MachineStatusMetricsHandler struct {
	state state.State
}

// NewMachineStatusMetricsHandler creates a new MachineStatusMetricsHandler
func NewMachineStatusMetricsHandler(s state.State) *MachineStatusMetricsHandler {
	return &MachineStatusMetricsHandler{state: s}
}

// GetMachineStatusMetrics godoc
// @Summary      Get machine status metrics
// @Description  Get aggregated metrics for all machines in Omni
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Machine ID (used for link generation, metrics are global)"
// @Success      200  {object}  MachineStatusMetricsResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machines/{id}/metrics [get]
func (h *MachineStatusMetricsHandler) GetMachineStatusMetrics(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	// MachineStatusMetrics is a singleton resource, typically with ID "default" or empty
	// Try common IDs
	possibleIDs := []string{"default", "omni", ""}
	
	var metricsRes resource.Resource
	var err error
	
	for _, metricsID := range possibleIDs {
		md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineStatusMetricsType, metricsID, resource.VersionUndefined)
		metricsRes, err = st.Get(c.Request.Context(), md)
		if err == nil {
			break
		}
	}
	
	// If not found, try listing to find the actual ID
	if err != nil {
		md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineStatusMetricsType, "", resource.VersionUndefined)
		list, listErr := st.List(c.Request.Context(), md)
		if listErr == nil && len(list.Items) > 0 {
			metricsRes = list.Items[0]
			err = nil
		}
	}
	
	if err != nil {
		log.Printf("Error getting machine status metrics: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "machine status metrics not found"})
		return
	}

	msm, ok := metricsRes.(*omni.MachineStatusMetrics)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := msm.TypedSpec().Value
	resp := MachineStatusMetricsResponse{
		ID:                      msm.Metadata().ID(),
		Namespace:               msm.Metadata().Namespace(),
		RegisteredMachinesCount: spec.RegisteredMachinesCount,
		ConnectedMachinesCount:  spec.ConnectedMachinesCount,
		AllocatedMachinesCount:  spec.AllocatedMachinesCount,
		PendingMachinesCount:    spec.PendingMachinesCount,
		Platforms:               spec.Platforms,
		SecureBootStatus:        spec.SecureBootStatus,
		UkiStatus:               spec.UkiStatus,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/machines/"+id+"/metrics"),
			"machine": buildURL(c, "/api/v1/machines/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
