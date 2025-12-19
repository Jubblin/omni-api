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

func TestMachineSetDestroyStatusHandler_GetMachineSetDestroyStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	msds := typed.NewResource[omni.MachineSetDestroyStatusSpec, omni.MachineSetDestroyStatusExtension](
		resource.NewMetadata("default", omni.MachineSetDestroyStatusType, "machineset-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.DestroyStatusSpec{
			Phase: "destroying",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(msds, nil)

	handler := NewMachineSetDestroyStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machineset-1"}}
	c.Request, _ = http.NewRequest("GET", "/machinesets/machineset-1/destroy-status", nil)

	handler.GetMachineSetDestroyStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineSetDestroyStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "machineset-1", resp.ID)
	assert.Equal(t, "destroying", resp.Phase)
}
