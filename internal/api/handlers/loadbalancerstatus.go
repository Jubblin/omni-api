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

// LoadBalancerStatusResponse represents the load balancer status information
type LoadBalancerStatusResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Healthy   bool              `json:"healthy,omitempty"`
	Stopped   bool              `json:"stopped,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// LoadBalancerStatusHandler handles load balancer status requests
type LoadBalancerStatusHandler struct {
	state state.State
}

// NewLoadBalancerStatusHandler creates a new LoadBalancerStatusHandler
func NewLoadBalancerStatusHandler(s state.State) *LoadBalancerStatusHandler {
	return &LoadBalancerStatusHandler{state: s}
}

// GetLoadBalancerStatus godoc
// @Summary      Get load balancer status
// @Description  Get status of a load balancer
// @Tags         infrastructure
// @Produce      json
// @Param        id   path      string  true  "Load Balancer ID"
// @Success      200  {object}  LoadBalancerStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /loadbalancers/{id}/status [get]
func (h *LoadBalancerStatusHandler) GetLoadBalancerStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.LoadBalancerStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting load balancer status %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "load balancer status not found"})
		return
	}

	lbs, ok := res.(*omni.LoadBalancerStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := lbs.TypedSpec().Value
	resp := LoadBalancerStatusResponse{
		ID:        lbs.Metadata().ID(),
		Namespace: lbs.Metadata().Namespace(),
		Healthy:   spec.Healthy,
		Stopped:   spec.Stopped,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/loadbalancers/"+id+"/status"),
		},
	}

	c.JSON(http.StatusOK, resp)
}
