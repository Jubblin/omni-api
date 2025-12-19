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

func TestClusterDiagnosticsHandler_GetClusterDiagnostics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	cd := typed.NewResource[omni.ClusterDiagnosticsSpec, omni.ClusterDiagnosticsExtension](
		resource.NewMetadata("default", omni.ClusterDiagnosticsType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ClusterDiagnosticsSpec{
			Nodes: []*specs.ClusterDiagnosticsSpec_Node{
				{
					Id:             "node-1",
					NumDiagnostics: 2,
				},
			},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(cd, nil)

	handler := NewClusterDiagnosticsHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/diagnostics", nil)

	handler.GetClusterDiagnostics(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ClusterDiagnosticsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cluster-1", resp.ID)
	assert.Len(t, resp.Nodes, 1)
	assert.Equal(t, "node-1", resp.Nodes[0].ID)
	assert.Equal(t, uint32(2), resp.Nodes[0].NumDiagnostics)
}
