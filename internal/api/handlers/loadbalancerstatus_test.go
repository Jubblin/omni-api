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

func TestLoadBalancerStatusHandler_GetLoadBalancerStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	lbs := typed.NewResource[omni.LoadBalancerStatusSpec, omni.LoadBalancerStatusExtension](
		resource.NewMetadata("default", omni.LoadBalancerStatusType, "lb-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.LoadBalancerStatusSpec{
			Healthy: true,
			Stopped: false,
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(lbs, nil)

	handler := NewLoadBalancerStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "lb-1"}}
	c.Request, _ = http.NewRequest("GET", "/loadbalancers/lb-1/status", nil)

	handler.GetLoadBalancerStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp LoadBalancerStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "lb-1", resp.ID)
	assert.True(t, resp.Healthy)
	assert.False(t, resp.Stopped)
}
