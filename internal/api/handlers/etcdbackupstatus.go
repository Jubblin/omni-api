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

// EtcdBackupStatusResponse represents the etcd backup status information
type EtcdBackupStatusResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Status    string            `json:"status,omitempty"`
	Error     string            `json:"error,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// EtcdBackupStatusHandler handles etcd backup status requests
type EtcdBackupStatusHandler struct {
	state state.State
}

// NewEtcdBackupStatusHandler creates a new EtcdBackupStatusHandler
func NewEtcdBackupStatusHandler(s state.State) *EtcdBackupStatusHandler {
	return &EtcdBackupStatusHandler{state: s}
}

// GetEtcdBackupStatus godoc
// @Summary      Get etcd backup status
// @Description  Get status of an etcd backup operation
// @Tags         etcdbackups
// @Produce      json
// @Param        id   path      string  true  "Etcd Backup ID"
// @Success      200  {object}  EtcdBackupStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /etcdbackups/{id}/status [get]
func (h *EtcdBackupStatusHandler) GetEtcdBackupStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.EtcdBackupStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting etcd backup status %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "etcd backup status not found"})
		return
	}

	ebs, ok := res.(*omni.EtcdBackupStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ebs.TypedSpec().Value
	resp := EtcdBackupStatusResponse{
		ID:        ebs.Metadata().ID(),
		Namespace: ebs.Metadata().Namespace(),
		Status:    spec.Status.String(),
		Links: map[string]string{
			"self":      buildURL(c, "/api/v1/etcdbackups/"+id+"/status"),
			"etcdbackup": buildURL(c, "/api/v1/etcdbackups/"+id),
		},
	}

	if spec.Error != "" {
		resp.Error = spec.Error
	}

	// Try to find cluster ID from labels
	if clusterID, ok := ebs.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
