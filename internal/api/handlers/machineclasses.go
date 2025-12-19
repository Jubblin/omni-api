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

// MachineClassResponse represents the machine class information returned by the API
type MachineClassResponse struct {
	ID            string            `json:"id"`
	Namespace     string            `json:"namespace"`
	MatchLabels   []string          `json:"match_labels,omitempty"`
	AutoProvision bool              `json:"auto_provision,omitempty"`
	Links         map[string]string `json:"_links,omitempty"`
}

// MachineClassHandler handles machine class requests
type MachineClassHandler struct {
	state state.State
}

// NewMachineClassHandler creates a new MachineClassHandler
func NewMachineClassHandler(s state.State) *MachineClassHandler {
	return &MachineClassHandler{state: s}
}

// ListMachineClasses godoc
// @Summary      List all machine classes
// @Description  Get a list of all machine classes in Omni
// @Tags         machineclasses
// @Produce      json
// @Success      200  {array}   MachineClassResponse
// @Failure      500  {object}  map[string]string
// @Router       /machineclasses [get]
func (h *MachineClassHandler) ListMachineClasses(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineClassType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing machine classes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var machineClasses []MachineClassResponse
	for _, item := range items.Items {
		mc, ok := item.(*omni.MachineClass)
		if !ok {
			log.Printf("Warning: resource is not a machine class: %T", item)
			continue
		}

		spec := mc.TypedSpec().Value
		classID := mc.Metadata().ID()
		resp := MachineClassResponse{
			ID:          classID,
			Namespace:   mc.Metadata().Namespace(),
			MatchLabels: spec.MatchLabels,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/machineclasses/"+classID),
			},
		}

		if spec.AutoProvision != nil {
			resp.AutoProvision = true
		}

		machineClasses = append(machineClasses, resp)
	}

	c.JSON(http.StatusOK, machineClasses)
}

// GetMachineClass godoc
// @Summary      Get a single machine class
// @Description  Get detailed information about a specific machine class
// @Tags         machineclasses
// @Produce      json
// @Param        id   path      string  true  "Machine Class ID"
// @Success      200  {object}  MachineClassResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machineclasses/{id} [get]
func (h *MachineClassHandler) GetMachineClass(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineClassType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine class %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mc, ok := res.(*omni.MachineClass)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := mc.TypedSpec().Value
	classID := mc.Metadata().ID()
	resp := MachineClassResponse{
		ID:          classID,
		Namespace:   mc.Metadata().Namespace(),
		MatchLabels: spec.MatchLabels,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/machineclasses/"+classID),
		},
	}

	if spec.AutoProvision != nil {
		resp.AutoProvision = true
	}

	c.JSON(http.StatusOK, resp)
}
