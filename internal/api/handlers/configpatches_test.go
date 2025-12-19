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

func TestConfigPatchHandler_ListConfigPatches(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cp := typed.NewResource[omni.ConfigPatchSpec, omni.ConfigPatchExtension](
		resource.NewMetadata("default", omni.ConfigPatchType, "patch-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ConfigPatchSpec{
			Data: "test data",
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{cp},
	}, nil)

	handler := NewConfigPatchHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/configpatches", nil)

	handler.ListConfigPatches(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []ConfigPatchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "patch-1", resp[0].ID)
	assert.Equal(t, "http://localhost:8080/api/v1/configpatches/patch-1", resp[0].Links["self"])
}

func TestConfigPatchHandler_GetConfigPatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cp := typed.NewResource[omni.ConfigPatchSpec, omni.ConfigPatchExtension](
		resource.NewMetadata("default", omni.ConfigPatchType, "patch-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ConfigPatchSpec{
			Data: "test data",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cp, nil)

	handler := NewConfigPatchHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "patch-1"}}
	c.Request, _ = http.NewRequest("GET", "/configpatches/patch-1", nil)

	handler.GetConfigPatch(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ConfigPatchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "patch-1", resp.ID)
	assert.Equal(t, "http://localhost:8080/api/v1/configpatches/patch-1", resp.Links["self"])
}
