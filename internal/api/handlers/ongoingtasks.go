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

// OngoingTaskResponse represents the ongoing task information returned by the API
type OngoingTaskResponse struct {
	ID         string            `json:"id"`
	Namespace  string            `json:"namespace"`
	Title      string            `json:"title"`
	ResourceID string            `json:"resource_id,omitempty"`
	TaskType   string            `json:"task_type,omitempty"`
	Links      map[string]string `json:"_links,omitempty"`
}

// OngoingTaskHandler handles ongoing task requests
type OngoingTaskHandler struct {
	state state.State
}

// NewOngoingTaskHandler creates a new OngoingTaskHandler
func NewOngoingTaskHandler(s state.State) *OngoingTaskHandler {
	return &OngoingTaskHandler{state: s}
}

// ListOngoingTasks godoc
// @Summary      List all ongoing tasks
// @Description  Get a list of all currently running tasks in Omni
// @Tags         ongoingtasks
// @Produce      json
// @Param        resource   query     string  false  "Filter by resource ID"
// @Success      200  {array}   OngoingTaskResponse
// @Failure      500  {object}  map[string]string
// @Router       /ongoingtasks [get]
func (h *OngoingTaskHandler) ListOngoingTasks(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.OngoingTaskType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing ongoing tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resourceFilter := c.Query("resource")

	var tasks []OngoingTaskResponse
	for _, item := range items.Items {
		ot, ok := item.(*omni.OngoingTask)
		if !ok {
			log.Printf("Warning: resource is not an ongoing task: %T", item)
			continue
		}

		spec := ot.TypedSpec().Value
		
		// Filter by resource if specified
		if resourceFilter != "" && spec.ResourceId != resourceFilter {
			continue
		}

		taskID := ot.Metadata().ID()
		resp := OngoingTaskResponse{
			ID:         taskID,
			Namespace:  ot.Metadata().Namespace(),
			Title:      spec.Title,
			ResourceID: spec.ResourceId,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/ongoingtasks/"+taskID),
			},
		}

		// Determine task type from details
		if spec.Details != nil {
			switch spec.Details.(type) {
			case interface{ GetTalosUpgrade() interface{} }:
				resp.TaskType = "talos_upgrade"
			case interface{ GetKubernetesUpgrade() interface{} }:
				resp.TaskType = "kubernetes_upgrade"
			case interface{ GetDestroy() interface{} }:
				resp.TaskType = "destroy"
			case interface{ GetMachineUpgrade() interface{} }:
				resp.TaskType = "machine_upgrade"
			default:
				resp.TaskType = "unknown"
			}
		}

		// Add link to resource if resource ID is set
		if spec.ResourceId != "" {
			// Try to determine resource type from labels or context
			if clusterID, ok := ot.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
				resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
			}
		}

		tasks = append(tasks, resp)
	}

	c.JSON(http.StatusOK, tasks)
}

// GetOngoingTask godoc
// @Summary      Get a single ongoing task
// @Description  Get detailed information about a specific ongoing task
// @Tags         ongoingtasks
// @Produce      json
// @Param        id   path      string  true  "Ongoing Task ID"
// @Success      200  {object}  OngoingTaskResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /ongoingtasks/{id} [get]
func (h *OngoingTaskHandler) GetOngoingTask(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.OngoingTaskType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting ongoing task %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ot, ok := res.(*omni.OngoingTask)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := ot.TypedSpec().Value
	taskID := ot.Metadata().ID()
	resp := OngoingTaskResponse{
		ID:         taskID,
		Namespace:  ot.Metadata().Namespace(),
		Title:      spec.Title,
		ResourceID: spec.ResourceId,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/ongoingtasks/"+taskID),
		},
	}

	// Determine task type from details
	if spec.Details != nil {
		switch spec.Details.(type) {
		case interface{ GetTalosUpgrade() interface{} }:
			resp.TaskType = "talos_upgrade"
		case interface{ GetKubernetesUpgrade() interface{} }:
			resp.TaskType = "kubernetes_upgrade"
		case interface{ GetDestroy() interface{} }:
			resp.TaskType = "destroy"
		case interface{ GetMachineUpgrade() interface{} }:
			resp.TaskType = "machine_upgrade"
		default:
			resp.TaskType = "unknown"
		}
	}

	// Add link to resource if resource ID is set
	if spec.ResourceId != "" {
		if clusterID, ok := ot.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
			resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
		}
	}

	c.JSON(http.StatusOK, resp)
}
