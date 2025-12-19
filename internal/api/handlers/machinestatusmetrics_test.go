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

func TestMachineStatusMetricsHandler_GetMachineStatusMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockState := new(MockState)

	msm := typed.NewResource[omni.MachineStatusMetricsSpec, omni.MachineStatusMetricsExtension](
		resource.NewMetadata("default", omni.MachineStatusMetricsType, "default", resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineStatusMetricsSpec{
			RegisteredMachinesCount: 10,
			ConnectedMachinesCount:  8,
			AllocatedMachinesCount: 5,
			PendingMachinesCount:   2,
			Platforms: map[string]uint32{
				"metal": 8,
				"cloud": 2,
			},
		}),
	)

	mockState.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(msm, nil)

	handler := NewMachineStatusMetricsHandler(mockState)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "machine-1"}}
	c.Request, _ = http.NewRequest("GET", "/machines/machine-1/metrics", nil)

	handler.GetMachineStatusMetrics(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MachineStatusMetricsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), resp.RegisteredMachinesCount)
	assert.Equal(t, uint32(8), resp.ConnectedMachinesCount)
	assert.Equal(t, uint32(5), resp.AllocatedMachinesCount)
	assert.Equal(t, uint32(2), resp.PendingMachinesCount)
	assert.Equal(t, uint32(8), resp.Platforms["metal"])
}
