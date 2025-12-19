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

// EtcdBackupResponse represents the etcd backup information returned by the API
type EtcdBackupResponse struct {
	ID          string            `json:"id"`
	Namespace   string            `json:"namespace"`
	Snapshot    string            `json:"snapshot,omitempty"`
	CreatedAt   string            `json:"created_at,omitempty"`
	Size        uint64            `json:"size,omitempty"`
	Links       map[string]string `json:"_links,omitempty"`
}

// EtcdBackupHandler handles etcd backup requests
type EtcdBackupHandler struct {
	state state.State
}

// NewEtcdBackupHandler creates a new EtcdBackupHandler
func NewEtcdBackupHandler(s state.State) *EtcdBackupHandler {
	return &EtcdBackupHandler{state: s}
}

// ListEtcdBackups godoc
// @Summary      List all etcd backups
// @Description  Get a list of all etcd backups in Omni
// @Tags         etcdbackups
// @Produce      json
// @Param        cluster   query     string  false  "Filter by cluster ID"
// @Success      200  {array}   EtcdBackupResponse
// @Failure      500  {object}  map[string]string
// @Router       /etcdbackups [get]
func (h *EtcdBackupHandler) ListEtcdBackups(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.EtcdBackupType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing etcd backups: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clusterFilter := c.Query("cluster")

	var backups []EtcdBackupResponse
	for _, item := range items.Items {
		eb, ok := item.(*omni.EtcdBackup)
		if !ok {
			log.Printf("Warning: resource is not an etcd backup: %T", item)
			continue
		}

		// Filter by cluster if specified
		if clusterFilter != "" {
			if clusterID, ok := eb.Metadata().Labels().Get("omni.sidero.dev/cluster"); !ok || clusterID != clusterFilter {
				continue
			}
		}

		spec := eb.TypedSpec().Value
		backupID := eb.Metadata().ID()
		resp := EtcdBackupResponse{
			ID:        backupID,
			Namespace: eb.Metadata().Namespace(),
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/etcdbackups/"+backupID),
			},
		}

		if spec.Snapshot != "" {
			resp.Snapshot = spec.Snapshot
		}
		if spec.CreatedAt != nil {
			resp.CreatedAt = spec.CreatedAt.AsTime().Format("2006-01-02T15:04:05Z")
		}
		if spec.Size != 0 {
			resp.Size = spec.Size
		}

		// Try to find cluster ID from labels
		if clusterID, ok := eb.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
			resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
		}

		backups = append(backups, resp)
	}

	c.JSON(http.StatusOK, backups)
}

// GetEtcdBackup godoc
// @Summary      Get a single etcd backup
// @Description  Get detailed information about a specific etcd backup
// @Tags         etcdbackups
// @Produce      json
// @Param        id   path      string  true  "Etcd Backup ID"
// @Success      200  {object}  EtcdBackupResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /etcdbackups/{id} [get]
func (h *EtcdBackupHandler) GetEtcdBackup(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.EtcdBackupType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting etcd backup %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	eb, ok := res.(*omni.EtcdBackup)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := eb.TypedSpec().Value
	backupID := eb.Metadata().ID()
	resp := EtcdBackupResponse{
		ID:        backupID,
		Namespace: eb.Metadata().Namespace(),
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/etcdbackups/"+backupID),
		},
	}

	if spec.Snapshot != "" {
		resp.Snapshot = spec.Snapshot
	}
	if spec.CreatedAt != nil {
		resp.CreatedAt = spec.CreatedAt.AsTime().Format("2006-01-02T15:04:05Z")
	}
	if spec.Size != 0 {
		resp.Size = spec.Size
	}

	// Try to find cluster ID from labels
	if clusterID, ok := eb.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
