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

func TestKubernetesStatusHandler_GetKubernetesStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ks := typed.NewResource[omni.KubernetesStatusSpec, omni.KubernetesStatusExtension](
		resource.NewMetadata("default", omni.KubernetesStatusType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.KubernetesStatusSpec{
			Nodes: []*specs.KubernetesStatusSpec_NodeStatus{
				{
					Nodename:       "node-1",
					KubeletVersion: "v1.28.0",
					Ready:          true,
				},
			},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ks, nil)

	handler := NewKubernetesStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/kubernetes-status", nil)

	handler.GetKubernetesStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp KubernetesStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cluster-1", resp.ID)
	assert.Len(t, resp.Nodes, 1)
	assert.Equal(t, "node-1", resp.Nodes[0].Nodename)
	assert.True(t, resp.Nodes[0].Ready)
}
