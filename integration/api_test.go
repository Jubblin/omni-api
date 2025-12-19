package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/richardw/talos-ctl/internal/api/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testServer wraps the API server for testing
type testServer struct {
	baseURL string
	client  *http.Client
}

// setupTestServer creates and starts a test API server
func setupTestServer(t *testing.T) *testServer {
	// Check if integration tests should run
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests. Set INTEGRATION_TESTS=true to run.")
	}

	// Verify Omni endpoint is configured
	endpoint := os.Getenv("OMNI_ENDPOINT")
	require.NotEmpty(t, endpoint, "OMNI_ENDPOINT must be set for integration tests")

	// Validate Omni endpoint is accessible (optional validation)
	// The actual API server will handle Omni client initialization

	// Get port from environment or use default
	port := os.Getenv("TEST_PORT")
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080" // Default API server port
		}
	}

	baseURL := fmt.Sprintf("http://localhost:%s", port)

	// Create test server instance
	// Note: The API server should be running separately
	// This test suite tests against a live API server
	server := &testServer{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Verify server is accessible
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", baseURL+"/swagger/index.html", nil)
	resp, err := server.client.Do(req)
	if err != nil {
		t.Fatalf("API server is not accessible at %s. Please start the server first: %v", baseURL, err)
	}
	resp.Body.Close()

	return server
}

// get performs a GET request and returns the response
func (ts *testServer) get(t *testing.T, path string) *http.Response {
	url := ts.baseURL + path
	resp, err := ts.client.Get(url)
	require.NoError(t, err, "Failed to GET %s", url)
	return resp
}

// getJSON performs a GET request and unmarshals JSON response
func (ts *testServer) getJSON(t *testing.T, path string, target interface{}) {
	resp := ts.get(t, path)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK for %s", path)
	contentType := resp.Header.Get("Content-Type")
	require.Contains(t, contentType, "application/json", "Content-Type should contain application/json, got: %s", contentType)

	err := json.NewDecoder(resp.Body).Decode(target)
	require.NoError(t, err, "Failed to decode JSON from %s", path)
}

// getJSONArray performs a GET request and unmarshals JSON array response
func (ts *testServer) getJSONArray(t *testing.T, path string) []map[string]interface{} {
	resp := ts.get(t, path)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK for %s", path)

	var result []map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode JSON array from %s", path)
	return result
}

// TestAPIHealth checks if the API server is accessible
func TestAPIHealth(t *testing.T) {
	ts := setupTestServer(t)

	// Test root redirect (may be 301 redirect or 200 if serving directly)
	resp := ts.get(t, "/")
	assert.Contains(t, []int{http.StatusMovedPermanently, http.StatusOK}, resp.StatusCode, "Root should redirect to Swagger or serve it directly")
	resp.Body.Close()

	// Test Swagger UI
	resp = ts.get(t, "/swagger/index.html")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Swagger UI should be accessible")
	resp.Body.Close()

	// Test health endpoint
	var health handlers.HealthResponse
	ts.getJSON(t, "/health", &health)
	assert.NotEmpty(t, health.Status)
	assert.NotEmpty(t, health.Timestamp)
	assert.NotEmpty(t, health.Links["self"])
	assert.NotEmpty(t, health.Links["metrics"])

	// Test metrics endpoint
	var metrics handlers.MetricsResponse
	ts.getJSON(t, "/metrics", &metrics)
	assert.Greater(t, metrics.ServerUptimeSeconds, 0.0)
	assert.NotEmpty(t, metrics.Links["self"])
	assert.NotEmpty(t, metrics.Links["health"])
}

// TestClustersEndpoints tests all cluster-related endpoints
func TestClustersEndpoints(t *testing.T) {
	ts := setupTestServer(t)

	// List clusters
	var clusters []handlers.ClusterResponse
	ts.getJSON(t, "/api/v1/clusters", &clusters)
	t.Logf("Found %d clusters", len(clusters))

	if len(clusters) > 0 {
		clusterID := clusters[0].ID
		t.Logf("Testing with cluster: %s", clusterID)

		// Get cluster details
		var cluster handlers.ClusterResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clusters/%s", clusterID), &cluster)
		assert.Equal(t, clusterID, cluster.ID)

		// Get cluster status
		var status handlers.ClusterStatusResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clusters/%s/status", clusterID), &status)
		assert.NotEmpty(t, status.Phase)

		// Get cluster metrics (may not be available)
		resp := ts.get(t, fmt.Sprintf("/api/v1/clusters/%s/metrics", clusterID))
		if resp.StatusCode == http.StatusOK {
			var metrics handlers.ClusterMetricsResponse
			json.NewDecoder(resp.Body).Decode(&metrics)
		}
		resp.Body.Close()

		// Get cluster bootstrap
		var bootstrap handlers.ClusterBootstrapResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clusters/%s/bootstrap", clusterID), &bootstrap)

		// Get cluster endpoints
		var endpoints handlers.ClusterEndpointResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clusters/%s/endpoints", clusterID), &endpoints)

		// Get Kubernetes status
		var k8sStatus handlers.KubernetesStatusResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clusters/%s/kubernetes-status", clusterID), &k8sStatus)

		// Get control plane status (may not be available)
		resp = ts.get(t, fmt.Sprintf("/api/v1/clusters/%s/controlplane-status", clusterID))
		if resp.StatusCode == http.StatusOK {
			var cpStatus handlers.ControlPlaneStatusResponse
			json.NewDecoder(resp.Body).Decode(&cpStatus)
		}
		resp.Body.Close()

		// Get diagnostics (may not be available)
		resp = ts.get(t, fmt.Sprintf("/api/v1/clusters/%s/diagnostics", clusterID))
		if resp.StatusCode == http.StatusOK {
			var diagnostics handlers.ClusterDiagnosticsResponse
			json.NewDecoder(resp.Body).Decode(&diagnostics)
		}
		resp.Body.Close()

		// Get destroy status (may not exist)
		resp = ts.get(t, fmt.Sprintf("/api/v1/clusters/%s/destroy-status", clusterID))
		resp.Body.Close()
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)

		// Get workload proxy status
		var wpStatus handlers.ClusterWorkloadProxyStatusResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clusters/%s/workload-proxy-status", clusterID), &wpStatus)
	}
}

