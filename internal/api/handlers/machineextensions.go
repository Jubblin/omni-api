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

// MachineExtensionsResponse represents the machine extensions information returned by the API
type MachineExtensionsResponse struct {
	ID         string            `json:"id"`
	Namespace  string            `json:"namespace"`
	Extensions []string          `json:"extensions,omitempty"`
	Links      map[string]string `json:"_links,omitempty"`
}

// MachineExtensionsHandler handles machine extensions requests
type MachineExtensionsHandler struct {
	state state.State
}

// NewMachineExtensionsHandler creates a new MachineExtensionsHandler
func NewMachineExtensionsHandler(s state.State) *MachineExtensionsHandler {
	return &MachineExtensionsHandler{state: s}
}

// GetMachineExtensions godoc
// @Summary      Get machine extensions
// @Description  Get the list of Talos extensions installed on a specific machine
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Machine ID"
// @Success      200  {object}  MachineExtensionsResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machines/{id}/extensions [get]
func (h *MachineExtensionsHandler) GetMachineExtensions(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineExtensionsType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine extensions %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	me, ok := res.(*omni.MachineExtensions)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := me.TypedSpec().Value
	resp := MachineExtensionsResponse{
		ID:         me.Metadata().ID(),
		Namespace:  me.Metadata().Namespace(),
		Extensions: spec.Extensions,
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/machines/"+id+"/extensions"),
			"machine": buildURL(c, "/api/v1/machines/"+id),
		},
	}

	// Try to find cluster ID from labels
	if clusterID, ok := me.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
