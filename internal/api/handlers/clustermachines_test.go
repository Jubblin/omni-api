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

func TestClusterMachineHandler_ListClusterMachines(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cm := typed.NewResource[omni.ClusterMachineSpec, omni.ClusterMachineExtension](
		resource.NewMetadata("default", omni.ClusterMachineType, "cm-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterMachineSpec{
			KubernetesVersion: "v1.28.0",
		}),
	)
	cm.Metadata().Labels().Set("omni.sidero.dev/cluster", "cluster-1")

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{cm},
	}, nil)

	handler := NewClusterMachineHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/clustermachines", nil)

	handler.ListClusterMachines(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []ClusterMachineResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "cm-1", resp[0].ID)
	assert.Equal(t, "v1.28.0", resp[0].KubernetesVersion)
	assert.Equal(t, "http://localhost:8080/api/v1/clustermachines/cm-1", resp[0].Links["self"])
}

func TestClusterMachineHandler_ListClusterMachines_WithFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cm1 := typed.NewResource[omni.ClusterMachineSpec, omni.ClusterMachineExtension](
		resource.NewMetadata("default", omni.ClusterMachineType, "cm-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterMachineSpec{}),
	)
	cm1.Metadata().Labels().Set("omni.sidero.dev/cluster", "cluster-1")

	cm2 := typed.NewResource[omni.ClusterMachineSpec, omni.ClusterMachineExtension](
		resource.NewMetadata("default", omni.ClusterMachineType, "cm-2", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterMachineSpec{}),
	)
	cm2.Metadata().Labels().Set("omni.sidero.dev/cluster", "cluster-2")

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{cm1, cm2},
	}, nil)

	handler := NewClusterMachineHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/clustermachines?cluster=cluster-1", nil)

	handler.ListClusterMachines(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []ClusterMachineResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "cm-1", resp[0].ID)
}

func TestClusterMachineHandler_GetClusterMachine(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cm := typed.NewResource[omni.ClusterMachineSpec, omni.ClusterMachineExtension](
		resource.NewMetadata("default", omni.ClusterMachineType, "cm-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterMachineSpec{
			KubernetesVersion: "v1.28.0",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cm, nil)

	handler := NewClusterMachineHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cm-1"}}
	c.Request, _ = http.NewRequest("GET", "/clustermachines/cm-1", nil)

	handler.GetClusterMachine(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterMachineResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cm-1", resp.ID)
	assert.Equal(t, "v1.28.0", resp.KubernetesVersion)
	assert.Equal(t, "http://localhost:8080/api/v1/clustermachines/cm-1", resp.Links["self"])
}

func TestClusterMachineStatusHandler_GetClusterMachineStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cms := typed.NewResource[omni.ClusterMachineStatusSpec, omni.ClusterMachineStatusExtension](
		resource.NewMetadata("default", omni.ClusterMachineStatusType, "cm-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterMachineStatusSpec{
			Ready:          true,
			Stage:          specs.ClusterMachineStatusSpec_RUNNING,
			ConfigUpToDate: true,
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cms, nil)

	handler := NewClusterMachineStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cm-1"}}
	c.Request, _ = http.NewRequest("GET", "/clustermachines/cm-1/status", nil)

	handler.GetClusterMachineStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterMachineStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cm-1", resp.ID)
	assert.True(t, resp.Ready)
	assert.True(t, resp.ConfigUpToDate)
}

func TestClusterMachineConfigStatusHandler_GetClusterMachineConfigStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cmcs := typed.NewResource[omni.ClusterMachineConfigStatusSpec, omni.ClusterMachineConfigStatusExtension](
		resource.NewMetadata("default", omni.ClusterMachineConfigStatusType, "cm-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterMachineConfigStatusSpec{
			TalosVersion: "v1.5.0",
			SchematicId:  "schematic-1",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cmcs, nil)

	handler := NewClusterMachineConfigStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cm-1"}}
	c.Request, _ = http.NewRequest("GET", "/clustermachines/cm-1/config-status", nil)

	handler.GetClusterMachineConfigStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterMachineConfigStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cm-1", resp.ID)
	assert.Equal(t, "v1.5.0", resp.TalosVersion)
	assert.Equal(t, "schematic-1", resp.SchematicID)
}

func TestClusterMachineTalosVersionHandler_GetClusterMachineTalosVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cmtv := typed.NewResource[omni.ClusterMachineTalosVersionSpec, omni.ClusterMachineTalosVersionExtension](
		resource.NewMetadata("default", omni.ClusterMachineTalosVersionType, "cm-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterMachineTalosVersionSpec{
			TalosVersion: "v1.5.0",
			SchematicId:  "schematic-1",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cmtv, nil)

	handler := NewClusterMachineTalosVersionHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cm-1"}}
	c.Request, _ = http.NewRequest("GET", "/clustermachines/cm-1/talos-version", nil)

	handler.GetClusterMachineTalosVersion(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterMachineTalosVersionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cm-1", resp.ID)
	assert.Equal(t, "v1.5.0", resp.TalosVersion)
	assert.Equal(t, "schematic-1", resp.SchematicID)
}
