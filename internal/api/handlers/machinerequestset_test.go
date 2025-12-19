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

func TestMachineRequestSetHandler_ListMachineRequestSets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	mrs := typed.NewResource[omni.MachineRequestSetSpec, omni.MachineRequestSetExtension](
		resource.NewMetadata("default", omni.MachineRequestSetType, "requestset-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineRequestSetSpec{
			ProviderId:   "provider-1",
			MachineCount: 3,
			TalosVersion: "v1.5.0",
			Extensions:   []string{"ext1"},
			KernelArgs:   []string{"arg1"},
			GrpcTunnel:   specs.GrpcTunnelMode_ENABLED,
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{mrs},
	}, nil)

	handler := NewMachineRequestSetHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/machine-request-sets", nil)

	handler.ListMachineRequestSets(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []MachineRequestSetResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "requestset-1", resp[0].ID)
	assert.Equal(t, "provider-1", resp[0].ProviderID)
	assert.Equal(t, int32(3), resp[0].MachineCount)
}

func TestMachineRequestSetHandler_GetMachineRequestSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	mrs := typed.NewResource[omni.MachineRequestSetSpec, omni.MachineRequestSetExtension](
		resource.NewMetadata("default", omni.MachineRequestSetType, "requestset-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineRequestSetSpec{
			ProviderId:   "provider-1",
			MachineCount: 3,
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(mrs, nil)

	handler := NewMachineRequestSetHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "requestset-1"}}
	c.Request, _ = http.NewRequest("GET", "/machine-request-sets/requestset-1", nil)

	handler.GetMachineRequestSet(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineRequestSetResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "requestset-1", resp.ID)
	assert.Equal(t, "provider-1", resp.ProviderID)
}
