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

func TestExtensionsConfigurationHandler_ListExtensionsConfigurations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ec := typed.NewResource[omni.ExtensionsConfigurationSpec, omni.ExtensionsConfigurationExtension](
		resource.NewMetadata("default", omni.ExtensionsConfigurationType, "config-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ExtensionsConfigurationSpec{
			Extensions: []string{"ext1", "ext2"},
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{ec},
	}, nil)

	handler := NewExtensionsConfigurationHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/extensions-configurations", nil)

	handler.ListExtensionsConfigurations(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []ExtensionsConfigurationResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "config-1", resp[0].ID)
	assert.Len(t, resp[0].Extensions, 2)
}

func TestExtensionsConfigurationHandler_GetExtensionsConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ec := typed.NewResource[omni.ExtensionsConfigurationSpec, omni.ExtensionsConfigurationExtension](
		resource.NewMetadata("default", omni.ExtensionsConfigurationType, "config-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ExtensionsConfigurationSpec{
			Extensions: []string{"ext1"},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ec, nil)

	handler := NewExtensionsConfigurationHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "config-1"}}
	c.Request, _ = http.NewRequest("GET", "/extensions-configurations/config-1", nil)

	handler.GetExtensionsConfiguration(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ExtensionsConfigurationResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "config-1", resp.ID)
	assert.Len(t, resp.Extensions, 1)
}
