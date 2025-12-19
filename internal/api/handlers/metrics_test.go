package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMetricsHandler_GetMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Reset metrics
	metricsMutex.Lock()
	requestCounts = make(map[string]uint64)
	responseTimes = make(map[string][]float64)
	errorCounts = make(map[string]uint64)
	metricsMutex.Unlock()

	// Record some test data
	RecordRequest("/api/v1/clusters", 100*time.Millisecond, http.StatusOK)
	RecordRequest("/api/v1/clusters", 150*time.Millisecond, http.StatusOK)
	RecordRequest("/api/v1/machines", 200*time.Millisecond, http.StatusNotFound)

	handler := NewMetricsHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/metrics", nil)

	handler.GetMetrics(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp MetricsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Greater(t, resp.ServerUptimeSeconds, 0.0)
	assert.Equal(t, uint64(2), resp.RequestCounts["/api/v1/clusters"])
	assert.Equal(t, uint64(1), resp.RequestCounts["/api/v1/machines"])
	assert.Equal(t, uint64(1), resp.ErrorCounts["/api/v1/machines"])
	assert.NotEmpty(t, resp.Links["self"])
	assert.NotEmpty(t, resp.Links["health"])
}

func TestRecordRequest(t *testing.T) {
	// Reset metrics
	metricsMutex.Lock()
	requestCounts = make(map[string]uint64)
	responseTimes = make(map[string][]float64)
	errorCounts = make(map[string]uint64)
	metricsMutex.Unlock()

	// Record requests
	RecordRequest("/test", 50*time.Millisecond, http.StatusOK)
	RecordRequest("/test", 100*time.Millisecond, http.StatusOK)
	RecordRequest("/test", 150*time.Millisecond, http.StatusInternalServerError)

	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	assert.Equal(t, uint64(3), requestCounts["/test"])
	assert.Equal(t, uint64(1), errorCounts["/test"])
	assert.Len(t, responseTimes["/test"], 3)
}