// TestMachinesEndpoints tests all machine-related endpoints
func TestMachinesEndpoints(t *testing.T) {
	ts := setupTestServer(t)

	// List machines
	var machines []handlers.MachineResponse
	ts.getJSON(t, "/api/v1/machines", &machines)
	t.Logf("Found %d machines", len(machines))

	if len(machines) > 0 {
		machineID := machines[0].ID
		t.Logf("Testing with machine: %s", machineID)

		// Get machine details
		var machine handlers.MachineResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/machines/%s", machineID), &machine)
		assert.Equal(t, machineID, machine.ID)

		// Get machine labels (may not be available)
		resp := ts.get(t, fmt.Sprintf("/api/v1/machines/%s/labels", machineID))
		if resp.StatusCode == http.StatusOK {
			var labels handlers.MachineLabelsResponse
			json.NewDecoder(resp.Body).Decode(&labels)
		}
		resp.Body.Close()

		// Get machine extensions (may not be available or may error)
		resp = ts.get(t, fmt.Sprintf("/api/v1/machines/%s/extensions", machineID))
		if resp.StatusCode == http.StatusOK {
			var extensions handlers.MachineExtensionsResponse
			json.NewDecoder(resp.Body).Decode(&extensions)
		}
		resp.Body.Close()

		// Get machine upgrade status
		resp = ts.get(t, fmt.Sprintf("/api/v1/machines/%s/upgrade-status", machineID))
		resp.Body.Close()
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)

		// Get machine metrics (may not be available)
		resp = ts.get(t, fmt.Sprintf("/api/v1/machines/%s/metrics", machineID))
		if resp.StatusCode == http.StatusOK {
			var metrics handlers.MachineStatusMetricsResponse
			json.NewDecoder(resp.Body).Decode(&metrics)
		}
		resp.Body.Close()

		// Get machine config diff
		resp2 := ts.get(t, fmt.Sprintf("/api/v1/machines/%s/config-diff", machineID))
		resp2.Body.Close()
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp2.StatusCode)
	}
}

// TestMachineSetsEndpoints tests all machine set-related endpoints
func TestMachineSetsEndpoints(t *testing.T) {
	ts := setupTestServer(t)

	// List machine sets
	var machineSets []handlers.MachineSetResponse
	ts.getJSON(t, "/api/v1/machinesets", &machineSets)
	t.Logf("Found %d machine sets", len(machineSets))

	if len(machineSets) > 0 {
		msID := machineSets[0].ID
		t.Logf("Testing with machine set: %s", msID)

		// Get machine set details
		var ms handlers.MachineSetResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/machinesets/%s", msID), &ms)
		assert.Equal(t, msID, ms.ID)

		// Get machine set status
		var status handlers.MachineSetStatusResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/machinesets/%s/status", msID), &status)

		// Get destroy status
		resp := ts.get(t, fmt.Sprintf("/api/v1/machinesets/%s/destroy-status", msID))
		resp.Body.Close()
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)
	}
}

// TestClusterMachinesEndpoints tests all cluster machine-related endpoints
func TestClusterMachinesEndpoints(t *testing.T) {
	ts := setupTestServer(t)

	// List cluster machines
	var clusterMachines []handlers.ClusterMachineResponse
	ts.getJSON(t, "/api/v1/clustermachines", &clusterMachines)
	t.Logf("Found %d cluster machines", len(clusterMachines))

	if len(clusterMachines) > 0 {
		cmID := clusterMachines[0].ID
		t.Logf("Testing with cluster machine: %s", cmID)

		// Get cluster machine details
		var cm handlers.ClusterMachineResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clustermachines/%s", cmID), &cm)
		assert.Equal(t, cmID, cm.ID)

		// Get cluster machine status
		var status handlers.ClusterMachineStatusResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clustermachines/%s/status", cmID), &status)

		// Get config status
		var configStatus handlers.ClusterMachineConfigStatusResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clustermachines/%s/config-status", cmID), &configStatus)

		// Get Talos version
		var talosVersion handlers.ClusterMachineTalosVersionResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clustermachines/%s/talos-version", cmID), &talosVersion)

		// Get config (may not be available)
		resp := ts.get(t, fmt.Sprintf("/api/v1/clustermachines/%s/config", cmID))
		resp.Body.Close()
		// 404 is acceptable for config endpoint
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)
	}
}

