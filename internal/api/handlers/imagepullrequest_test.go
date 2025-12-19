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

func TestImagePullRequestHandler_ListImagePullRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ipr := typed.NewResource[omni.ImagePullRequestSpec, omni.ImagePullRequestExtension](
		resource.NewMetadata("default", omni.ImagePullRequestType, "request-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ImagePullRequestSpec{
			NodeImageList: []*specs.ImagePullRequestSpec_NodeImageList{
				{
					Node:   "node-1",
					Images: []string{"image1", "image2"},
				},
			},
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{ipr},
	}, nil)

	handler := NewImagePullRequestHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/image-pull-requests", nil)

	handler.ListImagePullRequests(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []ImagePullRequestResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "request-1", resp[0].ID)
	assert.Len(t, resp[0].NodeImageList, 1)
	assert.Equal(t, "node-1", resp[0].NodeImageList[0].Node)
}

func TestImagePullRequestHandler_GetImagePullRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ipr := typed.NewResource[omni.ImagePullRequestSpec, omni.ImagePullRequestExtension](
		resource.NewMetadata("default", omni.ImagePullRequestType, "request-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.ImagePullRequestSpec{
			NodeImageList: []*specs.ImagePullRequestSpec_NodeImageList{
				{
					Node:   "node-1",
					Images: []string{"image1"},
				},
			},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ipr, nil)

	handler := NewImagePullRequestHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "request-1"}}
	c.Request, _ = http.NewRequest("GET", "/image-pull-requests/request-1", nil)

	handler.GetImagePullRequest(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ImagePullRequestResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "request-1", resp.ID)
	assert.Len(t, resp.NodeImageList, 1)
}
