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

// ClusterMachineConfigStatusResponse represents the cluster machine config status information
type ClusterMachineConfigStatusResponse struct {
	ID                           string            `json:"id"`
	Namespace                    string            `json:"namespace"`
	ClusterMachineConfigVersion  string            `json:"cluster_machine_config_version,omitempty"`
	ClusterMachineVersion        string            `json:"cluster_machine_version,omitempty"`
	ClusterMachineConfigSha256   string            `json:"cluster_machine_config_sha256,omitempty"`
	LastConfigError              string            `json:"last_config_error,omitempty"`
	TalosVersion                 string            `json:"talos_version,omitempty"`
	SchematicID                  string            `json:"schematic_id,omitempty"`
	Links                        map[string]string `json:"_links,omitempty"`
}

// ClusterMachineConfigStatusHandler handles cluster machine config status requests
type ClusterMachineConfigStatusHandler struct {
	state state.State
}

// NewClusterMachineConfigStatusHandler creates a new ClusterMachineConfigStatusHandler
func NewClusterMachineConfigStatusHandler(s state.State) *ClusterMachineConfigStatusHandler {
	return &ClusterMachineConfigStatusHandler{state: s}
}

// GetClusterMachineConfigStatus godoc
// @Summary      Get cluster machine config status
// @Description  Get configuration status information for a specific cluster machine
// @Tags         clustermachines
// @Produce      json
// @Param        id   path      string  true  "Cluster Machine ID"
// @Success      200  {object}  ClusterMachineConfigStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clustermachines/{id}/config-status [get]
func (h *ClusterMachineConfigStatusHandler) GetClusterMachineConfigStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterMachineConfigStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster machine config status %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster machine config status not found"})
		return
	}

	cmcs, ok := res.(*omni.ClusterMachineConfigStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := cmcs.TypedSpec().Value
	clusterMachineID := cmcs.Metadata().ID()
	resp := ClusterMachineConfigStatusResponse{
		ID:                          clusterMachineID,
		Namespace:                   cmcs.Metadata().Namespace(),
		ClusterMachineConfigVersion: spec.ClusterMachineConfigVersion,
		ClusterMachineVersion:       spec.ClusterMachineVersion,
		ClusterMachineConfigSha256:  spec.ClusterMachineConfigSha256,
		LastConfigError:             spec.LastConfigError,
		TalosVersion:                spec.TalosVersion,
		SchematicID:                 spec.SchematicId,
		Links: map[string]string{
			"self":           buildURL(c, "/api/v1/clustermachines/"+clusterMachineID+"/config-status"),
			"clustermachine": buildURL(c, "/api/v1/clustermachines/"+clusterMachineID),
		},
	}

	// Add schematic link if available
	if spec.SchematicId != "" {
		resp.Links["schematic"] = buildURL(c, "/api/v1/schematics/"+spec.SchematicId)
	}

	// Add cluster link if available
	if clusterID, ok := cmcs.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	// Add machine link
	resp.Links["machine"] = buildURL(c, "/api/v1/machines/"+clusterMachineID)

	c.JSON(http.StatusOK, resp)
}
