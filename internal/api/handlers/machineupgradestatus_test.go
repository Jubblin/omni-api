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

func TestMachineUpgradeStatusHandler_GetMachineUpgradeStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	mus := typed.NewResource[omni.MachineUpgradeStatusSpec, omni.MachineUpgradeStatusExtension](
		resource.NewMetadata("default", omni.MachineUpgradeStatusType, "machine-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineUpgradeStatusSpec{
			Phase:              specs.MachineUpgradeStatusSpec_Upgrading,
			CurrentTalosVersion: "v1.4.0",
			TalosVersion:        "v1.5.0",
			CurrentSchematicId:  "schematic-1",
			SchematicId:         "schematic-2",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(mus, nil)

	handler := NewMachineUpgradeStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machine-1"}}
	c.Request, _ = http.NewRequest("GET", "/machines/machine-1/upgrade-status", nil)

	handler.GetMachineUpgradeStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineUpgradeStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "machine-1", resp.ID)
	assert.Equal(t, "Upgrading", resp.Phase)
	assert.Equal(t, "v1.4.0", resp.CurrentTalosVersion)
	assert.Equal(t, "v1.5.0", resp.TalosVersion)
	assert.Equal(t, "schematic-1", resp.CurrentSchematicID)
	assert.Equal(t, "schematic-2", resp.SchematicID)
}
