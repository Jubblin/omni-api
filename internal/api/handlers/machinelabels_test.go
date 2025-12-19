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

func TestMachineLabelsHandler_GetMachineLabels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ml := typed.NewResource[omni.MachineLabelsSpec, omni.MachineLabelsExtension](
		resource.NewMetadata("default", omni.MachineLabelsType, "machine-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineLabelsSpec{}),
	)
	ml.Metadata().Labels().Set("key1", "value1")
	ml.Metadata().Labels().Set("key2", "value2")

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ml, nil)

	handler := NewMachineLabelsHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machine-1"}}
	c.Request, _ = http.NewRequest("GET", "/machines/machine-1/labels", nil)

	handler.GetMachineLabels(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineLabelsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "machine-1", resp.ID)
	assert.Equal(t, "value1", resp.Labels["key1"])
	assert.Equal(t, "value2", resp.Labels["key2"])
}
