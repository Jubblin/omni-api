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

func TestMachineExtensionsHandler_GetMachineExtensions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	me := typed.NewResource[omni.MachineExtensionsSpec, omni.MachineExtensionsExtension](
		resource.NewMetadata("default", omni.MachineExtensionsType, "machine-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineExtensionsSpec{
			Extensions: []string{"ext1", "ext2"},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(me, nil)

	handler := NewMachineExtensionsHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machine-1"}}
	c.Request, _ = http.NewRequest("GET", "/machines/machine-1/extensions", nil)

	handler.GetMachineExtensions(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineExtensionsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "machine-1", resp.ID)
	assert.Len(t, resp.Extensions, 2)
	assert.Contains(t, resp.Extensions, "ext1")
	assert.Contains(t, resp.Extensions, "ext2")
}
