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

func TestKernelArgsHandler_ListKernelArgs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ka := typed.NewResource[omni.KernelArgsSpec, omni.KernelArgsExtension](
		resource.NewMetadata("default", omni.KernelArgsType, "args-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.KernelArgsSpec{
			Args: []string{"arg1=value1", "arg2=value2"},
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{ka},
	}, nil)

	handler := NewKernelArgsHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/kernel-args", nil)

	handler.ListKernelArgs(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []KernelArgsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "args-1", resp[0].ID)
	assert.Len(t, resp[0].Args, 2)
}

func TestKernelArgsHandler_GetKernelArgs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ka := typed.NewResource[omni.KernelArgsSpec, omni.KernelArgsExtension](
		resource.NewMetadata("default", omni.KernelArgsType, "args-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.KernelArgsSpec{
			Args: []string{"arg1=value1"},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ka, nil)

	handler := NewKernelArgsHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "args-1"}}
	c.Request, _ = http.NewRequest("GET", "/kernel-args/args-1", nil)

	handler.GetKernelArgs(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp KernelArgsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "args-1", resp.ID)
	assert.Len(t, resp.Args, 1)
}
