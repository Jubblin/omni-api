package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jubblin/omni-api/internal/client"
)

// ClusterActionRequest represents a request for cluster actions
type ClusterActionRequest struct {
	Version string `json:"version,omitempty"` // For upgrade actions
}

// ClusterActionsHandler handles cluster action operations
type ClusterActionsHandler struct {
	state      interface{} // State interface
	management client.ManagementService // Management service
	talos      client.TalosService      // Talos service
}

// NewClusterActionsHandler creates a new ClusterActionsHandler
func NewClusterActionsHandler(state interface{}, mgmt client.ManagementService, talos client.TalosService) *ClusterActionsHandler {
	return &ClusterActionsHandler{
		state:      state,
		management: mgmt,
		talos:      talos,
	}
}

// TriggerKubernetesUpgrade godoc
// @Summary      Trigger Kubernetes upgrade
// @Description  Trigger a Kubernetes version upgrade for a cluster
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Cluster ID"
// @Param        request  body      ClusterActionRequest  true  "Upgrade request"
// @Success      202      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /clusters/{id}/actions/kubernetes-upgrade [post]
func (h *ClusterActionsHandler) TriggerKubernetesUpgrade(c *gin.Context) {
	id := c.Param("id")
	var req ClusterActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trigger Kubernetes upgrade using Management service
	err := h.management.UpgradeKubernetes(c.Request.Context(), id, req.Version)
	if err != nil {
		handleManagementError(c, err)
		return
	}
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Kubernetes upgrade initiated",
		"cluster_id": id,
		"version": req.Version,
	})
}

// TriggerTalosUpgrade godoc
// @Summary      Trigger Talos upgrade
// @Description  Trigger a Talos OS version upgrade for a cluster
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Cluster ID"
// @Param        request  body      ClusterActionRequest  true  "Upgrade request"
// @Success      202      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /clusters/{id}/actions/talos-upgrade [post]
func (h *ClusterActionsHandler) TriggerTalosUpgrade(c *gin.Context) {
	id := c.Param("id")
	var req ClusterActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trigger Talos upgrade using Management service
	err := h.management.UpgradeTalos(c.Request.Context(), id, req.Version)
	if err != nil {
		handleManagementError(c, err)
		return
	}
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Talos upgrade initiated",
		"cluster_id": id,
		"version": req.Version,
	})
}

// TriggerBootstrap godoc
// @Summary      Trigger cluster bootstrap
// @Description  Trigger bootstrap for a cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      202  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/actions/bootstrap [post]
func (h *ClusterActionsHandler) TriggerBootstrap(c *gin.Context) {
	id := c.Param("id")

	// Trigger bootstrap using Management service
	err := h.management.BootstrapCluster(c.Request.Context(), id)
	if err != nil {
		handleManagementError(c, err)
		return
	}
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Bootstrap initiated",
		"cluster_id": id,
	})
}

// TriggerDestroy godoc
// @Summary      Trigger cluster destruction
// @Description  Trigger destruction/teardown of a cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      202  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/actions/destroy [post]
func (h *ClusterActionsHandler) TriggerDestroy(c *gin.Context) {
	id := c.Param("id")

	// Trigger cluster destruction using Management service
	err := h.management.DeleteCluster(c.Request.Context(), id)
	if err != nil {
		handleManagementError(c, err)
		return
	}
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Cluster destruction initiated",
		"cluster_id": id,
	})
}
