package handlers

import (
	"encoding/json"
	"errors"
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

func TestMachineHandler_ListMachines(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	m := typed.NewResource[omni.MachineSpec, omni.MachineExtension](
		resource.NewMetadata("default", omni.MachineType, "machine-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineSpec{
			ManagementAddress: "192.168.1.10",
			Connected:         true,
			UseGrpcTunnel:     true,
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{m},
	}, nil)

	// Mock MachineStatus Get call (may be called when status consolidation is attempted)
	// Return an error to simulate status not being available
	mockState.On("Get", mock.Anything, mock.MatchedBy(func(md resource.Pointer) bool {
		return md.Type() == omni.MachineStatusType
	}), mock.Anything).Return(nil, errors.New("not found"))

	handler := NewMachineHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/machines", nil)

	handler.ListMachines(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []MachineResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "machine-1", resp[0].ID)
	assert.Equal(t, "192.168.1.10", resp[0].ManagementAddress)
	assert.True(t, resp[0].UseGrpcTunnel)
	assert.Equal(t, "http://localhost:8080/api/v1/machines/machine-1", resp[0].Links["self"])
}

func TestMachineHandler_GetMachine(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	m := typed.NewResource[omni.MachineSpec, omni.MachineExtension](
		resource.NewMetadata("default", omni.MachineType, "machine-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineSpec{
			ManagementAddress: "192.168.1.10",
			Connected:         true,
		}),
	)

	mockState.On("Get", mock.Anything, mock.MatchedBy(func(md resource.Pointer) bool {
		return md.Type() == omni.MachineType
	}), mock.Anything).Return(m, nil)

	// Mock MachineStatus Get call (may be called when status consolidation is attempted)
	// Return an error to simulate status not being available
	mockState.On("Get", mock.Anything, mock.MatchedBy(func(md resource.Pointer) bool {
		return md.Type() == omni.MachineStatusType
	}), mock.Anything).Return(nil, errors.New("not found"))

	handler := NewMachineHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machine-1"}}
	c.Request, _ = http.NewRequest("GET", "/machines/machine-1", nil)

	handler.GetMachine(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "machine-1", resp.ID)
	assert.Equal(t, "http://localhost:8080/api/v1/machines/machine-1", resp.Links["self"])
}
