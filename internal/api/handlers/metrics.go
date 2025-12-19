package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	metricsMutex     sync.RWMutex
	requestCounts    = make(map[string]uint64)
	responseTimes    = make(map[string][]float64)
	errorCounts      = make(map[string]uint64)
	serverStartTime  = time.Now()
)

func init() {
	serverStartTime = time.Now()
}

// MetricsResponse represents Prometheus-style metrics
type MetricsResponse struct {
	ServerUptimeSeconds float64            `json:"server_uptime_seconds"`
	RequestCounts       map[string]uint64   `json:"request_counts,omitempty"`
	ErrorCounts         map[string]uint64   `json:"error_counts,omitempty"`
	AverageResponseTime map[string]float64  `json:"average_response_time_seconds,omitempty"`
	Links               map[string]string   `json:"_links,omitempty"`
}

// MetricsHandler handles metrics requests
type MetricsHandler struct{}

// NewMetricsHandler creates a new MetricsHandler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// RecordRequest records a request for metrics
func RecordRequest(endpoint string, duration time.Duration, statusCode int) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	requestCounts[endpoint]++
	if statusCode >= 400 {
		errorCounts[endpoint]++
	}

	// Keep only last 100 response times per endpoint for average calculation
	if responseTimes[endpoint] == nil {
		responseTimes[endpoint] = make([]float64, 0, 100)
	}
	responseTimes[endpoint] = append(responseTimes[endpoint], duration.Seconds())
	if len(responseTimes[endpoint]) > 100 {
		responseTimes[endpoint] = responseTimes[endpoint][1:]
	}
}

// GetMetrics godoc
// @Summary      Get API metrics
// @Description  Get Prometheus-style metrics for the API server
// @Tags         metrics
// @Produce      json
// @Success      200  {object}  MetricsResponse
// @Router       /metrics [get]
func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	uptime := time.Since(serverStartTime).Seconds()

	// Calculate average response times
	avgResponseTimes := make(map[string]float64)
	for endpoint, times := range responseTimes {
		if len(times) > 0 {
			var sum float64
			for _, t := range times {
				sum += t
			}
			avgResponseTimes[endpoint] = sum / float64(len(times))
		}
	}

	// Create copies of maps to avoid race conditions
	requestCountsCopy := make(map[string]uint64)
	for k, v := range requestCounts {
		requestCountsCopy[k] = v
	}

	errorCountsCopy := make(map[string]uint64)
	for k, v := range errorCounts {
		errorCountsCopy[k] = v
	}

	resp := MetricsResponse{
		ServerUptimeSeconds: uptime,
		RequestCounts:       requestCountsCopy,
		ErrorCounts:         errorCountsCopy,
		AverageResponseTime: avgResponseTimes,
		Links: map[string]string{
			"self":   buildURL(c, "/metrics"),
			"health": buildURL(c, "/health"),
		},
	}

	c.JSON(http.StatusOK, resp)
}
