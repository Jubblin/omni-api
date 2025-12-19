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

// ClusterDiagnosticsNode represents a node with diagnostics
type ClusterDiagnosticsNode struct {
	ID              string `json:"id"`
	NumDiagnostics  uint32 `json:"num_diagnostics"`
}

// ClusterDiagnosticsResponse represents the cluster diagnostics information
type ClusterDiagnosticsResponse struct {
	ID        string                  `json:"id"`
	Namespace string                  `json:"namespace"`
	Nodes     []ClusterDiagnosticsNode `json:"nodes,omitempty"`
	Links     map[string]string       `json:"_links,omitempty"`
}

// ClusterDiagnosticsHandler handles cluster diagnostics requests
type ClusterDiagnosticsHandler struct {
	state state.State
}

// NewClusterDiagnosticsHandler creates a new ClusterDiagnosticsHandler
func NewClusterDiagnosticsHandler(s state.State) *ClusterDiagnosticsHandler {
	return &ClusterDiagnosticsHandler{state: s}
}

// GetClusterDiagnostics godoc
// @Summary      Get cluster diagnostics
// @Description  Get diagnostic information for a cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  ClusterDiagnosticsResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/diagnostics [get]
func (h *ClusterDiagnosticsHandler) GetClusterDiagnostics(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterDiagnosticsType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster diagnostics %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster diagnostics not found"})
		return
	}

	cd, ok := res.(*omni.ClusterDiagnostics)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := cd.TypedSpec().Value
	resp := ClusterDiagnosticsResponse{
		ID:        cd.Metadata().ID(),
		Namespace: cd.Metadata().Namespace(),
		Nodes:     make([]ClusterDiagnosticsNode, 0, len(spec.Nodes)),
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/clusters/"+id+"/diagnostics"),
			"cluster": buildURL(c, "/api/v1/clusters/"+id),
		},
	}

	for _, node := range spec.Nodes {
		resp.Nodes = append(resp.Nodes, ClusterDiagnosticsNode{
			ID:             node.Id,
			NumDiagnostics: node.NumDiagnostics,
		})
	}

	c.JSON(http.StatusOK, resp)
}
