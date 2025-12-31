# Omni Client Service API Reference

## Overview

This document provides a reference for the Omni client service APIs based on investigation and common patterns. This will be updated as actual API methods are discovered.

**Status**: 🔄 Investigation in progress - patterns documented, actual method signatures to be confirmed

---

## Management Service API

### Package
`github.com/siderolabs/omni/client/pkg/management`

### Client Type
```go
type Client struct {
    // Internal fields
}
```

### Expected Methods (Based on Patterns)

#### Resource Creation Methods

**CreateCluster**
```go
func (c *Client) CreateCluster(ctx context.Context, req *CreateClusterRequest) (*CreateClusterResponse, error)
```

**CreateMachineSet**
```go
func (c *Client) CreateMachineSet(ctx context.Context, req *CreateMachineSetRequest) (*CreateMachineSetResponse, error)
```

**CreateConfigPatch**
```go
func (c *Client) CreateConfigPatch(ctx context.Context, req *CreateConfigPatchRequest) (*CreateConfigPatchResponse, error)
```

#### Resource Update Methods

**UpdateCluster**
```go
func (c *Client) UpdateCluster(ctx context.Context, req *UpdateClusterRequest) (*UpdateClusterResponse, error)
```

**UpdateMachineSet**
```go
func (c *Client) UpdateMachineSet(ctx context.Context, req *UpdateMachineSetRequest) (*UpdateMachineSetResponse, error)
```

**UpdateConfigPatch**
```go
func (c *Client) UpdateConfigPatch(ctx context.Context, req *UpdateConfigPatchRequest) (*UpdateConfigPatchResponse, error)
```

**UpdateMachineLabels**
```go
func (c *Client) UpdateMachineLabels(ctx context.Context, machineID string, labels map[string]string) error
```

**UpdateMachineExtensions**
```go
func (c *Client) UpdateMachineExtensions(ctx context.Context, machineID string, extensions []string) error
```

**SetMachineMaintenance**
```go
func (c *Client) SetMachineMaintenance(ctx context.Context, machineID string, enabled bool) error
```

#### Resource Deletion Methods

**TeardownCluster** or **DeleteCluster**
```go
func (c *Client) TeardownCluster(ctx context.Context, clusterID string) error
// OR
func (c *Client) DeleteCluster(ctx context.Context, req *DeleteClusterRequest) error
```

**TeardownMachineSet** or **DeleteMachineSet**
```go
func (c *Client) TeardownMachineSet(ctx context.Context, machineSetID string) error
```

**DeleteConfigPatch**
```go
func (c *Client) DeleteConfigPatch(ctx context.Context, patchID string) error
```

#### Action Methods

**UpgradeKubernetes**
```go
func (c *Client) UpgradeKubernetes(ctx context.Context, clusterID string, version string) error
```

**UpgradeTalos**
```go
func (c *Client) UpgradeTalos(ctx context.Context, clusterID string, version string) error
```

**BootstrapCluster**
```go
func (c *Client) BootstrapCluster(ctx context.Context, clusterID string) error
```

**CreateEtcdManualBackup**
```go
func (c *Client) CreateEtcdManualBackup(ctx context.Context, clusterID string) (*EtcdBackupResponse, error)
```

### Request/Response Types (Expected)

```go
// Cluster Operations
type CreateClusterRequest struct {
    ID                string
    KubernetesVersion string
    TalosVersion      string
    Features          *ClusterFeatures
}

type UpdateClusterRequest struct {
    ID                string
    KubernetesVersion string
    TalosVersion      string
    Features          *ClusterFeatures
}

type ClusterFeatures struct {
    WorkloadProxy bool
    DiskEncryption bool
}

// MachineSet Operations
type CreateMachineSetRequest struct {
    ID             string
    Cluster        string
    MachineClass   string
    MachineCount   uint32
    UpdateStrategy string
    DeleteStrategy string
}

// ConfigPatch Operations
type CreateConfigPatchRequest struct {
    ID      string
    Cluster string
    Data    string
}
```

### Error Handling

Management service methods typically return gRPC status errors:
- `codes.NotFound` - Resource not found
- `codes.AlreadyExists` - Resource already exists
- `codes.InvalidArgument` - Invalid request parameters
- `codes.Internal` - Internal server error
- `codes.PermissionDenied` - Permission denied

---

## Talos Service API

### Package
`github.com/siderolabs/omni/client/pkg/talos`

### Client Type
```go
type Client struct {
    // Internal fields
}
```

### Expected Methods

#### Machine Control Methods

**RebootMachine**
```go
func (c *Client) RebootMachine(ctx context.Context, machineID string) error
```

**ShutdownMachine**
```go
func (c *Client) ShutdownMachine(ctx context.Context, machineID string) error
```

