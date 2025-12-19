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

func TestMachineConfigDiffHandler_GetMachineConfigDiff(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	mcd := typed.NewResource[omni.MachineConfigDiffSpec, omni.MachineConfigDiffExtension](
		resource.NewMetadata("default", omni.MachineConfigDiffType, "machine-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineConfigDiffSpec{
			Diff: "- old config\n+ new config",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(mcd, nil)

	handler := NewMachineConfigDiffHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machine-1"}}
	c.Request, _ = http.NewRequest("GET", "/machines/machine-1/config-diff", nil)

	handler.GetMachineConfigDiff(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineConfigDiffResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "machine-1", resp.ID)
	assert.Equal(t, "- old config\n+ new config", resp.Diff)
}