// TestOtherEndpoints tests various other endpoints
func TestOtherEndpoints(t *testing.T) {
	ts := setupTestServer(t)

	// Test endpoints that should always return lists
	endpoints := []struct {
		path string
		name string
	}{
		{"/api/v1/configpatches", "config patches"},
		{"/api/v1/machineclasses", "machine classes"},
		{"/api/v1/machinesetnodes", "machine set nodes"},
		{"/api/v1/etcdbackups", "etcd backups"},
		{"/api/v1/schematics", "schematics"},
		{"/api/v1/ongoingtasks", "ongoing tasks"},
		{"/api/v1/kubernetes-versions", "kubernetes versions"},
		{"/api/v1/extensions-configurations", "extensions configurations"},
		{"/api/v1/kernel-args", "kernel args"},
		{"/api/v1/loadbalancer-configs", "load balancer configs"},
		{"/api/v1/exposed-services", "exposed services"},
		{"/api/v1/machine-request-sets", "machine request sets"},
		{"/api/v1/image-pull-requests", "image pull requests"},
		{"/api/v1/installation-medias", "installation medias"},
		{"/api/v1/infra-machine-configs", "infra machine configs"},
	}

	for _, ep := range endpoints {
		t.Run(ep.name, func(t *testing.T) {
			resp := ts.get(t, ep.path)
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK for %s", ep.path)
			resp.Body.Close()
		})
	}
}

// TestFiltering tests query parameter filtering
func TestFiltering(t *testing.T) {
	ts := setupTestServer(t)

	// Get a cluster ID first
	var clusters []handlers.ClusterResponse
	ts.getJSON(t, "/api/v1/clusters", &clusters)

	if len(clusters) > 0 {
		clusterID := clusters[0].ID

		// Test cluster filtering on cluster machines
		var clusterMachines []handlers.ClusterMachineResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/clustermachines?cluster=%s", clusterID), &clusterMachines)
		t.Logf("Found %d cluster machines for cluster %s", len(clusterMachines), clusterID)

		// Test cluster filtering on etcd backups
		var backups []handlers.EtcdBackupResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/etcdbackups?cluster=%s", clusterID), &backups)
		t.Logf("Found %d etcd backups for cluster %s", len(backups), clusterID)
	}

	// Test machine set filtering
	var machineSets []handlers.MachineSetResponse
	ts.getJSON(t, "/api/v1/machinesets", &machineSets)

	if len(machineSets) > 0 {
		msID := machineSets[0].ID
		var nodes []handlers.MachineSetNodeResponse
		ts.getJSON(t, fmt.Sprintf("/api/v1/machinesetnodes?machineset=%s", msID), &nodes)
		t.Logf("Found %d nodes for machine set %s", len(nodes), msID)
	}
}

// TestLinks tests that all responses include proper _links
func TestLinks(t *testing.T) {
	ts := setupTestServer(t)

	// Test cluster links
	var clusters []handlers.ClusterResponse
	ts.getJSON(t, "/api/v1/clusters", &clusters)

	if len(clusters) > 0 {
		cluster := clusters[0]
		assert.NotEmpty(t, cluster.Links, "Cluster should have links")
		assert.NotEmpty(t, cluster.Links["self"], "Cluster should have self link")
	}

	// Test machine links
	var machines []handlers.MachineResponse
	ts.getJSON(t, "/api/v1/machines", &machines)

	if len(machines) > 0 {
		machine := machines[0]
		assert.NotEmpty(t, machine.Links, "Machine should have links")
		assert.NotEmpty(t, machine.Links["self"], "Machine should have self link")
	}
}

// TestErrorHandling tests error responses
func TestErrorHandling(t *testing.T) {
	ts := setupTestServer(t)

	// Test error response for non-existent resource (404 or 500 are both acceptable error responses)
	resp := ts.get(t, "/api/v1/clusters/non-existent-cluster")
	assert.Contains(t, []int{http.StatusNotFound, http.StatusInternalServerError}, resp.StatusCode, "Should return error (404 or 500) for non-existent cluster")
	resp.Body.Close()

	resp = ts.get(t, "/api/v1/machines/non-existent-machine")
	assert.Contains(t, []int{http.StatusNotFound, http.StatusInternalServerError}, resp.StatusCode, "Should return error (404 or 500) for non-existent machine")
	resp.Body.Close()
}
