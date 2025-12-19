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

func TestSchematicConfigurationHandler_ListSchematicConfigurations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	sc := typed.NewResource[omni.SchematicConfigurationSpec, omni.SchematicConfigurationExtension](
		resource.NewMetadata("default", omni.SchematicConfigurationType, "config-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.SchematicConfigurationSpec{
			SchematicId:  "schematic-1",
			TalosVersion: "v1.5.0",
			KernelArgs:   []string{"arg1", "arg2"},
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{sc},
	}, nil)

	handler := NewSchematicConfigurationHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/schematic-configurations", nil)

	handler.ListSchematicConfigurations(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []SchematicConfigurationResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "config-1", resp[0].ID)
	assert.Equal(t, "schematic-1", resp[0].SchematicID)
}

func TestSchematicConfigurationHandler_GetSchematicConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	sc := typed.NewResource[omni.SchematicConfigurationSpec, omni.SchematicConfigurationExtension](
		resource.NewMetadata("default", omni.SchematicConfigurationType, "config-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.SchematicConfigurationSpec{
			SchematicId:  "schematic-1",
			TalosVersion: "v1.5.0",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sc, nil)

	handler := NewSchematicConfigurationHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "config-1"}}
	c.Request, _ = http.NewRequest("GET", "/schematic-configurations/config-1", nil)

	handler.GetSchematicConfiguration(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp SchematicConfigurationResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "config-1", resp.ID)
	assert.Equal(t, "schematic-1", resp.SchematicID)
}
