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

// MachineSetResponse represents the machine set information returned by the API
type MachineSetResponse struct {
	ID          string            `json:"id"`
	Namespace   string            `json:"namespace"`
	MachineClass string           `json:"machine_class,omitempty"`
	UpdateStrategy string         `json:"update_strategy,omitempty"`
	DeleteStrategy string         `json:"delete_strategy,omitempty"`
	Links       map[string]string `json:"_links,omitempty"`
}

// MachineSetHandler handles machine set requests
type MachineSetHandler struct {
	state state.State
}

// NewMachineSetHandler creates a new MachineSetHandler
func NewMachineSetHandler(s state.State) *MachineSetHandler {
	return &MachineSetHandler{state: s}
}

// ListMachineSets godoc
// @Summary      List all machine sets
// @Description  Get a list of all machine sets in Omni
// @Tags         machinesets
// @Produce      json
// @Success      200  {array}   MachineSetResponse
// @Failure      500  {object}  map[string]string
// @Router       /machinesets [get]
func (h *MachineSetHandler) ListMachineSets(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineSetType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing machine sets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var machineSets []MachineSetResponse
	for _, item := range items.Items {
		ms, ok := item.(*omni.MachineSet)
		if !ok {
			log.Printf("Warning: resource is not a machine set: %T", item)
			continue
		}

		machineSetID := ms.Metadata().ID()
		resp := MachineSetResponse{
			ID:        machineSetID,
			Namespace: ms.Metadata().Namespace(),
			Links: map[string]string{
				"self":   buildURL(c, "/api/v1/machinesets/"+machineSetID),
				"status": buildURL(c, "/api/v1/machinesets/"+machineSetID+"/status"),
			},
		}

		spec := ms.TypedSpec().Value
		if spec.MachineAllocation != nil {
			resp.MachineClass = spec.MachineAllocation.Name
		} else if spec.MachineClass != nil {
			resp.MachineClass = spec.MachineClass.Name
		}
		resp.UpdateStrategy = spec.UpdateStrategy.String()
		resp.DeleteStrategy = spec.DeleteStrategy.String()

		// Try to find cluster ID from labels
		if clusterID, ok := ms.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
			resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
		}

		machineSets = append(machineSets, resp)
	}

	c.JSON(http.StatusOK, machineSets)
}

// GetMachineSet godoc
// @Summary      Get a single machine set
// @Description  Get detailed information about a specific machine set
// @Tags         machinesets
// @Produce      json
// @Param        id   path      string  true  "Machine Set ID"
// @Success      200  {object}  MachineSetResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /machinesets/{id} [get]
func (h *MachineSetHandler) GetMachineSet(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.MachineSetType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting machine set %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ms, ok := res.(*omni.MachineSet)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	machineSetID := ms.Metadata().ID()
	resp := MachineSetResponse{
		ID:        machineSetID,
		Namespace: ms.Metadata().Namespace(),
		Links: map[string]string{
			"self":   buildURL(c, "/api/v1/machinesets/"+machineSetID),
			"status": buildURL(c, "/api/v1/machinesets/"+machineSetID+"/status"),
		},
	}

	spec := ms.TypedSpec().Value
	if spec.MachineAllocation != nil {
		resp.MachineClass = spec.MachineAllocation.Name
	} else if spec.MachineClass != nil {
		resp.MachineClass = spec.MachineClass.Name
	}
	resp.UpdateStrategy = spec.UpdateStrategy.String()
	resp.DeleteStrategy = spec.DeleteStrategy.String()

	// Try to find cluster ID from labels
	if clusterID, ok := ms.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
