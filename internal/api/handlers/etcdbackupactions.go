package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EtcdBackupCreateRequest represents a request to trigger a manual etcd backup
type EtcdBackupCreateRequest struct {
	Cluster string `json:"cluster" binding:"required"`
}

// EtcdBackupActionsHandler handles etcd backup action operations
type EtcdBackupActionsHandler struct {
	state      interface{} // State interface
	management interface{} // Management service client
}

// NewEtcdBackupActionsHandler creates a new EtcdBackupActionsHandler
func NewEtcdBackupActionsHandler(state, mgmt interface{}) *EtcdBackupActionsHandler {
	return &EtcdBackupActionsHandler{
		state:      state,
		management: mgmt,
	}
}

// TriggerManualBackup godoc
// @Summary      Trigger manual etcd backup
// @Description  Trigger a manual etcd backup for a cluster
// @Tags         etcdbackups
// @Accept       json
// @Produce      json
// @Param        request  body      EtcdBackupCreateRequest  true  "Backup request"
// @Success      202      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /etcdbackups [post]
func (h *EtcdBackupActionsHandler) TriggerManualBackup(c *gin.Context) {
	var req EtcdBackupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use Management service to trigger manual backup
	// Typically: managementClient.CreateEtcdManualBackup(ctx, clusterID)
	log.Printf("Triggering manual etcd backup for cluster %s (Management service integration needed)", req.Cluster)
	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Manual etcd backup initiated",
		"cluster": req.Cluster,
		"note": "Management service integration required for actual backup",
	})
}
