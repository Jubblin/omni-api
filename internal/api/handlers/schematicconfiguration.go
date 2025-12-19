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

// SchematicConfigurationResponse represents the schematic configuration information
type SchematicConfigurationResponse struct {
	ID           string            `json:"id"`
	Namespace    string            `json:"namespace"`
	SchematicID  string            `json:"schematic_id,omitempty"`
	TalosVersion string            `json:"talos_version,omitempty"`
	KernelArgs   []string          `json:"kernel_args,omitempty"`
	Links        map[string]string `json:"_links,omitempty"`
}

// SchematicConfigurationHandler handles schematic configuration requests
type SchematicConfigurationHandler struct {
	state state.State
}

// NewSchematicConfigurationHandler creates a new SchematicConfigurationHandler
func NewSchematicConfigurationHandler(s state.State) *SchematicConfigurationHandler {
	return &SchematicConfigurationHandler{state: s}
}

// ListSchematicConfigurations godoc
// @Summary      List schematic configurations
// @Description  Get a list of all schematic configurations
// @Tags         schematics
// @Produce      json
// @Success      200  {array}   SchematicConfigurationResponse
// @Failure      500  {object}  map[string]string
// @Router       /schematic-configurations [get]
func (h *SchematicConfigurationHandler) ListSchematicConfigurations(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.SchematicConfigurationType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing schematic configurations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var configs []SchematicConfigurationResponse
	for _, item := range items.Items {
		sc, ok := item.(*omni.SchematicConfiguration)
		if !ok {
			log.Printf("Warning: resource is not a schematic configuration: %T", item)
			continue
		}

		spec := sc.TypedSpec().Value
		configID := sc.Metadata().ID()
		resp := SchematicConfigurationResponse{
			ID:           configID,
			Namespace:    sc.Metadata().Namespace(),
			SchematicID:  spec.SchematicId,
			TalosVersion: spec.TalosVersion,
			KernelArgs:   spec.KernelArgs,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/schematic-configurations/"+configID),
			},
		}

		if spec.SchematicId != "" {
			resp.Links["schematic"] = buildURL(c, "/api/v1/schematics/"+spec.SchematicId)
		}

		configs = append(configs, resp)
	}

	c.JSON(http.StatusOK, configs)
}

// GetSchematicConfiguration godoc
// @Summary      Get a schematic configuration
// @Description  Get detailed information about a specific schematic configuration
// @Tags         schematics
// @Produce      json
// @Param        id   path      string  true  "Schematic Configuration ID"
// @Success      200  {object}  SchematicConfigurationResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /schematic-configurations/{id} [get]
func (h *SchematicConfigurationHandler) GetSchematicConfiguration(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.SchematicConfigurationType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting schematic configuration %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sc, ok := res.(*omni.SchematicConfiguration)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := sc.TypedSpec().Value
	configID := sc.Metadata().ID()
	resp := SchematicConfigurationResponse{
		ID:           configID,
		Namespace:    sc.Metadata().Namespace(),
		SchematicID:  spec.SchematicId,
		TalosVersion: spec.TalosVersion,
		KernelArgs:   spec.KernelArgs,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/schematic-configurations/"+configID),
		},
	}

	if spec.SchematicId != "" {
		resp.Links["schematic"] = buildURL(c, "/api/v1/schematics/"+spec.SchematicId)
	}

	c.JSON(http.StatusOK, resp)
}
