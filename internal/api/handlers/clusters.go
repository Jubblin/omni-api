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

// ClusterStatusResponse represents the cluster status returned by the API
type ClusterStatusResponse struct {
	Available                 bool   `json:"available"`
	Phase                     string `json:"phase"`
	Ready                     bool   `json:"ready"`
	KubernetesAPIReady        bool   `json:"kubernetes_api_ready"`
	ControlplaneReady         bool   `json:"controlplane_ready"`
	HasConnectedControlPlanes bool   `json:"has_connected_control_planes"`
	Machines                  struct {
		Total     uint32 `json:"total"`
		Healthy   uint32 `json:"healthy"`
		Connected uint32 `json:"connected"`
		Requested uint32 `json:"requested"`
	} `json:"machines"`
}

// ClusterMetricsResponse represents the cluster metrics returned by the API
type ClusterMetricsResponse struct {
	Features map[string]uint32 `json:"features"`
}

// ClusterBootstrapResponse represents the cluster bootstrap status returned by the API
type ClusterBootstrapResponse struct {
	Bootstrapped bool `json:"bootstrapped"`
}

// ClusterResponse represents the cluster information returned by the API
type ClusterResponse struct {
	ID                string `json:"id"`
	Namespace         string `json:"namespace"`
	KubernetesVersion string `json:"kubernetes_version"`
	TalosVersion      string `json:"talos_version"`
	Features          struct {
		WorkloadProxy bool `json:"workload_proxy"`
		DiskEncryption bool `json:"disk_encryption"`
	} `json:"features"`
	Links map[string]string `json:"_links,omitempty"`
}

// ClusterHandler handles cluster requests
type ClusterHandler struct {
	state state.State
}

// NewClusterHandler creates a new ClusterHandler
func NewClusterHandler(s state.State) *ClusterHandler {
	return &ClusterHandler{state: s}
}

// ListClusters godoc
// @Summary      List all clusters
// @Description  Get a list of all clusters in Omni
// @Tags         clusters
// @Produce      json
// @Success      200  {array}   ClusterResponse
// @Failure      500  {object}  map[string]string
// @Router       /clusters [get]
func (h *ClusterHandler) ListClusters(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterType, "", resource.VersionUndefined)
	
	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing clusters: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var clusters []ClusterResponse
	for _, item := range items.Items {
		// Try to cast to the typed cluster resource
		cl, ok := item.(*omni.Cluster)
		if !ok {
			log.Printf("Warning: resource is not a cluster: %T", item)
			continue
		}

		clusterID := cl.Metadata().ID()
		res := ClusterResponse{
			ID:                clusterID,
			Namespace:         cl.Metadata().Namespace(),
			KubernetesVersion: cl.TypedSpec().Value.KubernetesVersion,
			TalosVersion:      cl.TypedSpec().Value.TalosVersion,
			Links: map[string]string{
				"self":      buildURL(c, "/api/v1/clusters/"+clusterID),
				"status":    buildURL(c, "/api/v1/clusters/"+clusterID+"/status"),
				"metrics":   buildURL(c, "/api/v1/clusters/"+clusterID+"/metrics"),
				"bootstrap": buildURL(c, "/api/v1/clusters/"+clusterID+"/bootstrap"),
				"machines":  buildURL(c, "/api/v1/machines?cluster="+clusterID),
			},
		}

		if cl.TypedSpec().Value.Features != nil {
			res.Features.WorkloadProxy = cl.TypedSpec().Value.Features.EnableWorkloadProxy
			res.Features.DiskEncryption = cl.TypedSpec().Value.Features.DiskEncryption
		}

		clusters = append(clusters, res)
	}

	c.JSON(http.StatusOK, clusters)
}

// GetCluster godoc
// @Summary      Get a single cluster
// @Description  Get detailed information about a specific cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  ClusterResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id} [get]
func (h *ClusterHandler) GetCluster(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cl, ok := res.(*omni.Cluster)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	clusterID := cl.Metadata().ID()
	resp := ClusterResponse{
		ID:                clusterID,
		Namespace:         cl.Metadata().Namespace(),
		KubernetesVersion: cl.TypedSpec().Value.KubernetesVersion,
		TalosVersion:      cl.TypedSpec().Value.TalosVersion,
		Links: map[string]string{
			"self":      buildURL(c, "/api/v1/clusters/"+clusterID),
			"status":    buildURL(c, "/api/v1/clusters/"+clusterID+"/status"),
			"metrics":   buildURL(c, "/api/v1/clusters/"+clusterID+"/metrics"),
			"bootstrap": buildURL(c, "/api/v1/clusters/"+clusterID+"/bootstrap"),
			"machines":  buildURL(c, "/api/v1/machines?cluster="+clusterID),
		},
	}

	if cl.TypedSpec().Value.Features != nil {
		resp.Features.WorkloadProxy = cl.TypedSpec().Value.Features.EnableWorkloadProxy
		resp.Features.DiskEncryption = cl.TypedSpec().Value.Features.DiskEncryption
	}

	c.JSON(http.StatusOK, resp)
}

// GetClusterStatus godoc
// @Summary      Get cluster status
// @Description  Get health and phase information for a specific cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  ClusterStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/status [get]
func (h *ClusterHandler) GetClusterStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster status %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cs, ok := res.(*omni.ClusterStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := cs.TypedSpec().Value
	resp := ClusterStatusResponse{
		Available:                 spec.Available,
		Phase:                     spec.Phase.String(),
		Ready:                     spec.Ready,
		KubernetesAPIReady:        spec.KubernetesAPIReady,
		ControlplaneReady:         spec.ControlplaneReady,
		HasConnectedControlPlanes: spec.HasConnectedControlPlanes,
	}

	if spec.Machines != nil {
		resp.Machines.Total = spec.Machines.Total
		resp.Machines.Healthy = spec.Machines.Healthy
		resp.Machines.Connected = spec.Machines.Connected
		resp.Machines.Requested = spec.Machines.Requested
	}

	c.JSON(http.StatusOK, resp)
}

// GetClusterMetrics godoc
// @Summary      Get cluster metrics
// @Description  Get real-time metrics for a specific cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  ClusterMetricsResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/metrics [get]
func (h *ClusterHandler) GetClusterMetrics(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterMetricsType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster metrics %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cm, ok := res.(*omni.ClusterMetrics)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	c.JSON(http.StatusOK, ClusterMetricsResponse{
		Features: cm.TypedSpec().Value.Features,
	})
}

// GetClusterBootstrap godoc
// @Summary      Get cluster bootstrap status
// @Description  Get the bootstrap status for a specific cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  ClusterBootstrapResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/bootstrap [get]
func (h *ClusterHandler) GetClusterBootstrap(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterBootstrapStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster bootstrap status %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cb, ok := res.(*omni.ClusterBootstrapStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	c.JSON(http.StatusOK, ClusterBootstrapResponse{
		Bootstrapped: cb.TypedSpec().Value.Bootstrapped,
	})
}

