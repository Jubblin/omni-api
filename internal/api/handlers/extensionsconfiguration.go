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

// ExtensionsConfigurationResponse represents the extensions configuration information
type ExtensionsConfigurationResponse struct {
	ID         string            `json:"id"`
	Namespace  string            `json:"namespace"`
	Extensions []string          `json:"extensions,omitempty"`
	Links      map[string]string `json:"_links,omitempty"`
}

// ExtensionsConfigurationHandler handles extensions configuration requests
type ExtensionsConfigurationHandler struct {
	state state.State
}

// NewExtensionsConfigurationHandler creates a new ExtensionsConfigurationHandler
func NewExtensionsConfigurationHandler(s state.State) *ExtensionsConfigurationHandler {
	return &ExtensionsConfigurationHandler{state: s}
}

// ListExtensionsConfigurations godoc
// @Summary      List extensions configurations
// @Description  Get a list of all extensions configurations
// @Tags         machines
// @Produce      json
// @Success      200  {array}   ExtensionsConfigurationResponse
// @Failure      500  {object}  map[string]string
// @Router       /extensions-configurations [get]
func (h *ExtensionsConfigurationHandler) ListExtensionsConfigurations(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ExtensionsConfigurationType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing extensions configurations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var configs []ExtensionsConfigurationResponse
	for _, item := range items.Items {
		ec, ok := item.(*omni.ExtensionsConfiguration)
		if !ok {
			log.Printf("Warning: resource is not an extensions configuration: %T", item)
			continue
		}

		spec := ec.TypedSpec().Value
		configID := ec.Metadata().ID()
		resp := ExtensionsConfigurationResponse{
			ID:         configID,
			Namespace:  ec.Metadata().Namespace(),
			Extensions: spec.Extensions,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/extensions-configurations/"+configID),
			},
		}

		configs = append(configs, resp)
	}

	c.JSON(http.StatusOK, configs)
}

// GetExtensionsConfiguration godoc
// @Summary      Get an extensions configuration
// @Description  Get detailed information about a specific extensions configuration
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Extensions Configuration ID"
// @Success      200  {object}  ExtensionsConfigurationResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /extensions-configurations/{id} [get]
func (h *ExtensionsConfigurationHandler) GetExtensionsConfiguration(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ExtensionsConfigurationType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting extensions configuration %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ec, ok := res.(*omni.ExtensionsConfiguration)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ec.TypedSpec().Value
	configID := ec.Metadata().ID()
	resp := ExtensionsConfigurationResponse{
		ID:         configID,
		Namespace:  ec.Metadata().Namespace(),
		Extensions: spec.Extensions,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/extensions-configurations/"+configID),
		},
	}

	c.JSON(http.StatusOK, resp)
}
