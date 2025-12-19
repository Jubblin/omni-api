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

// TalosUpgradeStatusResponse represents the Talos upgrade status information returned by the API
type TalosUpgradeStatusResponse struct {
	ID                    string   `json:"id"`
	Namespace             string   `json:"namespace"`
	Phase                 string   `json:"phase"`
	Error                 string   `json:"error,omitempty"`
	Step                  string   `json:"step,omitempty"`
	Status                string   `json:"status,omitempty"`
	LastUpgradeVersion    string   `json:"last_upgrade_version,omitempty"`
	CurrentUpgradeVersion string   `json:"current_upgrade_version,omitempty"`
	UpgradeVersions       []string `json:"upgrade_versions,omitempty"`
	Links                 map[string]string `json:"_links,omitempty"`
}

// TalosUpgradeHandler handles Talos upgrade status requests
type TalosUpgradeHandler struct {
	state state.State
}

// NewTalosUpgradeHandler creates a new TalosUpgradeHandler
func NewTalosUpgradeHandler(s state.State) *TalosUpgradeHandler {
	return &TalosUpgradeHandler{state: s}
}

// GetTalosUpgradeStatus godoc
// @Summary      Get Talos upgrade status
// @Description  Get the status of Talos OS upgrades for a cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  TalosUpgradeStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/talos-upgrade [get]
func (h *TalosUpgradeHandler) GetTalosUpgradeStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.TalosUpgradeStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting Talos upgrade status %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tus, ok := res.(*omni.TalosUpgradeStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := tus.TypedSpec().Value
	resp := TalosUpgradeStatusResponse{
		ID:                    tus.Metadata().ID(),
		Namespace:             tus.Metadata().Namespace(),
		Phase:                 spec.Phase.String(),
		Error:                 spec.Error,
		Step:                  spec.Step,
		Status:                spec.Status,
		LastUpgradeVersion:    spec.LastUpgradeVersion,
		CurrentUpgradeVersion: spec.CurrentUpgradeVersion,
		UpgradeVersions:       spec.UpgradeVersions,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/clusters/"+id+"/talos-upgrade"),
			"cluster": buildURL(c, "/api/v1/clusters/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
