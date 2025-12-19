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

// KernelArgsResponse represents the kernel args information
type KernelArgsResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Args      []string          `json:"args,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// KernelArgsHandler handles kernel args requests
type KernelArgsHandler struct {
	state state.State
}

// NewKernelArgsHandler creates a new KernelArgsHandler
func NewKernelArgsHandler(s state.State) *KernelArgsHandler {
	return &KernelArgsHandler{state: s}
}

// ListKernelArgs godoc
// @Summary      List kernel args
// @Description  Get a list of all kernel args configurations
// @Tags         machines
// @Produce      json
// @Success      200  {array}   KernelArgsResponse
// @Failure      500  {object}  map[string]string
// @Router       /kernel-args [get]
func (h *KernelArgsHandler) ListKernelArgs(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.KernelArgsType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing kernel args: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var kernelArgs []KernelArgsResponse
	for _, item := range items.Items {
		ka, ok := item.(*omni.KernelArgs)
		if !ok {
			log.Printf("Warning: resource is not kernel args: %T", item)
			continue
		}

		spec := ka.TypedSpec().Value
		argsID := ka.Metadata().ID()
		resp := KernelArgsResponse{
			ID:        argsID,
			Namespace: ka.Metadata().Namespace(),
			Args:      spec.Args,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/kernel-args/"+argsID),
			},
		}

		kernelArgs = append(kernelArgs, resp)
	}

	c.JSON(http.StatusOK, kernelArgs)
}

// GetKernelArgs godoc
// @Summary      Get kernel args
// @Description  Get detailed information about a specific kernel args configuration
// @Tags         machines
// @Produce      json
// @Param        id   path      string  true  "Kernel Args ID"
// @Success      200  {object}  KernelArgsResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /kernel-args/{id} [get]
func (h *KernelArgsHandler) GetKernelArgs(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.KernelArgsType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting kernel args %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ka, ok := res.(*omni.KernelArgs)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ka.TypedSpec().Value
	argsID := ka.Metadata().ID()
	resp := KernelArgsResponse{
		ID:        argsID,
		Namespace: ka.Metadata().Namespace(),
		Args:      spec.Args,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/kernel-args/"+argsID),
		},
	}

	c.JSON(http.StatusOK, resp)
}
