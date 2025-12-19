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

// KubernetesUpgradeStatusResponse represents the Kubernetes upgrade status information returned by the API
type KubernetesUpgradeStatusResponse struct {
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

// KubernetesUpgradeHandler handles Kubernetes upgrade status requests
type KubernetesUpgradeHandler struct {
	state state.State
}

// NewKubernetesUpgradeHandler creates a new KubernetesUpgradeHandler
func NewKubernetesUpgradeHandler(s state.State) *KubernetesUpgradeHandler {
	return &KubernetesUpgradeHandler{state: s}
}

// GetKubernetesUpgradeStatus godoc
// @Summary      Get Kubernetes upgrade status
// @Description  Get the status of Kubernetes version upgrades for a cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  KubernetesUpgradeStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/kubernetes-upgrade [get]
func (h *KubernetesUpgradeHandler) GetKubernetesUpgradeStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.KubernetesUpgradeStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting Kubernetes upgrade status %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kus, ok := res.(*omni.KubernetesUpgradeStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := kus.TypedSpec().Value
	resp := KubernetesUpgradeStatusResponse{
		ID:                    kus.Metadata().ID(),
		Namespace:             kus.Metadata().Namespace(),
		Phase:                 spec.Phase.String(),
		Error:                 spec.Error,
		Step:                  spec.Step,
		Status:                spec.Status,
		LastUpgradeVersion:    spec.LastUpgradeVersion,
		CurrentUpgradeVersion: spec.CurrentUpgradeVersion,
		UpgradeVersions:       spec.UpgradeVersions,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/clusters/"+id+"/kubernetes-upgrade"),
			"cluster": buildURL(c, "/api/v1/clusters/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
