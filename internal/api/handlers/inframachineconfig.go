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

// InfraMachineConfigResponse represents the infrastructure machine config information
type InfraMachineConfigResponse struct {
	ID                string            `json:"id"`
	Namespace         string            `json:"namespace"`
	PowerState        string            `json:"power_state,omitempty"`
	AcceptanceStatus  string            `json:"acceptance_status,omitempty"`
	ExtraKernelArgs   string            `json:"extra_kernel_args,omitempty"`
	RequestedRebootID string            `json:"requested_reboot_id,omitempty"`
	Cordoned          bool              `json:"cordoned,omitempty"`
	Links             map[string]string `json:"_links,omitempty"`
}

// InfraMachineConfigHandler handles infrastructure machine config requests
type InfraMachineConfigHandler struct {
	state state.State
}

// NewInfraMachineConfigHandler creates a new InfraMachineConfigHandler
func NewInfraMachineConfigHandler(s state.State) *InfraMachineConfigHandler {
	return &InfraMachineConfigHandler{state: s}
}

// ListInfraMachineConfigs godoc
// @Summary      List all infrastructure machine configs
// @Description  Get a list of all infrastructure machine configs in Omni
// @Tags         inframachineconfigs
// @Produce      json
// @Param        machine   query     string  false  "Filter by machine ID"
// @Success      200  {array}   InfraMachineConfigResponse
// @Failure      500  {object}  map[string]string
// @Router       /infra-machine-configs [get]
func (h *InfraMachineConfigHandler) ListInfraMachineConfigs(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.InfraMachineConfigType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing infrastructure machine configs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	machineFilter := c.Query("machine")

	var configs []InfraMachineConfigResponse
	for _, item := range items.Items {
		imc, ok := item.(*omni.InfraMachineConfig)
		if !ok {
			log.Printf("Warning: resource is not an infrastructure machine config: %T", item)
			continue
		}

		// Filter by machine if specified
		if machineFilter != "" {
			if imc.Metadata().ID() != machineFilter {
				continue
			}
		}

		configID := imc.Metadata().ID()
		spec := imc.TypedSpec().Value
		resp := InfraMachineConfigResponse{
			ID:                configID,
			Namespace:         imc.Metadata().Namespace(),
			PowerState:        spec.PowerState.String(),
			AcceptanceStatus:  spec.AcceptanceStatus.String(),
			ExtraKernelArgs:   spec.ExtraKernelArgs,
			RequestedRebootID: spec.RequestedRebootId,
			Cordoned:          spec.Cordoned,
			Links: map[string]string{
				"self":    buildURL(c, "/api/v1/infra-machine-configs/"+configID),
				"machine": buildURL(c, "/api/v1/machines/"+configID),
			},
		}

		configs = append(configs, resp)
	}

	c.JSON(http.StatusOK, configs)
}

// GetInfraMachineConfig godoc
// @Summary      Get a single infrastructure machine config
// @Description  Get detailed information about a specific infrastructure machine config
// @Tags         inframachineconfigs
// @Produce      json
// @Param        id   path      string  true  "Infrastructure Machine Config ID"
// @Success      200  {object}  InfraMachineConfigResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /infra-machine-configs/{id} [get]
func (h *InfraMachineConfigHandler) GetInfraMachineConfig(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.InfraMachineConfigType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting infrastructure machine config %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "infrastructure machine config not found"})
		return
	}

	imc, ok := res.(*omni.InfraMachineConfig)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := imc.TypedSpec().Value
	resp := InfraMachineConfigResponse{
		ID:                imc.Metadata().ID(),
		Namespace:         imc.Metadata().Namespace(),
		PowerState:        spec.PowerState.String(),
		AcceptanceStatus:  spec.AcceptanceStatus.String(),
		ExtraKernelArgs:   spec.ExtraKernelArgs,
		RequestedRebootID: spec.RequestedRebootId,
		Cordoned:          spec.Cordoned,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/infra-machine-configs/"+id),
			"machine": buildURL(c, "/api/v1/machines/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
