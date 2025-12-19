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

// ImagePullStatusResponse represents the image pull status information
type ImagePullStatusResponse struct {
	ID                 string            `json:"id"`
	Namespace          string            `json:"namespace"`
	LastProcessedNode  string            `json:"last_processed_node,omitempty"`
	LastProcessedImage string            `json:"last_processed_image,omitempty"`
	LastProcessedError string            `json:"last_processed_error,omitempty"`
	ProcessedCount     uint32            `json:"processed_count,omitempty"`
	TotalCount         uint32            `json:"total_count,omitempty"`
	RequestVersion     string            `json:"request_version,omitempty"`
	Links              map[string]string `json:"_links,omitempty"`
}

// ImagePullStatusHandler handles image pull status requests
type ImagePullStatusHandler struct {
	state state.State
}

// NewImagePullStatusHandler creates a new ImagePullStatusHandler
func NewImagePullStatusHandler(s state.State) *ImagePullStatusHandler {
	return &ImagePullStatusHandler{state: s}
}

// GetImagePullStatus godoc
// @Summary      Get image pull status
// @Description  Get the status of an image pull operation
// @Tags         imagepullrequests
// @Produce      json
// @Param        id   path      string  true  "Image Pull Request ID"
// @Success      200  {object}  ImagePullStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /image-pull-requests/{id}/status [get]
func (h *ImagePullStatusHandler) GetImagePullStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ImagePullStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting image pull status %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "image pull status not found"})
		return
	}

	ips, ok := res.(*omni.ImagePullStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ips.TypedSpec().Value
	resp := ImagePullStatusResponse{
		ID:                 ips.Metadata().ID(),
		Namespace:          ips.Metadata().Namespace(),
		LastProcessedNode:  spec.LastProcessedNode,
		LastProcessedImage: spec.LastProcessedImage,
		LastProcessedError: spec.LastProcessedError,
		ProcessedCount:     spec.ProcessedCount,
		TotalCount:         spec.TotalCount,
		RequestVersion:     spec.RequestVersion,
		Links: map[string]string{
			"self":   buildURL(c, "/api/v1/image-pull-requests/"+id+"/status"),
			"request": buildURL(c, "/api/v1/image-pull-requests/"+id),
		},
	}

	c.JSON(http.StatusOK, resp)
}
