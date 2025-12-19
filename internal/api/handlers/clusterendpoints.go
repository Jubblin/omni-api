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

// ClusterEndpointResponse represents the cluster endpoint information returned by the API
type ClusterEndpointResponse struct {
	ID                  string            `json:"id"`
	Namespace           string            `json:"namespace"`
	ManagementAddresses []string          `json:"management_addresses"`
	Links               map[string]string `json:"_links,omitempty"`
}

// ClusterEndpointHandler handles cluster endpoint requests
type ClusterEndpointHandler struct {
	state state.State
}

// NewClusterEndpointHandler creates a new ClusterEndpointHandler
func NewClusterEndpointHandler(s state.State) *ClusterEndpointHandler {
	return &ClusterEndpointHandler{state: s}
}

// GetClusterEndpoints godoc
// @Summary      Get cluster endpoints
// @Description  Get management endpoints for a specific cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  ClusterEndpointResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/endpoints [get]
func (h *ClusterEndpointHandler) GetClusterEndpoints(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterEndpointType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster endpoints %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ce, ok := res.(*omni.ClusterEndpoint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ce.TypedSpec().Value
	resp := ClusterEndpointResponse{
		ID:                  ce.Metadata().ID(),
		Namespace:           ce.Metadata().Namespace(),
		ManagementAddresses: spec.ManagementAddresses,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/clusters/"+id+"/endpoints"),
			"cluster": buildURL(c, "/api/v1/clusters/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
