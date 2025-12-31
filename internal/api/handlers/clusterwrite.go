package handlers

import (
	"net/http"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/gin-gonic/gin"
	"github.com/jubblin/omni-api/internal/client"
	omniresources "github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
)

// ClusterCreateRequest represents a request to create a cluster
type ClusterCreateRequest struct {
	ID                string `json:"id" binding:"required"`
	KubernetesVersion string `json:"kubernetes_version" binding:"required"`
	TalosVersion      string `json:"talos_version,omitempty"`
	Features          struct {
		WorkloadProxy bool `json:"workload_proxy,omitempty"`
		DiskEncryption bool `json:"disk_encryption,omitempty"`
	} `json:"features,omitempty"`
}

// ClusterUpdateRequest represents a request to update a cluster
type ClusterUpdateRequest struct {
	KubernetesVersion string `json:"kubernetes_version,omitempty"`
	TalosVersion      string `json:"talos_version,omitempty"`
	Features          struct {
		WorkloadProxy bool `json:"workload_proxy,omitempty"`
		DiskEncryption bool `json:"disk_encryption,omitempty"`
	} `json:"features,omitempty"`
}

// ClusterWriteHandler handles cluster write operations
type ClusterWriteHandler struct {
	state      state.State
	management client.ManagementService // Management service interface
}

// NewClusterWriteHandler creates a new ClusterWriteHandler
func NewClusterWriteHandler(s state.State, mgmt client.ManagementService) *ClusterWriteHandler {
	return &ClusterWriteHandler{
		state:      s,
		management: mgmt,
	}
}

// CreateCluster godoc
// @Summary      Create a new cluster
// @Description  Create a new cluster in Omni
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        cluster  body      ClusterCreateRequest  true  "Cluster creation request"
// @Success      201      {object}  ClusterResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /clusters [post]
func (h *ClusterWriteHandler) CreateCluster(c *gin.Context) {
	var req ClusterCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create cluster using Management service
	features := &client.ClusterFeatures{
		WorkloadProxy: req.Features.WorkloadProxy,
		DiskEncryption: req.Features.DiskEncryption,
	}
	
	err := h.management.CreateCluster(
		c.Request.Context(),
		req.ID,
		req.KubernetesVersion,
		req.TalosVersion,
		features,
	)
	
	if err != nil {
		handleManagementError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Cluster created successfully",
		"id": req.ID,
	})
}

// UpdateCluster godoc
// @Summary      Update a cluster
// @Description  Update an existing cluster in Omni
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Cluster ID"
// @Param        cluster  body      ClusterUpdateRequest  true  "Cluster update request"
// @Success      200      {object}  ClusterResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /clusters/{id} [put]
func (h *ClusterWriteHandler) UpdateCluster(c *gin.Context) {
	id := c.Param("id")
	var req ClusterUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing cluster (verify it exists)
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterType, id, resource.VersionUndefined)
	res, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
		return
	}

	// Verify it's a cluster resource
	_, ok := res.(*omni.Cluster)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	// Update cluster using Management service
	features := &client.ClusterFeatures{
		WorkloadProxy: req.Features.WorkloadProxy,
		DiskEncryption: req.Features.DiskEncryption,
	}
	
	err = h.management.UpdateCluster(
		c.Request.Context(),
		id,
		req.KubernetesVersion,
		req.TalosVersion,
		features,
	)
	
	if err != nil {
		handleManagementError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Cluster updated successfully",
		"id": id,
	})
}

// DeleteCluster godoc
// @Summary      Delete a cluster
// @Description  Delete a cluster from Omni
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      202  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id} [delete]
func (h *ClusterWriteHandler) DeleteCluster(c *gin.Context) {
	id := c.Param("id")

	// Verify cluster exists
	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterType, id, resource.VersionUndefined)
	_, err := h.state.Get(c.Request.Context(), md)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
		return
	}

	// Delete cluster using Management service
	err = h.management.DeleteCluster(c.Request.Context(), id)
	if err != nil {
		handleManagementError(c, err)
		return
	}
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Cluster deletion initiated",
		"id": id,
	})
}
