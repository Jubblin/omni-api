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

func TestMachineClassHandler_ListMachineClasses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	mc := typed.NewResource[omni.MachineClassSpec, omni.MachineClassExtension](
		resource.NewMetadata("default", omni.MachineClassType, "class-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineClassSpec{
			MatchLabels: []string{"label1=value1"},
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{mc},
	}, nil)

	handler := NewMachineClassHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/machineclasses", nil)

	handler.ListMachineClasses(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []MachineClassResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "class-1", resp[0].ID)
	assert.Equal(t, "http://localhost:8080/api/v1/machineclasses/class-1", resp[0].Links["self"])
}

func TestMachineClassHandler_GetMachineClass(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	mc := typed.NewResource[omni.MachineClassSpec, omni.MachineClassExtension](
		resource.NewMetadata("default", omni.MachineClassType, "class-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineClassSpec{
			MatchLabels: []string{"label1=value1"},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(mc, nil)

	handler := NewMachineClassHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "class-1"}}
	c.Request, _ = http.NewRequest("GET", "/machineclasses/class-1", nil)

	handler.GetMachineClass(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineClassResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "class-1", resp.ID)
	assert.Equal(t, "http://localhost:8080/api/v1/machineclasses/class-1", resp.Links["self"])
}
