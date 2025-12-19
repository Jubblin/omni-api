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

// ConfigPatchResponse represents the config patch information returned by the API
type ConfigPatchResponse struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Data      string            `json:"data"`
	Links     map[string]string `json:"_links,omitempty"`
}

// ConfigPatchHandler handles config patch requests
type ConfigPatchHandler struct {
	state state.State
}

// NewConfigPatchHandler creates a new ConfigPatchHandler
func NewConfigPatchHandler(s state.State) *ConfigPatchHandler {
	return &ConfigPatchHandler{state: s}
}

// ListConfigPatches godoc
// @Summary      List all config patches
// @Description  Get a list of all config patches in Omni
// @Tags         configpatches
// @Produce      json
// @Success      200  {array}   ConfigPatchResponse
// @Failure      500  {object}  map[string]string
// @Router       /configpatches [get]
func (h *ConfigPatchHandler) ListConfigPatches(c *gin.Context) {
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ConfigPatchType, "", resource.VersionUndefined)

	items, err := st.List(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error listing config patches: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var patches []ConfigPatchResponse
	for _, item := range items.Items {
		cp, ok := item.(*omni.ConfigPatch)
		if !ok {
			log.Printf("Warning: resource is not a config patch: %T", item)
			continue
		}

		data, err := cp.TypedSpec().Value.GetUncompressedData()
		dataStr := ""
		if err == nil {
			dataStr = string(data.Data())
			data.Free()
		} else {
			// Fallback to Data field if GetUncompressedData fails
			dataStr = cp.TypedSpec().Value.Data
		}
		
		patchID := cp.Metadata().ID()
		resp := ConfigPatchResponse{
			ID:        patchID,
			Namespace: cp.Metadata().Namespace(),
			Data:      dataStr,
			Links: map[string]string{
				"self": buildURL(c, "/api/v1/configpatches/"+patchID),
			},
		}

		// Try to find cluster ID from labels
		if clusterID, ok := cp.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
			resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
		}

		patches = append(patches, resp)
	}

	c.JSON(http.StatusOK, patches)
}

// GetConfigPatch godoc
// @Summary      Get a single config patch
// @Description  Get detailed information about a specific config patch
// @Tags         configpatches
// @Produce      json
// @Param        id   path      string  true  "Config Patch ID"
// @Success      200  {object}  ConfigPatchResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /configpatches/{id} [get]
func (h *ConfigPatchHandler) GetConfigPatch(c *gin.Context) {
	id := c.Param("id")
	st := h.state

	md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ConfigPatchType, id, resource.VersionUndefined)

	res, err := st.Get(c.Request.Context(), md)
	if err != nil {
		log.Printf("Error getting config patch %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cp, ok := res.(*omni.ConfigPatch)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: unexpected resource type"})
		return
	}

	data, err := cp.TypedSpec().Value.GetUncompressedData()
	dataStr := ""
	if err == nil {
		dataStr = string(data.Data())
		data.Free()
	} else {
		// Fallback to Data field if GetUncompressedData fails
		dataStr = cp.TypedSpec().Value.Data
	}
	
	patchID := cp.Metadata().ID()
	resp := ConfigPatchResponse{
		ID:        patchID,
		Namespace: cp.Metadata().Namespace(),
		Data:      dataStr,
		Links: map[string]string{
			"self": buildURL(c, "/api/v1/configpatches/"+patchID),
		},
	}

	// Try to find cluster ID from labels
	if clusterID, ok := cp.Metadata().Labels().Get("omni.sidero.dev/cluster"); ok {
		resp.Links["cluster"] = buildURL(c, "/api/v1/clusters/"+clusterID)
	}

	c.JSON(http.StatusOK, resp)
}
