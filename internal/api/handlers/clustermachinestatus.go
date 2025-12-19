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

// ClusterMachineStatusResponse represents the cluster machine status information
type ClusterMachineStatusResponse struct {
	ID                string            `json:"id"`
	Namespace         string            `json:"namespace"`
	Ready             bool               `json:"ready"`
	Stage             string             `json:"stage"`
	ApidAvailable     bool               `json:"apid_available,omitempty"`
	ConfigUpToDate    bool               `json:"config_up_to_date"`
	LastConfigError   string             `json:"last_config_error,omitempty"`
	ManagementAddress string             `json:"management_address,omitempty"`
	ConfigApplyStatus string             `json:"config_apply_status,omitempty"`
	IsRemoved         bool               `json:"is_removed"`
	ProvisionStatus   *ProvisionStatus   `json:"provision_status,omitempty"`
	Links             map[string]string  `json:"_links,omitempty"`
}

// ProvisionStatus represents the provision status details
type ProvisionStatus struct {
	ProviderID string `json:"provider_id,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
}

// ClusterMachineStatusHandler handles cluster machine status requests
type ClusterMachineStatusHandler struct {
	state state.State
}

// NewClusterMachineStatusHandler creates a new ClusterMachineStatusHandler
func NewClusterMachineStatusHandler(s state.State) *ClusterMachineStatusHandler {
	return &ClusterMachineStatusHandler{state: s}
}

// GetClusterMachineStatus godoc
// @Summary      Get cluster machine status
// @Description  Get detailed status information about a specific cluster machine
// @Tags         clustermachines
// @Produce      json
// @Param        id   path      string  true  "Cluster Machine ID"
// @Success      200  {object}  ClusterMachineStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clustermachines/{id}/status [get]
func (h *ClusterMachineStatusHandler) GetClusterMachineStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterMachineStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster machine status %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster machine status not found"})
		return
	}

	cms, ok := res.(*omni.ClusterMachineStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := cms.TypedSpec().Value
	clusterMachineID := cms.Metadata().ID()
	resp := ClusterMachineStatusResponse{
		ID:                clusterMachineID,
		Namespace:         cms.Metadata().Namespace(),
		Ready:             spec.Ready,
		Stage:             spec.Stage.String(),
		ApidAvailable:     spec.ApidAvailable,
		ConfigUpToDate:    spec.ConfigUpToDate,
		LastConfigError:   spec.LastConfigError,
		ManagementAddress: spec.ManagementAddress,
		ConfigApplyStatus: spec.ConfigApplyStatus.String(),
		IsRemoved:         spec.IsRemoved,
		Links: map[string]string{
			"self":         buildURL(c, "/api/v1/clustermachines/"+clusterMachineID+"/status"),
			"clustermachine": buildURL(c, "/api/v1/clustermachines/"+clusterMachineID),
		},
	}

	// Add provision status if available
	if spec.ProvisionStatus != nil {
		resp.ProvisionStatus = &ProvisionStatus{
			ProviderID: spec.ProvisionStatus.ProviderId,
			RequestID:  spec.ProvisionStatus.RequestId,
		}
	}

	// Add cluster link if available
	if clusterID, ok := cms.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	// Add machine link
	resp.Links["machine"] = buildURL(c, "/api/v1/machines/"+clusterMachineID)

	c.JSON(http.StatusOK, resp)
}
