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

// ControlPlaneStatusResponse represents the control plane status information
type ControlPlaneStatusResponse struct {
	ID         string            `json:"id"`
	Namespace  string            `json:"namespace"`
	Conditions []Condition       `json:"conditions,omitempty"`
	Links      map[string]string `json:"_links,omitempty"`
}

// Condition represents a control plane condition
type Condition struct {
	Type    string `json:"type,omitempty"`
	Status  string `json:"status,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// ControlPlaneStatusHandler handles control plane status requests
type ControlPlaneStatusHandler struct {
	state state.State
}

// NewControlPlaneStatusHandler creates a new ControlPlaneStatusHandler
func NewControlPlaneStatusHandler(s state.State) *ControlPlaneStatusHandler {
	return &ControlPlaneStatusHandler{state: s}
}

// GetControlPlaneStatus godoc
// @Summary      Get control plane status
// @Description  Get control plane health status for a cluster
// @Tags         clusters
// @Produce      json
// @Param        id   path      string  true  "Cluster ID"
// @Success      200  {object}  ControlPlaneStatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clusters/{id}/controlplane-status [get]
func (h *ControlPlaneStatusHandler) GetControlPlaneStatus(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ControlPlaneStatusType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting control plane status for cluster %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "control plane status not found"})
		return
	}

	cps, ok := res.(*omni.ControlPlaneStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	spec := cps.TypedSpec().Value
	resp := ControlPlaneStatusResponse{
		ID:        cps.Metadata().ID(),
		Namespace: cps.Metadata().Namespace(),
		Links: map[string]string{
			"self":    buildURL(c, "/api/v1/clusters/"+id+"/controlplane-status"),
			"cluster": buildURL(c, "/api/v1/clusters/"+id),
		},
	}

	for _, cond := range spec.Conditions {
		resp.Conditions = append(resp.Conditions, Condition{
			Type:    cond.Type.String(),
			Status:   cond.Status.String(),
			Reason:   cond.Reason,
			Message: "", // Condition doesn't have Message field
		})
	}

	c.JSON(http.StatusOK, resp)
}
