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

// InstallationMediaResponse represents the installation media information
type InstallationMediaResponse struct {
	ID              string            `json:"id"`
	Namespace       string            `json:"namespace"`
	Name            string            `json:"name,omitempty"`
	Architecture    string            `json:"architecture,omitempty"`
	Profile         string            `json:"profile,omitempty"`
	ContentType     string            `json:"content_type,omitempty"`
	SrcFilePrefix   string            `json:"src_file_prefix,omitempty"`
	DestFilePrefix  string            `json:"dest_file_prefix,omitempty"`
	Extension       string            `json:"extension,omitempty"`
	NoSecureBoot    bool              `json:"no_secure_boot,omitempty"`
	Overlay         string            `json:"overlay,omitempty"`
	MinTalosVersion string            `json:"min_talos_version,omitempty"`
	Links           map[string]string `json:"_links,omitempty"`
}

// InstallationMediaHandler handles installation media requests
type InstallationMediaHandler struct {
	state state.State
}

// NewInstallationMediaHandler creates a new InstallationMediaHandler
func NewInstallationMediaHandler(s state.State) *InstallationMediaHandler {
	return &InstallationMediaHandler{state: s}
}

// ListInstallationMedias godoc
// @Summary      List all installation medias
// @Description  Get a list of all installation medias in Omni
// @Tags         installationmedias
// @Produce      json
// @Success      200  {array}   InstallationMediaResponse
// @Failure      500  {object}  map[string]string
// @Router       /installation-medias [get]
func (h *InstallationMediaHandler) ListInstallationMedias(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.InstallationMediaType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing installation medias: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var medias []InstallationMediaResponse
	for _, item := range items.Items {
		im, ok := item.(*omni.InstallationMedia)
		if !ok {
			log.Printf("Warning: resource is not an installation media: %T", item)
			continue
		}

		mediaID := im.Metadata().ID()
		spec := im.TypedSpec().Value
		resp := InstallationMediaResponse{
			ID:              mediaID,
			Namespace:       im.Metadata().Namespace(),
			Name:            spec.Name,
			Architecture:    spec.Architecture,
			Profile:         spec.Profile,
			ContentType:     spec.ContentType,
			SrcFilePrefix:   spec.SrcFilePrefix,
			DestFilePrefix:  spec.DestFilePrefix,
			Extension:       spec.Extension,
			NoSecureBoot:    spec.NoSecureBoot,
			Overlay:         spec.Overlay,
			MinTalosVersion: spec.MinTalosVersion,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/installation-medias/"+mediaID),
			},
		}

		medias = append(medias, resp)
	}

	c.JSON(http.StatusOK, medias)
}

// GetInstallationMedia godoc
// @Summary      Get a single installation media
// @Description  Get detailed information about a specific installation media
// @Tags         installationmedias
// @Produce      json
// @Param        id   path      string  true  "Installation Media ID"
// @Success      200  {object}  InstallationMediaResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /installation-medias/{id} [get]
func (h *InstallationMediaHandler) GetInstallationMedia(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.InstallationMediaType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting installation media %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "installation media not found"})
		return
	}

	im, ok := res.(*omni.InstallationMedia)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := im.TypedSpec().Value
	resp := InstallationMediaResponse{
		ID:              im.Metadata().ID(),
		Namespace:       im.Metadata().Namespace(),
		Name:            spec.Name,
		Architecture:    spec.Architecture,
		Profile:         spec.Profile,
		ContentType:     spec.ContentType,
		SrcFilePrefix:   spec.SrcFilePrefix,
		DestFilePrefix:  spec.DestFilePrefix,
		Extension:       spec.Extension,
		NoSecureBoot:    spec.NoSecureBoot,
		Overlay:         spec.Overlay,
		MinTalosVersion: spec.MinTalosVersion,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/installation-medias/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
