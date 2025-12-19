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

func TestKubernetesUpgradeHandler_GetKubernetesUpgradeStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	kus := typed.NewResource[omni.KubernetesUpgradeStatusSpec, omni.KubernetesUpgradeStatusExtension](
		resource.NewMetadata("default", omni.KubernetesUpgradeStatusType, "cluster-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.KubernetesUpgradeStatusSpec{
			Phase:                specs.KubernetesUpgradeStatusSpec_Upgrading,
			CurrentUpgradeVersion: "v1.29.0",
			LastUpgradeVersion:    "v1.28.0",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(kus, nil)

	handler := NewKubernetesUpgradeHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "cluster-1"}}
	c.Request, _ = http.NewRequest("GET", "/clusters/cluster-1/kubernetes-upgrade", nil)

	handler.GetKubernetesUpgradeStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp KubernetesUpgradeStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cluster-1", resp.ID)
	assert.Equal(t, "Upgrading", resp.Phase)
	assert.Equal(t, "v1.28.0", resp.LastUpgradeVersion)
	assert.Equal(t, "v1.29.0", resp.CurrentUpgradeVersion)
}
