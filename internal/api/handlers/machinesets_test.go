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

func TestMachineSetHandler_ListMachineSets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ms := typed.NewResource[omni.MachineSetSpec, omni.MachineSetExtension](
		resource.NewMetadata("default", omni.MachineSetType, "machineset-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineSetSpec{
			UpdateStrategy: specs.MachineSetSpec_Rolling,
			DeleteStrategy: specs.MachineSetSpec_Unset,
			MachineAllocation: &specs.MachineSetSpec_MachineAllocation{
				Name: "class-1",
			},
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{ms},
	}, nil)

	handler := NewMachineSetHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/machinesets", nil)

	handler.ListMachineSets(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []MachineSetResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "machineset-1", resp[0].ID)
	assert.Equal(t, "class-1", resp[0].MachineClass)
	assert.Equal(t, "http://localhost:8080/api/v1/machinesets/machineset-1", resp[0].Links["self"])
}

func TestMachineSetHandler_GetMachineSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ms := typed.NewResource[omni.MachineSetSpec, omni.MachineSetExtension](
		resource.NewMetadata("default", omni.MachineSetType, "machineset-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineSetSpec{
			UpdateStrategy: specs.MachineSetSpec_Rolling,
			DeleteStrategy: specs.MachineSetSpec_Unset,
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ms, nil)

	handler := NewMachineSetHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machineset-1"}}
	c.Request, _ = http.NewRequest("GET", "/machinesets/machineset-1", nil)

	handler.GetMachineSet(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineSetResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "machineset-1", resp.ID)
	assert.Equal(t, "http://localhost:8080/api/v1/machinesets/machineset-1", resp.Links["self"])
}

func TestMachineSetStatusHandler_GetMachineSetStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	mss := typed.NewResource[omni.MachineSetStatusSpec, omni.MachineSetStatusExtension](
		resource.NewMetadata("default", omni.MachineSetStatusType, "machineset-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineSetStatusSpec{
			Machines: &specs.Machines{
				Total:   3,
				Healthy: 2,
			},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(mss, nil)

	handler := NewMachineSetStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machineset-1"}}
	c.Request, _ = http.NewRequest("GET", "/machinesets/machineset-1/status", nil)

	handler.GetMachineSetStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineSetStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "machineset-1", resp.ID)
	assert.Equal(t, uint32(3), resp.Machines.Total)
	assert.Equal(t, uint32(2), resp.Machines.Healthy)
}
