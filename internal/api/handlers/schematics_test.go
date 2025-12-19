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

func TestSchematicHandler_ListSchematics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	s := typed.NewResource[omni.SchematicSpec, omni.SchematicExtension](
		resource.NewMetadata("default", omni.SchematicType, "schematic-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.SchematicSpec{}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{s},
	}, nil)

	handler := NewSchematicHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/schematics", nil)

	handler.ListSchematics(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []SchematicResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "schematic-1", resp[0].ID)
	assert.Equal(t, "http://localhost:8080/api/v1/schematics/schematic-1", resp[0].Links["self"])
}

func TestSchematicHandler_GetSchematic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	s := typed.NewResource[omni.SchematicSpec, omni.SchematicExtension](
		resource.NewMetadata("default", omni.SchematicType, "schematic-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.SchematicSpec{}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(s, nil)

	handler := NewSchematicHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "schematic-1"}}
	c.Request, _ = http.NewRequest("GET", "/schematics/schematic-1", nil)

	handler.GetSchematic(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp SchematicResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "schematic-1", resp.ID)
	assert.Equal(t, "http://localhost:8080/api/v1/schematics/schematic-1", resp.Links["self"])
}
