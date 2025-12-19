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

// KubernetesVersionResponse represents the Kubernetes version information
type KubernetesVersionResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Version   string            `json:"version,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// KubernetesVersionHandler handles Kubernetes version requests
type KubernetesVersionHandler struct {
	state state.State
}

// NewKubernetesVersionHandler creates a new KubernetesVersionHandler
func NewKubernetesVersionHandler(s state.State) *KubernetesVersionHandler {
	return &KubernetesVersionHandler{state: s}
}

// ListKubernetesVersions godoc
// @Summary      List all Kubernetes versions
// @Description  Get a list of all available Kubernetes versions
// @Tags         kubernetes
// @Produce      json
// @Success      200  {array}   KubernetesVersionResponse
// @Failure      500  {object}  map[string]string
// @Router       /kubernetes-versions [get]
func (h *KubernetesVersionHandler) ListKubernetesVersions(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.KubernetesVersionType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing Kubernetes versions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var versions []KubernetesVersionResponse
	for _, item := range items.Items {
		kv, ok := item.(*omni.KubernetesVersion)
		if !ok {
			log.Printf("Warning: resource is not a Kubernetes version: %T", item)
			continue
		}

		spec := kv.TypedSpec().Value
		versionID := kv.Metadata().ID()
		resp := KubernetesVersionResponse{
			ID:        versionID,
			Namespace: kv.Metadata().Namespace(),
			Version:   spec.Version,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/kubernetes-versions/"+versionID),
			},
		}

		versions = append(versions, resp)
	}

	c.JSON(http.StatusOK, versions)
}

// GetKubernetesVersion godoc
// @Summary      Get a Kubernetes version
// @Description  Get detailed information about a specific Kubernetes version
// @Tags         kubernetes
// @Produce      json
// @Param        id   path      string  true  "Kubernetes Version ID"
// @Success      200  {object}  KubernetesVersionResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /kubernetes-versions/{id} [get]
func (h *KubernetesVersionHandler) GetKubernetesVersion(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.KubernetesVersionType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting Kubernetes version %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kv, ok := res.(*omni.KubernetesVersion)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := kv.TypedSpec().Value
	versionID := kv.Metadata().ID()
	resp := KubernetesVersionResponse{
		ID:        versionID,
		Namespace: kv.Metadata().Namespace(),
		Version:   spec.Version,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/kubernetes-versions/"+versionID),
		},
	}

	c.JSON(http.StatusOK, resp)
}
