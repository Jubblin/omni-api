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

// SchematicResponse represents the schematic information returned by the API
type SchematicResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Links     map[string]string `json:"_links,omitempty"`
}

// SchematicHandler handles schematic requests
type SchematicHandler struct {
	state state.State
}

// NewSchematicHandler creates a new SchematicHandler
func NewSchematicHandler(s state.State) *SchematicHandler {
	return &SchematicHandler{state: s}
}

// ListSchematics godoc
// @Summary      List all schematics
// @Description  Get a list of all Talos image schematics in Omni
// @Tags         schematics
// @Produce      json
// @Success      200  {array}   SchematicResponse
// @Failure      500  {object}  map[string]string
// @Router       /schematics [get]
func (h *SchematicHandler) ListSchematics(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.SchematicType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing schematics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var schematics []SchematicResponse
	for _, item := range items.Items {
		s, ok := item.(*omni.Schematic)
		if !ok {
			log.Printf("Warning: resource is not a schematic: %T", item)
			continue
		}

		schematicID := s.Metadata().ID()
		resp := SchematicResponse{
			ID:        schematicID,
			Namespace: s.Metadata().Namespace(),
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/schematics/"+schematicID),
			},
		}

		schematics = append(schematics, resp)
	}

	c.JSON(http.StatusOK, schematics)
}

// GetSchematic godoc
// @Summary      Get a single schematic
// @Description  Get detailed information about a specific schematic
// @Tags         schematics
// @Produce      json
// @Param        id   path      string  true  "Schematic ID"
// @Success      200  {object}  SchematicResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /schematics/{id} [get]
func (h *SchematicHandler) GetSchematic(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.SchematicType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting schematic %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s, ok := res.(*omni.Schematic)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	schematicID := s.Metadata().ID()
	resp := SchematicResponse{
		ID:        schematicID,
		Namespace: s.Metadata().Namespace(),
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/schematics/"+schematicID),
		},
	}

	c.JSON(http.StatusOK, resp)
}
