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

func TestImagePullStatusHandler_GetImagePullStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ips := typed.NewResource[omni.ImagePullStatusSpec, omni.ImagePullStatusExtension](
		resource.NewMetadata("default", omni.ImagePullStatusType, "request-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ImagePullStatusSpec{
			LastProcessedNode:  "node-1",
			LastProcessedImage: "image1",
			ProcessedCount:     5,
			TotalCount:         10,
			RequestVersion:     "v1",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ips, nil)

	handler := NewImagePullStatusHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "request-1"}}
	c.Request, _ = http.NewRequest("GET", "/image-pull-requests/request-1/status", nil)

	handler.GetImagePullStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ImagePullStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "request-1", resp.ID)
	assert.Equal(t, "node-1", resp.LastProcessedNode)
	assert.Equal(t, uint32(5), resp.ProcessedCount)
	assert.Equal(t, uint32(10), resp.TotalCount)
}
