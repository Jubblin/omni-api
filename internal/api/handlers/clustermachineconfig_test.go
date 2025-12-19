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

func TestClusterMachineConfigHandler_GetClusterMachineConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cmc := typed.NewResource[omni.ClusterMachineConfigSpec, omni.ClusterMachineConfigExtension](
		resource.NewMetadata("default", omni.ClusterMachineConfigType, "cm-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterMachineConfigSpec{
			ClusterMachineVersion: "v1",
			WithoutComments:       true,
			GrubUseUkiCmdline:     false,
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cmc, nil)

	handler := NewClusterMachineConfigHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cm-1"}}
	c.Request, _ = http.NewRequest("GET", "/clustermachines/cm-1/config", nil)

	handler.GetClusterMachineConfig(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterMachineConfigResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cm-1", resp.ID)
	assert.Equal(t, "v1", resp.ClusterMachineVersion)
	assert.True(t, resp.WithoutComments)
	assert.False(t, resp.GrubUseUkiCmdline)
}
