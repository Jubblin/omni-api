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

// KubernetesStatusResponse represents the Kubernetes status information
type KubernetesStatusResponse struct {
	ID         string            `json:"id"`
	Namespace  string            `json:"namespace"`
	Nodes      []NodeStatus      `json:"nodes,omitempty"`
	StaticPods []NodeStaticPods  `json:"static_pods,omitempty"`
	Links      map[string]string `json:"_links,omitempty"`
}

// NodeStatus represents the status of a Kubernetes node
type NodeStatus struct {
	Nodename string `json:"nodename,omitempty"`
	Ready    bool   `json:"ready,omitempty"`
	Kubelet  bool   `json:"kubelet,omitempty"`
	APIServer bool  `json:"api_server,omitempty"`
}

// NodeStaticPods represents static pods on a node
type NodeStaticPods struct {
	Nodename string   `json:"nodename,omitempty"`
	Pods     []PodStatus `json:"pods,omitempty"`
}

// PodStatus represents the status of a static pod
type PodStatus struct {
	App     string `json:"app,omitempty"`
	Version string `json:"version,omitempty"`
	Ready   bool   `json:"ready,omitempty"`
}

// KubernetesStatusHandler handles Kubernetes status requests
type KubernetesStatusHandler struct {
	state state.State
}

// NewKubernetesStatusHandler creates a new KubernetesStatusHandler
func NewKubernetesStatusHandler(s state.State) *KubernetesStatusHandler {
	return &KubernetesStatusHandler{state: s}
}

// GetKubernetesStatus godoc
// @Summary      Get Kubernetes status
// @Description  Get Kubernetes cluster status including nodes and static pods
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  KubernetesStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/kubernetes-status [get]
func (h *KubernetesStatusHandler) GetKubernetesStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.KubernetesStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting Kubernetes status for cluster %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "kubernetes status not found"})
		return
	}

	ks, ok := res.(*omni.KubernetesStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ks.TypedSpec().Value
	resp := KubernetesStatusResponse{
		ID:        ks.Metadata().ID(),
		Namespace: ks.Metadata().Namespace(),
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/clusters/"+id+"/kubernetes-status"),
			"cluster": buildURL(c, "/api/v1/clusters/"+id),
		},
	}

	// Convert nodes
	for _, node := range spec.Nodes {
		resp.Nodes = append(resp.Nodes, NodeStatus{
			Nodename:  node.Nodename,
			Ready:     node.Ready,
			Kubelet:   node.KubeletVersion != "",
			APIServer: false, // Not available in NodeStatus
		})
	}

	// Convert static pods
	for _, nodePods := range spec.StaticPods {
		pods := []PodStatus{}
		for _, pod := range nodePods.StaticPods {
			pods = append(pods, PodStatus{
				App:     pod.App,
				Version: pod.Version,
				Ready:   pod.Ready,
			})
		}
		resp.StaticPods = append(resp.StaticPods, NodeStaticPods{
			Nodename: nodePods.Nodename,
			Pods:     pods,
		})
	}

	c.JSON(http.StatusOK, resp)
}
