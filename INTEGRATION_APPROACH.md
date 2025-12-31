# Integration Approach - Practical Implementation Guide

## Overview

Since the exact Management/Talos/Auth/OIDC service APIs need to be discovered, this document provides a practical approach to integration that will work once the APIs are known.

## Strategy: Type-Safe Wrapper Pattern

Instead of using `interface{}` everywhere, we'll create type-safe wrappers that can be easily updated once the actual API is known.

---

## Step 1: Create Service Interface Wrappers

### Management Service Wrapper

```go
// internal/client/management_wrapper.go
package client

import (
    "context"
    "github.com/siderolabs/omni/client/pkg/client"
)

// ManagementService defines the interface for Management operations
type ManagementService interface {
    // Cluster operations
    CreateCluster(ctx context.Context, id, k8sVersion, talosVersion string, features *ClusterFeatures) error
    UpdateCluster(ctx context.Context, id, k8sVersion, talosVersion string, features *ClusterFeatures) error
    DeleteCluster(ctx context.Context, id string) error
    
    // MachineSet operations
    CreateMachineSet(ctx context.Context, id, cluster, machineClass string, count uint32) error
    UpdateMachineSet(ctx context.Context, id string, updates *MachineSetUpdates) error
    DeleteMachineSet(ctx context.Context, id string) error
    
    // ConfigPatch operations
    CreateConfigPatch(ctx context.Context, id, cluster, data string) error
    UpdateConfigPatch(ctx context.Context, id, data string) error
    DeleteConfigPatch(ctx context.Context, id string) error
    
    // Machine operations
    UpdateMachineLabels(ctx context.Context, machineID string, labels map[string]string) error
    UpdateMachineExtensions(ctx context.Context, machineID string, extensions []string) error
    SetMachineMaintenance(ctx context.Context, machineID string, enabled bool) error
    
    // Actions
    UpgradeKubernetes(ctx context.Context, clusterID, version string) error
    UpgradeTalos(ctx context.Context, clusterID, version string) error
    BootstrapCluster(ctx context.Context, clusterID string) error
    CreateEtcdManualBackup(ctx context.Context, clusterID string) error
}

// managementService implements ManagementService
type managementService struct {
    client *client.Client
}

// NewManagementService creates a new ManagementService wrapper
func NewManagementService(c *client.Client) ManagementService {
    return &managementService{client: c.Management()}
}

// Implementation methods will call actual Management client methods
// These will be implemented once we know the actual API
```

### Talos Service Wrapper

```go
// internal/client/talos_wrapper.go
package client

import (
    "context"
    "github.com/siderolabs/omni/client/pkg/client"
)

// TalosService defines the interface for Talos operations
type TalosService interface {
    RebootMachine(ctx context.Context, machineID string) error
    ShutdownMachine(ctx context.Context, machineID string) error
    ResetMachine(ctx context.Context, machineID string) error
}

// talosService implements TalosService
type talosService struct {
    client *client.Client
}

func NewTalosService(c *client.Client) TalosService {
    return &talosService{client: c.Talos()}
}
```

---

## Step 2: Update Handlers to Use Wrappers

### Example: Cluster Write Handler

```go
// internal/api/handlers/clusterwrite.go
package handlers

import (
    "context"
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/jubblin/omni-api/internal/client"
    // ... other imports
)

type ClusterWriteHandler struct {
    state      state.State
    management client.ManagementService  // Use interface instead of interface{}
}

func NewClusterWriteHandler(s state.State, mgmt client.ManagementService) *ClusterWriteHandler {
    return &ClusterWriteHandler{
        state:      s,
        management: mgmt,
    }
}

func (h *ClusterWriteHandler) CreateCluster(c *gin.Context) {
    var req ClusterCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Call through interface - implementation handles actual API call
    err := h.management.CreateCluster(
        c.Request.Context(),
        req.ID,
        req.KubernetesVersion,
        req.TalosVersion,
        &client.ClusterFeatures{
            WorkloadProxy: req.Features.WorkloadProxy,
            DiskEncryption: req.Features.DiskEncryption,
        },
    )
    
    if err != nil {
        // Handle errors
        handleManagementError(c, err)
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id": req.ID,
        "message": "Cluster created successfully",
    })
}
```

---

## Step 3: Implement Wrappers Once API is Known

Once we discover the actual Management service API, we implement the wrapper methods:

```go
// internal/client/management_wrapper.go (implementation)

func (m *managementService) CreateCluster(ctx context.Context, id, k8sVersion, talosVersion string, features *ClusterFeatures) error {
    // Type assert to actual client
    mgmtClient := m.client.(*management.Client)
    
    // Create request based on actual API
    req := &management.CreateClusterRequest{
        ID: id,
        KubernetesVersion: k8sVersion,
        TalosVersion: talosVersion,
        // ... map features
    }
    
    // Call actual API
    _, err := mgmtClient.CreateCluster(ctx, req)
    return err
}
```

---

## Alternative Approach: Direct Integration

If the API is straightforward, we can integrate directly:

### Step 1: Update Handler Structs

```go
type ClusterWriteHandler struct {
    state      state.State
    management *management.Client  // Direct type instead of interface{}
}
```

### Step 2: Update Handler Creation

```go
// In main.go
import "github.com/siderolabs/omni/client/pkg/management"

clusterWriteHandler := handlers.NewClusterWriteHandler(
    client.Omni().State(),
    client.Management(),  // Direct client
)
```

### Step 3: Use Directly in Handlers

```go
func (h *ClusterWriteHandler) CreateCluster(c *gin.Context) {
    // ... validation ...
    
    req := &management.CreateClusterRequest{
        // ... map request ...
    }
    
    resp, err := h.management.CreateCluster(c.Request.Context(), req)
    // ... handle response ...
}
```

---

## Recommended Approach

**Use the Wrapper Pattern** because:
1. ✅ Type-safe interfaces
2. ✅ Easy to test (can mock interfaces)
3. ✅ Can adapt to API changes
4. ✅ Cleaner handler code
5. ✅ Better error handling abstraction

---

## Implementation Steps

1. **Create wrapper interfaces** (Step 1 above)
2. **Update handlers to use interfaces** (Step 2 above)
3. **Investigate actual API** (Phase 1)
4. **Implement wrapper methods** (Step 3 above)
5. **Test integration**
6. **Update documentation**

---

## Error Handling Helper

Create a helper function to map gRPC errors to HTTP:

```go
// internal/api/handlers/errors.go
package handlers

import (
    "net/http"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "github.com/gin-gonic/gin"
)

func handleManagementError(c *gin.Context, err error) {
    if err == nil {
        return
    }
    
    code := status.Code(err)
    switch code {
    case codes.NotFound:
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    case codes.AlreadyExists:
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    case codes.InvalidArgument:
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    case codes.PermissionDenied:
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
    case codes.Unauthenticated:
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    case codes.Unavailable:
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    }
}
```

---

## Next Steps

1. Create wrapper interfaces (if using wrapper pattern)
2. OR update handlers to use direct client types (if using direct approach)
3. Investigate actual API methods
4. Implement actual API calls
5. Add comprehensive error handling
6. Write tests

---

*This approach allows us to proceed with integration while the exact API is being discovered*
