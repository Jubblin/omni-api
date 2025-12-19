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

func TestEtcdBackupStatusHandler_GetEtcdBackupStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ebs := typed.NewResource[omni.EtcdBackupStatusSpec, omni.EtcdBackupStatusExtension](
		resource.NewMetadata("default", omni.EtcdBackupStatusType, "backup-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.EtcdBackupStatusSpec{
			Status: specs.EtcdBackupStatusSpec_Unknown,
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ebs, nil)

	handler := NewEtcdBackupStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "backup-1"}}
	c.Request, _ = http.NewRequest("GET", "/etcdbackups/backup-1/status", nil)

	handler.GetEtcdBackupStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp EtcdBackupStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "backup-1", resp.ID)
	assert.Equal(t, "Unknown", resp.Status)
}
