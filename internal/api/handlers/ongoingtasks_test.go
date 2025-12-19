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

func TestOngoingTaskHandler_ListOngoingTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ot := typed.NewResource[omni.OngoingTaskSpec, omni.OngoingTaskExtension](
		resource.NewMetadata("default", omni.OngoingTaskType, "task-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.OngoingTaskSpec{
			Title:      "Test Task",
			ResourceId: "cluster-1",
		}),
	)

	mockState.On("List", mock.Anything, mock.Anything, mock.Anything).Return(resource.List{
		Items: []resource.Resource{ot},
	}, nil)

	handler := NewOngoingTaskHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/ongoingtasks", nil)

	handler.ListOngoingTasks(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []OngoingTaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "task-1", resp[0].ID)
	assert.Equal(t, "Test Task", resp[0].Title)
	assert.Equal(t, "cluster-1", resp[0].ResourceID)
}

func TestOngoingTaskHandler_GetOngoingTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	ot := typed.NewResource[omni.OngoingTaskSpec, omni.OngoingTaskExtension](
		resource.NewMetadata("default", omni.OngoingTaskType, "task-1", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.OngoingTaskSpec{
			Title:      "Test Task",
			ResourceId: "cluster-1",
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(ot, nil)

	handler := NewOngoingTaskHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "task-1"}}
	c.Request, _ = http.NewRequest("GET", "/ongoingtasks/task-1", nil)

	handler.GetOngoingTask(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp OngoingTaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "task-1", resp.ID)
	assert.Equal(t, "Test Task", resp.Title)
}
