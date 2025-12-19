package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/resource/typed"
	"github.com/gin-gonic/gin"
	"github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestControlPlaneStatusHandler_GetControlPlaneStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cps := typed.NewResource[omni.ControlPlaneStatusSpec, omni.ControlPlaneStatusExtension](
		resource.NewMetadata("default", omni.ControlPlaneStatusType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ControlPlaneStatusSpec{
			Conditions: []*specs.ControlPlaneStatusSpec_Condition{
				{
					Type:     specs.ConditionType_UnknownCondition,
					Status:   specs.ControlPlaneStatusSpec_Condition_Unknown,
					Reason:   "AllNodesReady",
					Severity: specs.ControlPlaneStatusSpec_Condition_Info,
				},
			},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cps, nil)

	handler := NewControlPlaneStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/controlplane-status", nil)

	handler.GetControlPlaneStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ControlPlaneStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cluster-1", resp.ID)
	assert.Len(t, resp.Conditions, 1)
	assert.Equal(t, "UnknownCondition", resp.Conditions[0].Type)
	assert.Equal(t, "Unknown", resp.Conditions[0].Status)
	assert.Equal(t, "AllNodesReady", resp.Conditions[0].Reason)
}
