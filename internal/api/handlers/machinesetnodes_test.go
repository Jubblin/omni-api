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

func TestMachineSetNodeHandler_ListMachineSetNodes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	msn := typed.NewResource[omni.MachineSetNodeSpec, omni.MachineSetNodeExtension](
		resource.NewMetadata("default", omni.MachineSetNodeType, "node-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineSetNodeSpec{}),
	)
	msn.Metadata().Labels().Set("omni.sidero.dev/machine-set", "machineset-1")

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{msn},
	}, nil)

	handler := NewMachineSetNodeHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/machinesetnodes", nil)

	handler.ListMachineSetNodes(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []MachineSetNodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "node-1", resp[0].ID)
	assert.Equal(t, "http://localhost:8080/api/v1/machinesetnodes/node-1", resp[0].Links["self"])
}

func TestMachineSetNodeHandler_GetMachineSetNode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	msn := typed.NewResource[omni.MachineSetNodeSpec, omni.MachineSetNodeExtension](
		resource.NewMetadata("default", omni.MachineSetNodeType, "node-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineSetNodeSpec{}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(msn, nil)

	handler := NewMachineSetNodeHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "node-1"}}
	c.Request, _ = http.NewRequest("GET", "/machinesetnodes/node-1", nil)

	handler.GetMachineSetNode(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineSetNodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "node-1", resp.ID)
	assert.Equal(t, "http://localhost:8080/api/v1/machinesetnodes/node-1", resp.Links["self"])
}
