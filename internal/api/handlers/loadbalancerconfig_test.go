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

func TestLoadBalancerConfigHandler_ListLoadBalancerConfigs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	lb := typed.NewResource[omni.LoadBalancerConfigSpec, omni.LoadBalancerConfigExtension](
		resource.NewMetadata("default", omni.LoadBalancerConfigType, "lb-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.LoadBalancerConfigSpec{
			BindPort:          "6443",
			SiderolinkEndpoint: "https://lb.example.com",
			Endpoints:         []string{"192.168.1.10:6443"},
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{lb},
	}, nil)

	handler := NewLoadBalancerConfigHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/loadbalancer-configs", nil)

	handler.ListLoadBalancerConfigs(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []LoadBalancerConfigResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "lb-1", resp[0].ID)
	assert.Equal(t, "6443", resp[0].BindPort)
}

func TestLoadBalancerConfigHandler_GetLoadBalancerConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	lb := typed.NewResource[omni.LoadBalancerConfigSpec, omni.LoadBalancerConfigExtension](
		resource.NewMetadata("default", omni.LoadBalancerConfigType, "lb-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.LoadBalancerConfigSpec{
			BindPort: "6443",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(lb, nil)

	handler := NewLoadBalancerConfigHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "lb-1"}}
	c.Request, _ = http.NewRequest("GET", "/loadbalancer-configs/lb-1", nil)

	handler.GetLoadBalancerConfig(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp LoadBalancerConfigResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "lb-1", resp.ID)
	assert.Equal(t, "6443", resp.BindPort)
}
