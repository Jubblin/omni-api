package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	omniresources "github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
)

func TestHealthHandler_GetHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		omniConnected  bool
		expectedStatus string
	}{
		{
			name:           "Omni connected",
			omniConnected:  true,
			expectedStatus: "healthy",
		},
		{
			name:           "Omni not connected",
			omniConnected:  false,
			expectedStatus: "degraded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockState := new(MockState)
			handler := NewHealthHandler(mockState)

			md := resource.NewMetadata(omniresources.DefaultNamespace, omni.ClusterType, "", resource.VersionUndefined)
			if tt.omniConnected {
				mockState.On("List", mock.Anything, mock.MatchedBy(func(k resource.Kind) bool {
					return k.Namespace() == md.Namespace() && k.Type() == md.Type()
				}), mock.Anything).Return(resource.List{Items: []resource.Resource{}}, nil)
			} else {
				mockState.On("List", mock.Anything, mock.MatchedBy(func(k resource.Kind) bool {
					return k.Namespace() == md.Namespace() && k.Type() == md.Type()
				}), mock.Anything).Return(resource.List{}, assert.AnError)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/health", nil)

			handler.GetHealth(c)

			assert.Equal(t, http.StatusOK, w.Code)

			var resp HealthResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.Status)
			assert.Equal(t, tt.omniConnected, resp.Omni.Connected)
			assert.NotEmpty(t, resp.Timestamp)
			assert.Equal(t, "0.0.1", resp.Version)
			assert.NotEmpty(t, resp.Links["self"])
			assert.NotEmpty(t, resp.Links["metrics"])

			mockState.AssertExpectations(t)
		})
	}
}
