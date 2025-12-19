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

// MachineRequestSetResponse represents the machine request set information
type MachineRequestSetResponse struct {
	ID           string            `json:"id"`
	Namespace    string            `json:"namespace"`
	ProviderID   string            `json:"provider_id,omitempty"`
	MachineCount int32             `json:"machine_count,omitempty"`
	TalosVersion string            `json:"talos_version,omitempty"`
	Extensions   []string          `json:"extensions,omitempty"`
	KernelArgs   []string          `json:"kernel_args,omitempty"`
	ProviderData string            `json:"provider_data,omitempty"`
	GrpcTunnel   string            `json:"grpc_tunnel,omitempty"`
	Links        map[string]string `json:"_links,omitempty"`
}

// MachineRequestSetHandler handles machine request set requests
type MachineRequestSetHandler struct {
	state state.State
}

// NewMachineRequestSetHandler creates a new MachineRequestSetHandler
func NewMachineRequestSetHandler(s state.State) *MachineRequestSetHandler {
	return &MachineRequestSetHandler{state: s}
}

// ListMachineRequestSets godoc
// @Summary      List machine request sets
// @Description  Get a list of all machine request sets
// @Tags         machines
// @Produce      json
// @Success      200  {array}   MachineRequestSetResponse
// @Failure      500  {object}  map[string]string
// @Router       /machine-request-sets [get]
func (h *MachineRequestSetHandler) ListMachineRequestSets(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineRequestSetType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing machine request sets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var requestSets []MachineRequestSetResponse
	for _, item := range items.Items {
		mrs, ok := item.(*omni.MachineRequestSet)
		if !ok {
			log.Printf("Warning: resource is not a machine request set: %T", item)
			continue
		}

		spec := mrs.TypedSpec().Value
		setID := mrs.Metadata().ID()
		resp := MachineRequestSetResponse{
			ID:           setID,
			Namespace:    mrs.Metadata().Namespace(),
			ProviderID:   spec.ProviderId,
			MachineCount: spec.MachineCount,
			TalosVersion: spec.TalosVersion,
			Extensions:   spec.Extensions,
			KernelArgs:   spec.KernelArgs,
			ProviderData: spec.ProviderData,
			GrpcTunnel:   spec.GrpcTunnel.String(),
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/machine-request-sets/"+setID),
			},
		}

		requestSets = append(requestSets, resp)
	}

	c.JSON(http.StatusOK, requestSets)
}

// GetMachineRequestSet godoc
// @Summary      Get a machine request set
// @Description  Get detailed information about a specific machine request set
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Machine Request Set ID"
// @Success      200  {object}  MachineRequestSetResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machine-request-sets/{id} [get]
func (h *MachineRequestSetHandler) GetMachineRequestSet(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineRequestSetType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine request set %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mrs, ok := res.(*omni.MachineRequestSet)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := mrs.TypedSpec().Value
	setID := mrs.Metadata().ID()
	resp := MachineRequestSetResponse{
		ID:           setID,
		Namespace:    mrs.Metadata().Namespace(),
		ProviderID:   spec.ProviderId,
		MachineCount: spec.MachineCount,
		TalosVersion: spec.TalosVersion,
		Extensions:   spec.Extensions,
		KernelArgs:   spec.KernelArgs,
		ProviderData: spec.ProviderData,
		GrpcTunnel:   spec.GrpcTunnel.String(),
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/machine-request-sets/"+setID),
		},
	}

	c.JSON(http.StatusOK, resp)
}
