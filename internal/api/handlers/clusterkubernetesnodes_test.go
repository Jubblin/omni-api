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

func TestClusterKubernetesNodesHandler_ListClusterKubernetesNodes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ckn := typed.NewResource[omni.ClusterKubernetesNodesSpec, omni.ClusterKubernetesNodesExtension](
		resource.NewMetadata("default", omni.ClusterKubernetesNodesType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterKubernetesNodesSpec{
			Nodes: []string{"node-1", "node-2"},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ckn, nil)

	handler := NewClusterKubernetesNodesHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/kubernetes-nodes", nil)

	handler.ListClusterKubernetesNodes(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []ClusterKubernetesNodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "node-1", resp[0].ID)
	assert.Equal(t, "node-2", resp[1].ID)
}

func TestClusterKubernetesNodesHandler_GetClusterKubernetesNode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ckn := typed.NewResource[omni.ClusterKubernetesNodesSpec, omni.ClusterKubernetesNodesExtension](
		resource.NewMetadata("default", omni.ClusterKubernetesNodesType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterKubernetesNodesSpec{
			Nodes: []string{"node-1"},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ckn, nil)

	handler := NewClusterKubernetesNodesHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}, {Key: "node", Value: "node-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/kubernetes-nodes/node-1", nil)

	handler.GetClusterKubernetesNode(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterKubernetesNodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "node-1", resp.ID)
	assert.Equal(t, "node-1", resp.Name)
}
