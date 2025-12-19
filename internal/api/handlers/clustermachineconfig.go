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

// ClusterMachineConfigResponse represents the cluster machine config information
type ClusterMachineConfigResponse struct {
	ID                    string            `json:"id"`
	Namespace             string            `json:"namespace"`
	ClusterMachineVersion string            `json:"cluster_machine_version,omitempty"`
	GenerationError       string            `json:"generation_error,omitempty"`
	WithoutComments       bool              `json:"without_comments,omitempty"`
	GrubUseUkiCmdline     bool              `json:"grub_use_uki_cmdline,omitempty"`
	ConfigSize            int               `json:"config_size,omitempty"`
	Links                 map[string]string `json:"_links,omitempty"`
}

// ClusterMachineConfigHandler handles cluster machine config requests
type ClusterMachineConfigHandler struct {
	state state.State
}

// NewClusterMachineConfigHandler creates a new ClusterMachineConfigHandler
func NewClusterMachineConfigHandler(s state.State) *ClusterMachineConfigHandler {
	return &ClusterMachineConfigHandler{state: s}
}

// GetClusterMachineConfig godoc
// @Summary      Get cluster machine config
// @Description  Get machine configuration for a specific cluster machine
// @Tags         clustermachines
// @Produce      json
// @Param        id   path      string  true  "Cluster Machine ID"
// @Success      200  {object}  ClusterMachineConfigResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clustermachines/{id}/config [get]
func (h *ClusterMachineConfigHandler) GetClusterMachineConfig(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterMachineConfigType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster machine config %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster machine config not found"})
		return
	}

	cmc, ok := res.(*omni.ClusterMachineConfig)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := cmc.TypedSpec().Value
	resp := ClusterMachineConfigResponse{
		ID:                    cmc.Metadata().ID(),
		Namespace:             cmc.Metadata().Namespace(),
		ClusterMachineVersion: spec.ClusterMachineVersion,
		WithoutComments:       spec.WithoutComments,
		GrubUseUkiCmdline:     spec.GrubUseUkiCmdline,
		Links: map[string]string{
			"self":          buildURL(c, "/api/v1/clustermachines/"+id+"/config"),
			"clustermachine": buildURL(c, "/api/v1/clustermachines/"+id),
		},
	}

	if spec.GenerationError != "" {
		resp.GenerationError = spec.GenerationError
	}

	// Get config size (from data or compressed data)
	data, err := spec.GetUncompressedData()
	if err == nil {
		defer data.Free()
		configData := data.Data()
		if len(configData) > 0 {
			resp.ConfigSize = len(configData)
		}
	}

	// Try to find cluster ID from labels
	if clusterID, ok := cmc.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
