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

func TestEtcdBackupHandler_ListEtcdBackups(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	eb := typed.NewResource[omni.EtcdBackupSpec, omni.EtcdBackupExtension](
		resource.NewMetadata("default", omni.EtcdBackupType, "backup-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.EtcdBackupSpec{
			Snapshot:  "snapshot-1",
			CreatedAt: timestamppb.New(time.Now()),
			Size:      1024,
		}),
	)
	eb.Metadata().Labels().Set("omni.sidero.dev/cluster", "cluster-1")

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{eb},
	}, nil)

	handler := NewEtcdBackupHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/etcdbackups", nil)

	handler.ListEtcdBackups(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []EtcdBackupResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "backup-1", resp[0].ID)
	assert.Equal(t, "snapshot-1", resp[0].Snapshot)
	assert.Equal(t, uint64(1024), resp[0].Size)
}

func TestEtcdBackupHandler_GetEtcdBackup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	eb := typed.NewResource[omni.EtcdBackupSpec, omni.EtcdBackupExtension](
		resource.NewMetadata("default", omni.EtcdBackupType, "backup-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.EtcdBackupSpec{
			Snapshot: "snapshot-1",
			Size:     1024,
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(eb, nil)

	handler := NewEtcdBackupHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "backup-1"}}
	c.Request, _ = http.NewRequest("GET", "/etcdbackups/backup-1", nil)

	handler.GetEtcdBackup(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp EtcdBackupResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "backup-1", resp.ID)
	assert.Equal(t, "snapshot-1", resp.Snapshot)
}
