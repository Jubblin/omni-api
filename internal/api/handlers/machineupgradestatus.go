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

// MachineUpgradeStatusResponse represents the machine upgrade status information returned by the API
type MachineUpgradeStatusResponse struct {
	ID                  string            `json:"id"`
	Namespace           string            `json:"namespace"`
	SchematicID         string            `json:"schematic_id,omitempty"`
	TalosVersion        string            `json:"talos_version,omitempty"`
	CurrentSchematicID  string            `json:"current_schematic_id,omitempty"`
	CurrentTalosVersion string            `json:"current_talos_version,omitempty"`
	Phase               string            `json:"phase"`
	Status              string            `json:"status,omitempty"`
	Error               string            `json:"error,omitempty"`
	IsMaintenance       bool              `json:"is_maintenance,omitempty"`
	Links               map[string]string `json:"_links,omitempty"`
}

// MachineUpgradeStatusHandler handles machine upgrade status requests
type MachineUpgradeStatusHandler struct {
	state state.State
}

// NewMachineUpgradeStatusHandler creates a new MachineUpgradeStatusHandler
func NewMachineUpgradeStatusHandler(s state.State) *MachineUpgradeStatusHandler {
	return &MachineUpgradeStatusHandler{state: s}
}

// GetMachineUpgradeStatus godoc
// @Summary      Get machine upgrade status
// @Description  Get the upgrade status and progress for a specific machine
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Machine ID"
// @Success      200  {object}  MachineUpgradeStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machines/{id}/upgrade-status [get]
func (h *MachineUpgradeStatusHandler) GetMachineUpgradeStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineUpgradeStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine upgrade status %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mus, ok := res.(*omni.MachineUpgradeStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := mus.TypedSpec().Value
	resp := MachineUpgradeStatusResponse{
		ID:                  mus.Metadata().ID(),
		Namespace:           mus.Metadata().Namespace(),
		SchematicID:         spec.SchematicId,
		TalosVersion:        spec.TalosVersion,
		CurrentSchematicID:  spec.CurrentSchematicId,
		CurrentTalosVersion: spec.CurrentTalosVersion,
		Phase:               spec.Phase.String(),
		Status:              spec.Status,
		Error:               spec.Error,
		IsMaintenance:       spec.IsMaintenance,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/machines/"+id+"/upgrade-status"),
			"machine": buildURL(c, "/api/v1/machines/"+id),
		},
	}

	// Add schematic link if schematic ID is present
	if spec.SchematicId != "" {
		resp.Links["schematic"] = buildURL(c, "/api/v1/schematics/"+spec.SchematicId)
	}
	if spec.CurrentSchematicId != "" {
		resp.Links["current_schematic"] = buildURL(c, "/api/v1/schematics/"+spec.CurrentSchematicId)
	}

	// Try to find cluster ID from labels
	if clusterID, ok := mus.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
