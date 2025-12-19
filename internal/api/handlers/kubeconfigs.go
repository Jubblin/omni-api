package handlers

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/gin-gonic/gin"
	omniresources "github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
)

// KubeconfigResponse represents the kubeconfig information returned by the API
// WARNING: This contains sensitive cluster credentials
type KubeconfigResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Data      string            `json:"data"` // Base64 encoded kubeconfig
	Links     map[string]string `json:"_links,omitempty"`
}

// KubeconfigHandler handles kubeconfig requests
type KubeconfigHandler struct {
	state state.State
}

// NewKubeconfigHandler creates a new KubeconfigHandler
func NewKubeconfigHandler(s state.State) *KubeconfigHandler {
	return &KubeconfigHandler{state: s}
}

// GetKubeconfig godoc
// @Summary      Get cluster kubeconfig
// @Description  Get the Kubernetes kubeconfig for a cluster. WARNING: Contains sensitive credentials.
// @Tags         kubeconfigs
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  KubeconfigResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/kubeconfig [get]
func (h *KubeconfigHandler) GetKubeconfig(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.KubeconfigType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting kubeconfig %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kc, ok := res.(*omni.Kubeconfig)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := kc.TypedSpec().Value
	resp := KubeconfigResponse{
		ID:        kc.Metadata().ID(),
		Namespace: kc.Metadata().Namespace(),
		Data:      base64.StdEncoding.EncodeToString(spec.Data),
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/clusters/"+id+"/kubeconfig"),
			"cluster": buildURL(c, "/api/v1/clusters/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
