# Integration Quick Reference

## Code Pattern Examples

### Current Placeholder Pattern

```go
// Current placeholder in clusterwrite.go
log.Printf("Creating cluster %s (Management service integration needed)", req.ID)
c.JSON(http.StatusCreated, gin.H{
    "message": "Cluster creation initiated",
    "id": req.ID,
    "note": "Management service integration required for actual creation",
})
```

### Expected Integration Pattern

```go
// Expected pattern after integration
mgmtClient := h.management.(*management.Client)
ctx := c.Request.Context()

// Example: Create cluster
createReq := &management.CreateClusterRequest{
    ID: req.ID,
    KubernetesVersion: req.KubernetesVersion,
    TalosVersion: req.TalosVersion,
    Features: &management.ClusterFeatures{
        WorkloadProxy: req.Features.WorkloadProxy,
        DiskEncryption: req.Features.DiskEncryption,
    },
}

resp, err := mgmtClient.CreateCluster(ctx, createReq)
if err != nil {
    // Map gRPC errors to HTTP errors
    if status.Code(err) == codes.AlreadyExists {
        c.JSON(http.StatusConflict, gin.H{"error": "cluster already exists"})
        return
    }
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}

// Map response to REST model
clusterResp := ClusterResponse{
    ID: resp.ID,
    KubernetesVersion: resp.KubernetesVersion,
    // ...
}
c.JSON(http.StatusCreated, clusterResp)
```

## Integration Checklist Per Handler

### For Each Write Handler:

- [ ] Import Management service types
- [ ] Type assert `h.management` to actual client type
- [ ] Map request model to service request type
- [ ] Call service method with context
- [ ] Handle gRPC errors (map to HTTP status codes)
- [ ] Map service response to REST response model
- [ ] Remove placeholder log/response
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Update Swagger docs if needed

## Error Mapping Guide

| gRPC Status Code | HTTP Status Code | Example |
|-----------------|------------------|---------|
| `codes.OK` | 200/201/202 | Success |
| `codes.InvalidArgument` | 400 | Bad Request |
| `codes.NotFound` | 404 | Resource not found |
| `codes.AlreadyExists` | 409 | Conflict |
| `codes.PermissionDenied` | 403 | Forbidden |
| `codes.Unauthenticated` | 401 | Unauthorized |
| `codes.Internal` | 500 | Internal Server Error |
| `codes.Unavailable` | 503 | Service Unavailable |
| `codes.DeadlineExceeded` | 504 | Gateway Timeout |

## Files Requiring Integration

### Management Service (44-64 hours)
1. `internal/api/handlers/clusterwrite.go` - 3 methods
2. `internal/api/handlers/machinewrite.go` - 1 method
3. `internal/api/handlers/machinesetwrite.go` - 3 methods
4. `internal/api/handlers/configpatchwrite.go` - 3 methods
5. `internal/api/handlers/clusteractions.go` - 4 methods
6. `internal/api/handlers/machinesetactions.go` - 1 method
7. `internal/api/handlers/etcdbackupactions.go` - 1 method

**Total**: ~16 methods to integrate

### Talos Service (19-27 hours)
1. `internal/api/handlers/machineactions.go` - 4 methods

**Total**: ~4 methods to integrate

### Auth Service (20-27 hours)
1. `internal/api/handlers/auth.go` - 4 methods

**Total**: ~4 methods to integrate

### OIDC Service (23-31 hours)
1. `internal/api/handlers/oidc.go` - 5 methods

**Total**: ~5 methods to integrate

## Investigation Steps

### Step 1: Find Service Client Types

```bash
# In Go module cache or source
find ~/go/pkg/mod/github.com/siderolabs/omni/client@v1.4.6 -name "*.go" | \
  grep -E "(management|talos|auth|oidc)" | \
  head -20
```

### Step 2: Identify Method Signatures

Look for patterns like:
- `Create*` methods
- `Update*` methods
- `Delete*` or `Teardown*` methods
- `*Request` and `*Response` types

### Step 3: Check Examples

Look for:
- Test files in the Omni client package
- Example code in documentation
- Existing usage in other projects

## Testing Strategy

### Unit Tests
- Mock the service clients
- Test request/response mapping
- Test error handling

### Integration Tests
- Use test Omni instance
- Test with real resources
- Test error scenarios

### Example Test Structure

```go
func TestCreateCluster(t *testing.T) {
    // Setup mock Management client
    mockMgmt := &MockManagementClient{}
    handler := NewClusterWriteHandler(mockState, mockMgmt)
    
    // Test request
    req := ClusterCreateRequest{
        ID: "test-cluster",
        KubernetesVersion: "v1.28.0",
    }
    
    // Verify service called correctly
    // Verify response mapping
    // Verify error handling
}
```

## Priority Order

1. **Management Service - Clusters** (highest priority)
   - Most commonly used resource
   - Foundation for other operations
   - ~8-12 hours

2. **Management Service - MachineSets** (high priority)
   - Core infrastructure resource
   - ~6-8 hours

3. **Management Service - ConfigPatches** (high priority)
   - Common configuration task
   - ~4-6 hours

4. **Talos Service - Machine Actions** (medium-high priority)
   - Important for machine management
   - ~9-12 hours

5. **Management Service - Machine Updates** (medium priority)
   - Less frequent but useful
   - ~3-4 hours

6. **Auth Service** (medium priority, if needed)
   - Only if service account management needed
   - ~20-27 hours

7. **OIDC Service** (low priority, if needed)
   - Only if OIDC management needed
   - ~23-31 hours

---

*See INTEGRATION_ANALYSIS.md for detailed breakdown*
