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

// NodeImageList represents a node and its images
type NodeImageList struct {
	Node   string   `json:"node"`
	Images []string `json:"images,omitempty"`
}

// ImagePullRequestResponse represents the image pull request information
type ImagePullRequestResponse struct {
	ID            string            `json:"id"`
	Namespace     string            `json:"namespace"`
	NodeImageList []NodeImageList   `json:"node_image_list,omitempty"`
	Links         map[string]string `json:"_links,omitempty"`
}

// ImagePullRequestHandler handles image pull request requests
type ImagePullRequestHandler struct {
	state state.State
}

// NewImagePullRequestHandler creates a new ImagePullRequestHandler
func NewImagePullRequestHandler(s state.State) *ImagePullRequestHandler {
	return &ImagePullRequestHandler{state: s}
}

// ListImagePullRequests godoc
// @Summary      List all image pull requests
// @Description  Get a list of all image pull requests in Omni
// @Tags         imagepullrequests
// @Produce      json
// @Success      200  {array}   ImagePullRequestResponse
// @Failure      500  {object}  map[string]string
// @Router       /image-pull-requests [get]
func (h *ImagePullRequestHandler) ListImagePullRequests(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ImagePullRequestType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing image pull requests: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var requests []ImagePullRequestResponse
	for _, item := range items.Items {
		ipr, ok := item.(*omni.ImagePullRequest)
		if !ok {
			log.Printf("Warning: resource is not an image pull request: %T", item)
			continue
		}

		requestID := ipr.Metadata().ID()
		spec := ipr.TypedSpec().Value
		resp := ImagePullRequestResponse{
			ID:            requestID,
			Namespace:     ipr.Metadata().Namespace(),
			NodeImageList: make([]NodeImageList, 0, len(spec.NodeImageList)),
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/image-pull-requests/"+requestID),
				"status": buildURL(c, "/api/v1/image-pull-requests/"+requestID+"/status"),
			},
		}

		for _, nodeImageList := range spec.NodeImageList {
			resp.NodeImageList = append(resp.NodeImageList, NodeImageList{
				Node:   nodeImageList.Node,
				Images: nodeImageList.Images,
			})
		}

		requests = append(requests, resp)
	}

	c.JSON(http.StatusOK, requests)
}

// GetImagePullRequest godoc
// @Summary      Get a single image pull request
// @Description  Get detailed information about a specific image pull request
// @Tags         imagepullrequests
// @Produce      json
// @Param        id   path      string  true  "Image Pull Request ID"
// @Success      200  {object}  ImagePullRequestResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /image-pull-requests/{id} [get]
func (h *ImagePullRequestHandler) GetImagePullRequest(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ImagePullRequestType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting image pull request %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "image pull request not found"})
		return
	}

	ipr, ok := res.(*omni.ImagePullRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ipr.TypedSpec().Value
	resp := ImagePullRequestResponse{
		ID:            ipr.Metadata().ID(),
		Namespace:     ipr.Metadata().Namespace(),
		NodeImageList: make([]NodeImageList, 0, len(spec.NodeImageList)),
		Links: map[string]string{
			"self":   buildURL(c, "/api/v1/image-pull-requests/"+id),
			"status": buildURL(c, "/api/v1/image-pull-requests/"+id+"/status"),
		},
	}

	for _, nil := range spec.NodeImageList {
		resp.NodeImageList = append(resp.NodeImageList, NodeImageList{
			Node:   nil.Node,
			Images: nil.Images,
		})
	}

	c.JSON(http.StatusOK, resp)
}
