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

// EtcdManualBackupResponse represents the etcd manual backup information
type EtcdManualBackupResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	BackupAt  string            `json:"backup_at,omitempty"`
	Links     map[string]string `json:"_links,omitempty"`
}

// EtcdManualBackupHandler handles etcd manual backup requests
type EtcdManualBackupHandler struct {
	state state.State
}

// NewEtcdManualBackupHandler creates a new EtcdManualBackupHandler
func NewEtcdManualBackupHandler(s state.State) *EtcdManualBackupHandler {
	return &EtcdManualBackupHandler{state: s}
}

// ListEtcdManualBackups godoc
// @Summary      List etcd manual backups
// @Description  Get a list of all etcd manual backup requests
// @Tags         etcdbackups
// @Produce      json
// @Param        cluster   query     string  false  "Filter by cluster ID"
// @Success      200  {array}   EtcdManualBackupResponse
// @Failure      500  {object}  map[string]string
// @Router       /etcd-manual-backups [get]
func (h *EtcdManualBackupHandler) ListEtcdManualBackups(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.EtcdManualBackupType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing etcd manual backups: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clusterFilter := c.Query("cluster")

	var backups []EtcdManualBackupResponse
	for _, item := range items.Items {
		emb, ok := item.(*omni.EtcdManualBackup)
		if !ok {
			log.Printf("Warning: resource is not an etcd manual backup: %T", item)
			continue
		}

		// Filter by cluster if specified
		if clusterFilter != "" {
			if clusterID, ok := emb.Metadata().Labels().Get("omni.sidero.dev/cluster"); !ok || clusterID != clusterFilter {
				continue
			}
		}

		spec := emb.TypedSpec().Value
		backupID := emb.Metadata().ID()
		resp := EtcdManualBackupResponse{
			ID:        backupID,
			Namespace: emb.Metadata().Namespace(),
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/etcd-manual-backups/"+backupID),
			},
		}

		if spec.BackupAt != nil {
			resp.BackupAt = spec.BackupAt.AsTime().Format("2006-01-02T15:04:05Z")
		}

		// Try to find cluster ID from labels
		if clusterID, ok := emb.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
			resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
		}

		backups = append(backups, resp)
	}

	c.JSON(http.StatusOK, backups)
}

// GetEtcdManualBackup godoc
// @Summary      Get an etcd manual backup
// @Description  Get detailed information about a specific etcd manual backup request
// @Tags         etcdbackups
// @Produce      json
// @Param        id   path      string  true  "Etcd Manual Backup ID"
// @Success      200  {object}  EtcdManualBackupResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /etcd-manual-backups/{id} [get]
func (h *EtcdManualBackupHandler) GetEtcdManualBackup(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.EtcdManualBackupType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting etcd manual backup %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	emb, ok := res.(*omni.EtcdManualBackup)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := emb.TypedSpec().Value
	backupID := emb.Metadata().ID()
	resp := EtcdManualBackupResponse{
		ID:        backupID,
		Namespace: emb.Metadata().Namespace(),
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/etcd-manual-backups/"+backupID),
		},
	}

	if spec.BackupAt != nil {
		resp.BackupAt = spec.BackupAt.AsTime().Format("2006-01-02T15:04:05Z")
	}

	// Try to find cluster ID from labels
	if clusterID, ok := emb.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
