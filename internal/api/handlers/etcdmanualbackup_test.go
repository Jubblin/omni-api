package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/resource/typed"
	"github.com/gin-gonic/gin"
	"github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEtcdManualBackupHandler_ListEtcdManualBackups(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	emb := typed.NewResource[omni.EtcdManualBackupSpec, omni.EtcdManualBackupExtension](
		resource.NewMetadata("default", omni.EtcdManualBackupType, "backup-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.EtcdManualBackupSpec{
			BackupAt: timestamppb.New(time.Now()),
		}),
	)
	emb.Metadata().Labels().Set("omni.sidero.dev/cluster", "cluster-1")

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{emb},
	}, nil)

	handler := NewEtcdManualBackupHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/etcd-manual-backups", nil)

	handler.ListEtcdManualBackups(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []EtcdManualBackupResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "backup-1", resp[0].ID)
}

func TestEtcdManualBackupHandler_GetEtcdManualBackup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	emb := typed.NewResource[omni.EtcdManualBackupSpec, omni.EtcdManualBackupExtension](
		resource.NewMetadata("default", omni.EtcdManualBackupType, "backup-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.EtcdManualBackupSpec{
			BackupAt: timestamppb.New(time.Now()),
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(emb, nil)

	handler := NewEtcdManualBackupHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "backup-1"}}
	c.Request, _ = http.NewRequest("GET", "/etcd-manual-backups/backup-1", nil)

	handler.GetEtcdManualBackup(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp EtcdManualBackupResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "backup-1", resp.ID)
}
