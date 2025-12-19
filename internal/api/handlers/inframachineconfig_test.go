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

func TestInfraMachineConfigHandler_ListInfraMachineConfigs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	imc := typed.NewResource[omni.InfraMachineConfigSpec, omni.InfraMachineConfigExtension](
		resource.NewMetadata("default", omni.InfraMachineConfigType, "machine-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.InfraMachineConfigSpec{
			PowerState:       specs.InfraMachineConfigSpec_POWER_STATE_DEFAULT,
			AcceptanceStatus: specs.InfraMachineConfigSpec_PENDING,
			Cordoned:         false,
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{imc},
	}, nil)

	handler := NewInfraMachineConfigHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/infra-machine-configs", nil)

	handler.ListInfraMachineConfigs(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []InfraMachineConfigResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "machine-1", resp[0].ID)
	assert.False(t, resp[0].Cordoned)
}

func TestInfraMachineConfigHandler_GetInfraMachineConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	imc := typed.NewResource[omni.InfraMachineConfigSpec, omni.InfraMachineConfigExtension](
		resource.NewMetadata("default", omni.InfraMachineConfigType, "machine-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.InfraMachineConfigSpec{
			PowerState: specs.InfraMachineConfigSpec_POWER_STATE_DEFAULT,
			Cordoned:   true,
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(imc, nil)

	handler := NewInfraMachineConfigHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machine-1"}}
	c.Request, _ = http.NewRequest("GET", "/infra-machine-configs/machine-1", nil)

	handler.GetInfraMachineConfig(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp InfraMachineConfigResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "machine-1", resp.ID)
	assert.True(t, resp.Cordoned)
}
