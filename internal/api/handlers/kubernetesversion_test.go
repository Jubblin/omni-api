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

func TestKubernetesVersionHandler_ListKubernetesVersions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	kv := typed.NewResource[omni.KubernetesVersionSpec, omni.KubernetesVersionExtension](
		resource.NewMetadata("default", omni.KubernetesVersionType, "v1.28.0", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.KubernetesVersionSpec{
			Version: "v1.28.0",
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{kv},
	}, nil)

	handler := NewKubernetesVersionHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/kubernetes-versions", nil)

	handler.ListKubernetesVersions(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []KubernetesVersionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "v1.28.0", resp[0].ID)
	assert.Equal(t, "v1.28.0", resp[0].Version)
}

func TestKubernetesVersionHandler_GetKubernetesVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	kv := typed.NewResource[omni.KubernetesVersionSpec, omni.KubernetesVersionExtension](
		resource.NewMetadata("default", omni.KubernetesVersionType, "v1.28.0", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.KubernetesVersionSpec{
			Version: "v1.28.0",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(kv, nil)

	handler := NewKubernetesVersionHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "v1.28.0"}}
	c.Request, _ = http.NewRequest("GET", "/kubernetes-versions/v1.28.0", nil)

	handler.GetKubernetesVersion(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp KubernetesVersionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "v1.28.0", resp.ID)
	assert.Equal(t, "v1.28.0", resp.Version)
}
