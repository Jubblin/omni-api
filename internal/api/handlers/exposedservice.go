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

// ExposedServiceResponse represents the exposed service information
type ExposedServiceResponse struct {
	ID              string            `json:"id"`
	Namespace       string            `json:"namespace"`
	Port            uint32            `json:"port,omitempty"`
	Label           string            `json:"label,omitempty"`
	IconBase64      string            `json:"icon_base64,omitempty"`
	URL             string            `json:"url,omitempty"`
	Error           string            `json:"error,omitempty"`
	HasExplicitAlias bool             `json:"has_explicit_alias,omitempty"`
	Links           map[string]string `json:"_links,omitempty"`
}

// ExposedServiceHandler handles exposed service requests
type ExposedServiceHandler struct {
	state state.State
}

// NewExposedServiceHandler creates a new ExposedServiceHandler
func NewExposedServiceHandler(s state.State) *ExposedServiceHandler {
	return &ExposedServiceHandler{state: s}
}

// ListExposedServices godoc
// @Summary      List exposed services
// @Description  Get a list of all exposed services
// @Tags         infrastructure
// @Produce      json
// @Success      200  {array}   ExposedServiceResponse
// @Failure      500  {object}  map[string]string
// @Router       /exposed-services [get]
func (h *ExposedServiceHandler) ListExposedServices(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ExposedServiceType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing exposed services: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var services []ExposedServiceResponse
	for _, item := range items.Items {
		es, ok := item.(*omni.ExposedService)
		if !ok {
			log.Printf("Warning: resource is not an exposed service: %T", item)
			continue
		}

		spec := es.TypedSpec().Value
		serviceID := es.Metadata().ID()
		resp := ExposedServiceResponse{
			ID:               serviceID,
			Namespace:        es.Metadata().Namespace(),
			Port:             spec.Port,
			Label:            spec.Label,
			IconBase64:       spec.IconBase64,
			URL:              spec.Url,
			HasExplicitAlias: spec.HasExplicitAlias,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/exposed-services/"+serviceID),
			},
		}

		if spec.Error != "" {
			resp.Error = spec.Error
		}

		services = append(services, resp)
	}

	c.JSON(http.StatusOK, services)
}

// GetExposedService godoc
// @Summary      Get an exposed service
// @Description  Get detailed information about a specific exposed service
// @Tags         infrastructure
// @Produce      json
// @Param        id   path      string  true  "Exposed Service ID"
// @Success      200  {object}  ExposedServiceResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /exposed-services/{id} [get]
func (h *ExposedServiceHandler) GetExposedService(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ExposedServiceType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting exposed service %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	es, ok := res.(*omni.ExposedService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := es.TypedSpec().Value
	serviceID := es.Metadata().ID()
	resp := ExposedServiceResponse{
		ID:               serviceID,
		Namespace:        es.Metadata().Namespace(),
		Port:             spec.Port,
		Label:            spec.Label,
		IconBase64:       spec.IconBase64,
		URL:              spec.Url,
		HasExplicitAlias: spec.HasExplicitAlias,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/exposed-services/"+serviceID),
		},
	}

	if spec.Error != "" {
		resp.Error = spec.Error
	}

	c.JSON(http.StatusOK, resp)
}
