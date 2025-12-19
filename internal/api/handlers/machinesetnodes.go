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

// MachineSetNodeResponse represents the machine set node information returned by the API
type MachineSetNodeResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Links     map[string]string `json:"_links,omitempty"`
}

// MachineSetNodeHandler handles machine set node requests
type MachineSetNodeHandler struct {
	state state.State
}

// NewMachineSetNodeHandler creates a new MachineSetNodeHandler
func NewMachineSetNodeHandler(s state.State) *MachineSetNodeHandler {
	return &MachineSetNodeHandler{state: s}
}

// ListMachineSetNodes godoc
// @Summary      List all machine set nodes
// @Description  Get a list of all machine set nodes in Omni
// @Tags         machinesetnodes
// @Produce      json
// @Param        machineset   query     string  false  "Filter by machine set ID"
// @Success      200  {array}   MachineSetNodeResponse
// @Failure      500  {object}  map[string]string
// @Router       /machinesetnodes [get]
func (h *MachineSetNodeHandler) ListMachineSetNodes(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineSetNodeType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing machine set nodes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	machineSetFilter := c.Query("machineset")

	var nodes []MachineSetNodeResponse
	for _, item := range items.Items {
		msn, ok := item.(*omni.MachineSetNode)
		if !ok {
			log.Printf("Warning: resource is not a machine set node: %T", item)
			continue
		}

		// Filter by machine set if specified
		if machineSetFilter != "" {
			if machineSetID, ok := msn.Metadata().Labels().Get("omni.sidero.dev/machine-set"); !ok || machineSetID != machineSetFilter {
				continue
			}
		}

		nodeID := msn.Metadata().ID()
		resp := MachineSetNodeResponse{
			ID:        nodeID,
			Namespace: msn.Metadata().Namespace(),
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/machinesetnodes/"+nodeID),
			},
		}

		// Try to find machine set ID from labels
		if machineSetID, ok := msn.Metadata().Labels().Get("omni.sidero.dev/machine-set"); ok {
			resp.Links["machineset"] = buildURL(c, "/api/v1/machinesets/"+machineSetID)
		}

		// Try to find cluster machine ID (MachineSetNode ID is typically the ClusterMachine ID)
		resp.Links["clustermachine"] = buildURL(c, "/api/v1/clustermachines/"+nodeID)

		// Try to find cluster ID from labels
		if clusterID, ok := msn.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
			resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
		}

		nodes = append(nodes, resp)
	}

	c.JSON(http.StatusOK, nodes)
}

// GetMachineSetNode godoc
// @Summary      Get a single machine set node
// @Description  Get detailed information about a specific machine set node
// @Tags         machinesetnodes
// @Produce      json
// @Param        id   path      string  true  "Machine Set Node ID"
// @Success      200  {object}  MachineSetNodeResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machinesetnodes/{id} [get]
func (h *MachineSetNodeHandler) GetMachineSetNode(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineSetNodeType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine set node %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msn, ok := res.(*omni.MachineSetNode)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	nodeID := msn.Metadata().ID()
	resp := MachineSetNodeResponse{
		ID:        nodeID,
		Namespace: msn.Metadata().Namespace(),
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/machinesetnodes/"+nodeID),
		},
	}

	// Try to find machine set ID from labels
	if machineSetID, ok := msn.Metadata().Labels().Get("omni.sidero.dev/machine-set"); ok {
		resp.Links["machineset"] = buildURL(c, "/api/v1/machinesets/"+machineSetID)
	}

	// Try to find cluster machine ID
	resp.Links["clustermachine"] = buildURL(c, "/api/v1/clustermachines/"+nodeID)

	// Try to find cluster ID from labels
	if clusterID, ok := msn.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
