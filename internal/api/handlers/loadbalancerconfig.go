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

// LoadBalancerConfigResponse represents the load balancer config information
type LoadBalancerConfigResponse struct {
	ID                string            `json:"id"`
	Namespace         string            `json:"namespace"`
	BindPort          string            `json:"bind_port,omitempty"`
	SiderolinkEndpoint string           `json:"siderolink_endpoint,omitempty"`
	Endpoints         []string          `json:"endpoints,omitempty"`
	Links             map[string]string `json:"_links,omitempty"`
}

// LoadBalancerConfigHandler handles load balancer config requests
type LoadBalancerConfigHandler struct {
	state state.State
}

// NewLoadBalancerConfigHandler creates a new LoadBalancerConfigHandler
func NewLoadBalancerConfigHandler(s state.State) *LoadBalancerConfigHandler {
	return &LoadBalancerConfigHandler{state: s}
}

// ListLoadBalancerConfigs godoc
// @Summary      List load balancer configs
// @Description  Get a list of all load balancer configurations
// @Tags         infrastructure
// @Produce      json
// @Success      200  {array}   LoadBalancerConfigResponse
// @Failure      500  {object}  map[string]string
// @Router       /loadbalancer-configs [get]
func (h *LoadBalancerConfigHandler) ListLoadBalancerConfigs(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.LoadBalancerConfigType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing load balancer configs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var configs []LoadBalancerConfigResponse
	for _, item := range items.Items {
		lb, ok := item.(*omni.LoadBalancerConfig)
		if !ok {
			log.Printf("Warning: resource is not a load balancer config: %T", item)
			continue
		}

		spec := lb.TypedSpec().Value
		configID := lb.Metadata().ID()
		resp := LoadBalancerConfigResponse{
			ID:                 configID,
			Namespace:          lb.Metadata().Namespace(),
			BindPort:           spec.BindPort,
			SiderolinkEndpoint: spec.SiderolinkEndpoint,
			Endpoints:          spec.Endpoints,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/loadbalancer-configs/"+configID),
			},
		}

		configs = append(configs, resp)
	}

	c.JSON(http.StatusOK, configs)
}

// GetLoadBalancerConfig godoc
// @Summary      Get a load balancer config
// @Description  Get detailed information about a specific load balancer configuration
// @Tags         infrastructure
// @Produce      json
// @Param        id   path      string  true  "Load Balancer Config ID"
// @Success      200  {object}  LoadBalancerConfigResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /loadbalancer-configs/{id} [get]
func (h *LoadBalancerConfigHandler) GetLoadBalancerConfig(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.LoadBalancerConfigType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting load balancer config %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lb, ok := res.(*omni.LoadBalancerConfig)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := lb.TypedSpec().Value
	configID := lb.Metadata().ID()
	resp := LoadBalancerConfigResponse{
		ID:                 configID,
		Namespace:          lb.Metadata().Namespace(),
		BindPort:           spec.BindPort,
		SiderolinkEndpoint: spec.SiderolinkEndpoint,
		Endpoints:          spec.Endpoints,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/loadbalancer-configs/"+configID),
		},
	}

	c.JSON(http.StatusOK, resp)
}
