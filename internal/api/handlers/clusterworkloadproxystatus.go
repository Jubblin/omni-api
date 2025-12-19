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

// ClusterWorkloadProxyStatusResponse represents the cluster workload proxy status information
type ClusterWorkloadProxyStatusResponse struct {
	ID                 string            `json:"id"`
	Namespace          string            `json:"namespace"`
	NumExposedServices uint32           `json:"num_exposed_services,omitempty"`
	Links              map[string]string `json:"_links,omitempty"`
}

// ClusterWorkloadProxyStatusHandler handles cluster workload proxy status requests
type ClusterWorkloadProxyStatusHandler struct {
	state state.State
}

// NewClusterWorkloadProxyStatusHandler creates a new ClusterWorkloadProxyStatusHandler
func NewClusterWorkloadProxyStatusHandler(s state.State) *ClusterWorkloadProxyStatusHandler {
	return &ClusterWorkloadProxyStatusHandler{state: s}
}

// GetClusterWorkloadProxyStatus godoc
// @Summary      Get cluster workload proxy status
// @Description  Get the status of workload proxy for a cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  ClusterWorkloadProxyStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/workload-proxy-status [get]
func (h *ClusterWorkloadProxyStatusHandler) GetClusterWorkloadProxyStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterWorkloadProxyStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster workload proxy status %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster workload proxy status not found"})
		return
	}

	cwps, ok := res.(*omni.ClusterWorkloadProxyStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := cwps.TypedSpec().Value
	resp := ClusterWorkloadProxyStatusResponse{
		ID:                 cwps.Metadata().ID(),
		Namespace:          cwps.Metadata().Namespace(),
		NumExposedServices: spec.NumExposedServices,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/clusters/"+id+"/workload-proxy-status"),
			"cluster": buildURL(c, "/api/v1/clusters/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
