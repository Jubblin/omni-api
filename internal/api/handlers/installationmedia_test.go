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

func TestInstallationMediaHandler_ListInstallationMedias(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	im := typed.NewResource[omni.InstallationMediaSpec, omni.InstallationMediaExtension](
		resource.NewMetadata("default", omni.InstallationMediaType, "media-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.InstallationMediaSpec{
			Name:         "Talos ISO",
			Architecture: "amd64",
			Profile:      "metal",
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{im},
	}, nil)

	handler := NewInstallationMediaHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/installation-medias", nil)

	handler.ListInstallationMedias(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []InstallationMediaResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "media-1", resp[0].ID)
	assert.Equal(t, "Talos ISO", resp[0].Name)
}

func TestInstallationMediaHandler_GetInstallationMedia(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	im := typed.NewResource[omni.InstallationMediaSpec, omni.InstallationMediaExtension](
		resource.NewMetadata("default", omni.InstallationMediaType, "media-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.InstallationMediaSpec{
			Name: "Talos ISO",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(im, nil)

	handler := NewInstallationMediaHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "media-1"}}
	c.Request, _ = http.NewRequest("GET", "/installation-medias/media-1", nil)

	handler.GetInstallationMedia(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp InstallationMediaResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "media-1", resp.ID)
	assert.Equal(t, "Talos ISO", resp.Name)
}
