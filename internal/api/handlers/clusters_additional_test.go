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

func TestClusterHandler_GetClusterMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cm := typed.NewResource[omni.ClusterMetricsSpec, omni.ClusterMetricsExtension](
		resource.NewMetadata("default", omni.ClusterMetricsType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterMetricsSpec{
			Features: map[string]uint32{
				"feature1": 1,
				"feature2": 2,
			},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cm, nil)

	handler := NewClusterHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/metrics", nil)

	handler.GetClusterMetrics(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterMetricsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), resp.Features["feature1"])
	assert.Equal(t, uint32(2), resp.Features["feature2"])
}

func TestClusterHandler_GetClusterBootstrap(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cbs := typed.NewResource[omni.ClusterBootstrapStatusSpec, omni.ClusterBootstrapStatusExtension](
		resource.NewMetadata("default", omni.ClusterBootstrapStatusType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterBootstrapStatusSpec{
			Bootstrapped: true,
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cbs, nil)

	handler := NewClusterHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/bootstrap", nil)

	handler.GetClusterBootstrap(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterBootstrapResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Bootstrapped)
}
