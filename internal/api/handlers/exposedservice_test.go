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

func TestExposedServiceHandler_ListExposedServices(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	es := typed.NewResource[omni.ExposedServiceSpec, omni.ExposedServiceExtension](
		resource.NewMetadata("default", omni.ExposedServiceType, "service-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ExposedServiceSpec{
			Port:  8080,
			Label: "Test Service",
			Url:   "http://example.com",
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{es},
	}, nil)

	handler := NewExposedServiceHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/exposed-services", nil)

	handler.ListExposedServices(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []ExposedServiceResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "service-1", resp[0].ID)
	assert.Equal(t, uint32(8080), resp[0].Port)
	assert.Equal(t, "Test Service", resp[0].Label)
}

func TestExposedServiceHandler_GetExposedService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	es := typed.NewResource[omni.ExposedServiceSpec, omni.ExposedServiceExtension](
		resource.NewMetadata("default", omni.ExposedServiceType, "service-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ExposedServiceSpec{
			Port: 8080,
			Url:  "http://example.com",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(es, nil)

	handler := NewExposedServiceHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "service-1"}}
	c.Request, _ = http.NewRequest("GET", "/exposed-services/service-1", nil)

	handler.GetExposedService(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ExposedServiceResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "service-1", resp.ID)
	assert.Equal(t, uint32(8080), resp.Port)
}
