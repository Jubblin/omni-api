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

func TestClusterHandler_ListClusters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	// Create a dummy cluster
	cl := typed.NewResource[omni.ClusterSpec, omni.ClusterExtension](
		resource.NewMetadata("default", omni.ClusterType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterSpec{
			KubernetesVersion: "v1.28.0",
			TalosVersion:      "v1.5.0",
			Features: &specs.ClusterSpec_Features{
				EnableWorkloadProxy: true,
				DiskEncryption:      false,
			},
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{cl},
	}, nil)

	handler := NewClusterHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/clusters", nil)

	handler.ListClusters(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []ClusterResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "cluster-1", resp[0].ID)
	assert.Equal(t, "v1.28.0", resp[0].KubernetesVersion)
	assert.True(t, resp[0].Features.WorkloadProxy)
	assert.Equal(t, "http://localhost:8080/api/v1/clusters/cluster-1", resp[0].Links["self"])
	assert.Equal(t, "http://localhost:8080/api/v1/clusters/cluster-1/status", resp[0].Links["status"])
}

func TestClusterHandler_GetCluster(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cl := typed.NewResource[omni.ClusterSpec, omni.ClusterExtension](
		resource.NewMetadata("default", omni.ClusterType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterSpec{
			KubernetesVersion: "v1.28.0",
			TalosVersion:      "v1.5.0",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cl, nil)

	handler := NewClusterHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1", nil)

	handler.GetCluster(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cluster-1", resp.ID)
	assert.Equal(t, "http://localhost:8080/api/v1/clusters/cluster-1", resp.Links["self"])
}

func TestClusterHandler_GetClusterStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cs := typed.NewResource[omni.ClusterStatusSpec, omni.ClusterStatusExtension](
		resource.NewMetadata("default", omni.ClusterStatusType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterStatusSpec{
			Available: true,
			Phase:     specs.ClusterStatusSpec_RUNNING,
			Ready:     true,
			Machines: &specs.Machines{
				Total:   3,
				Healthy: 3,
			},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cs, nil)

	handler := NewClusterHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/status", nil)

	handler.GetClusterStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "RUNNING", resp.Phase)
	assert.Equal(t, uint32(3), resp.Machines.Total)
}