**ResetMachine**
```go
func (c *Client) ResetMachine(ctx context.Context, machineID string) error
```

### Error Handling

Talos service methods may return:
- `codes.NotFound` - Machine not found
- `codes.Unavailable` - Machine unreachable
- `codes.DeadlineExceeded` - Operation timeout
- `codes.Internal` - Internal error

---

## Auth Service API

### Package
`github.com/siderolabs/omni/client/pkg/auth`

### Client Type
```go
type Client struct {
    // Internal fields
}
```

### Expected Methods

**ListServiceAccounts**
```go
func (c *Client) ListServiceAccounts(ctx context.Context) ([]*ServiceAccount, error)
```

**GetServiceAccount**
```go
func (c *Client) GetServiceAccount(ctx context.Context, id string) (*ServiceAccount, error)
```

**CreateServiceAccount**
```go
func (c *Client) CreateServiceAccount(ctx context.Context, req *CreateServiceAccountRequest) (*ServiceAccount, error)
```

**DeleteServiceAccount**
```go
func (c *Client) DeleteServiceAccount(ctx context.Context, id string) error
```

---

## OIDC Service API

### Package
`github.com/siderolabs/omni/client/pkg/oidc`

### Client Type
```go
type Client struct {
    // Internal fields
}
```

### Expected Methods

**ListOIDCProviders**
```go
func (c *Client) ListOIDCProviders(ctx context.Context) ([]*OIDCProvider, error)
```

**GetOIDCProvider**
```go
func (c *Client) GetOIDCProvider(ctx context.Context, id string) (*OIDCProvider, error)
```

**CreateOIDCProvider**
```go
func (c *Client) CreateOIDCProvider(ctx context.Context, req *CreateOIDCProviderRequest) (*OIDCProvider, error)
```

**UpdateOIDCProvider**
```go
func (c *Client) UpdateOIDCProvider(ctx context.Context, req *UpdateOIDCProviderRequest) (*OIDCProvider, error)
```

**DeleteOIDCProvider**
```go
func (c *Client) DeleteOIDCProvider(ctx context.Context, id string) error
```

---

## Common Patterns

### Context Usage
All methods accept `context.Context` as the first parameter for:
- Request cancellation
- Timeout handling
- Request tracing

### Error Handling Pattern
```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

resp, err := mgmtClient.CreateCluster(ctx, req)
if err != nil {
    if status.Code(err) == codes.AlreadyExists {
        // Handle already exists
    } else if status.Code(err) == codes.InvalidArgument {
        // Handle invalid argument
    }
    return err
}
```

### Type Assertion Pattern
```go
// In handlers, we need to type assert the interface{} to actual client type
mgmtClient, ok := h.management.(*management.Client)
if !ok {
    return fmt.Errorf("invalid management client type")
}
```

---

## Integration Notes

### Current Handler Pattern
Handlers currently receive `interface{}` for service clients. We need to:
1. Type assert to actual client type
2. Call appropriate method
3. Handle errors
4. Map responses

### Example Integration (Cluster Create)
```go
func (h *ClusterWriteHandler) CreateCluster(c *gin.Context) {
    var req ClusterCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Type assert Management client
    mgmtClient, ok := h.management.(*management.Client)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid management client"})
        return
    }

    // Create request
    createReq := &management.CreateClusterRequest{
        ID:                req.ID,
        KubernetesVersion: req.KubernetesVersion,
        TalosVersion:      req.TalosVersion,
    }
    
    if req.Features.WorkloadProxy || req.Features.DiskEncryption {
        createReq.Features = &management.ClusterFeatures{
            WorkloadProxy: req.Features.WorkloadProxy,
            DiskEncryption: req.Features.DiskEncryption,
        }
    }

    // Call Management service
    resp, err := mgmtClient.CreateCluster(c.Request.Context(), createReq)
    if err != nil {
        // Map gRPC errors to HTTP
        if status.Code(err) == codes.AlreadyExists {
            c.JSON(http.StatusConflict, gin.H{"error": "cluster already exists"})
            return
        }
        if status.Code(err) == codes.InvalidArgument {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Map response
    clusterResp := ClusterResponse{
        ID:                resp.Cluster.ID,
        KubernetesVersion: resp.Cluster.KubernetesVersion,
        // ... map other fields
    }
    
    c.JSON(http.StatusCreated, clusterResp)
}
```

---

## Next Steps

1. **Verify Actual API**: Run investigation commands to get actual method signatures
2. **Update This Document**: Replace expected patterns with actual API
3. **Implement Integration**: Update handlers with actual API calls
4. **Test**: Verify all operations work correctly

---

*This document will be updated as actual API methods are discovered*
