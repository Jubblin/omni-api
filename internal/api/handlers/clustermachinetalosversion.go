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

// ClusterMachineTalosVersionResponse represents the cluster machine Talos version information
type ClusterMachineTalosVersionResponse struct {
	ID          string            `json:"id"`
	Namespace   string            `json:"namespace"`
	TalosVersion string           `json:"talos_version,omitempty"`
	SchematicID string            `json:"schematic_id,omitempty"`
	Links       map[string]string `json:"_links,omitempty"`
}

// ClusterMachineTalosVersionHandler handles cluster machine Talos version requests
type ClusterMachineTalosVersionHandler struct {
	state state.State
}

// NewClusterMachineTalosVersionHandler creates a new ClusterMachineTalosVersionHandler
func NewClusterMachineTalosVersionHandler(s state.State) *ClusterMachineTalosVersionHandler {
	return &ClusterMachineTalosVersionHandler{state: s}
}

// GetClusterMachineTalosVersion godoc
// @Summary      Get cluster machine Talos version
// @Description  Get Talos OS version and schematic information for a specific cluster machine
// @Tags         clustermachines
// @Produce      json
// @Param        id   path      string  true  "Cluster Machine ID"
// @Success      200  {object}  ClusterMachineTalosVersionResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clustermachines/{id}/talos-version [get]
func (h *ClusterMachineTalosVersionHandler) GetClusterMachineTalosVersion(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterMachineTalosVersionType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting cluster machine Talos version %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster machine Talos version not found"})
		return
	}

	cmtv, ok := res.(*omni.ClusterMachineTalosVersion)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := cmtv.TypedSpec().Value
	clusterMachineID := cmtv.Metadata().ID()
	resp := ClusterMachineTalosVersionResponse{
		ID:           clusterMachineID,
		Namespace:    cmtv.Metadata().Namespace(),
		TalosVersion: spec.TalosVersion,
		SchematicID:  spec.SchematicId,
		Links: map[string]string{
			"self":           buildURL(c, "/api/v1/clustermachines/"+clusterMachineID+"/talos-version"),
			"clustermachine": buildURL(c, "/api/v1/clustermachines/"+clusterMachineID),
		},
	}

	// Add schematic link if available
	if spec.SchematicId != "" {
		resp.Links["schematic"] = buildURL(c, "/api/v1/schematics/"+spec.SchematicId)
	}

	// Add cluster link if available
	if clusterID, ok := cmtv.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	// Add machine link
	resp.Links["machine"] = buildURL(c, "/api/v1/machines/"+clusterMachineID)

	c.JSON(http.StatusOK, resp)
}
