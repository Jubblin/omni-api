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

// ClusterKubernetesNodeResponse represents the cluster Kubernetes node information
type ClusterKubernetesNodeResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Name      string            `json:"name,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// ClusterKubernetesNodesHandler handles cluster Kubernetes nodes requests
type ClusterKubernetesNodesHandler struct {
	state state.State
}

// NewClusterKubernetesNodesHandler creates a new ClusterKubernetesNodesHandler
func NewClusterKubernetesNodesHandler(s state.State) *ClusterKubernetesNodesHandler {
	return &ClusterKubernetesNodesHandler{state: s}
}

// ListClusterKubernetesNodes godoc
// @Summary      List cluster Kubernetes nodes
// @Description  Get a list of all Kubernetes nodes in a cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {array}   ClusterKubernetesNodeResponse
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/kubernetes-nodes [get]
func (h *ClusterKubernetesNodesHandler) ListClusterKubernetesNodes(c *gin.Context) {
	clusterID := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterKubernetesNodesType, clusterID, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster Kubernetes nodes for cluster %s: %v", clusterID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster kubernetes nodes not found"})
		return
	}

	ckn, ok := res.(*omni.ClusterKubernetesNodes)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ckn.TypedSpec().Value
	var nodes []ClusterKubernetesNodeResponse

	for _, nodeName := range spec.Nodes {
		nodeResp := ClusterKubernetesNodeResponse{
			ID:        nodeName,
			Namespace: ckn.Metadata().Namespace(),
			Name:      nodeName,
			Labels:    make(map[string]string),
			Links: map[string]string{
				"self":    buildURL(c, "/api/v1/clusters/"+clusterID+"/kubernetes-nodes/"+nodeName),
				"cluster": buildURL(c, "/api/v1/clusters/"+clusterID),
			},
		}
		nodes = append(nodes, nodeResp)
	}

	c.JSON(http.StatusOK, nodes)
}

// GetClusterKubernetesNode godoc
// @Summary      Get a cluster Kubernetes node
// @Description  Get detailed information about a specific Kubernetes node in a cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Param        node   path      string  true  "Node name"
// @Success      200  {object}  ClusterKubernetesNodeResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/kubernetes-nodes/{node} [get]
func (h *ClusterKubernetesNodesHandler) GetClusterKubernetesNode(c *gin.Context) {
	clusterID := c.Param("id")
	nodeName := c.Param("node")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterKubernetesNodesType, clusterID, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster Kubernetes nodes for cluster %s: %v", clusterID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster kubernetes nodes not found"})
		return
	}

	ckn, ok := res.(*omni.ClusterKubernetesNodes)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ckn.TypedSpec().Value
	for _, n := range spec.Nodes {
		if n == nodeName {
			resp := ClusterKubernetesNodeResponse{
				ID:        nodeName,
				Namespace: ckn.Metadata().Namespace(),
				Name:      nodeName,
				Labels:    make(map[string]string),
				Links: map[string]string{
					"self":    buildURL(c, "/api/v1/clusters/"+clusterID+"/kubernetes-nodes/"+nodeName),
					"cluster": buildURL(c, "/api/v1/clusters/"+clusterID),
				},
			}
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "kubernetes node not found"})
}
