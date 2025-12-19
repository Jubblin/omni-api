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

func TestClusterEndpointHandler_GetClusterEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ce := typed.NewResource[omni.ClusterEndpointSpec, omni.ClusterEndpointExtension](
		resource.NewMetadata("default", omni.ClusterEndpointType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterEndpointSpec{
			ManagementAddresses: []string{"192.168.1.10:6443"},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ce, nil)

	handler := NewClusterEndpointHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/endpoints", nil)

	handler.GetClusterEndpoints(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterEndpointResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cluster-1", resp.ID)
	assert.Len(t, resp.ManagementAddresses, 1)
	assert.Equal(t, "192.168.1.10:6443", resp.ManagementAddresses[0])
}
