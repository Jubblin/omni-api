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

// ClusterDestroyStatusResponse represents the cluster destroy status information
type ClusterDestroyStatusResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Phase     string            `json:"phase,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// ClusterDestroyStatusHandler handles cluster destroy status requests
type ClusterDestroyStatusHandler struct {
	state state.State
}

// NewClusterDestroyStatusHandler creates a new ClusterDestroyStatusHandler
func NewClusterDestroyStatusHandler(s state.State) *ClusterDestroyStatusHandler {
	return &ClusterDestroyStatusHandler{state: s}
}

// GetClusterDestroyStatus godoc
// @Summary      Get cluster destroy status
// @Description  Get the status of cluster destruction operation
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  ClusterDestroyStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/destroy-status [get]
func (h *ClusterDestroyStatusHandler) GetClusterDestroyStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterDestroyStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster destroy status %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster destroy status not found"})
		return
	}

	cds, ok := res.(*omni.ClusterDestroyStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := cds.TypedSpec().Value
	resp := ClusterDestroyStatusResponse{
		ID:        cds.Metadata().ID(),
		Namespace: cds.Metadata().Namespace(),
		Phase:     spec.Phase,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/clusters/"+id+"/destroy-status"),
			"cluster": buildURL(c, "/api/v1/clusters/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
