# Implementation Status - Remaining Services

## Overview

This document tracks the implementation status of the remaining gRPC services identified in the audit.

## ✅ Completed

### 1. Service Wrapper Created

- **File**: `internal/client/services.go`
- **Status**: Complete
- **Purpose**: Provides access to all Omni client services (State, Management, Talos, Auth, OIDC)

### 2. Write Operation Handlers Created

#### Cluster Write Operations

- **File**: `internal/api/handlers/clusterwrite.go`
- **Endpoints**:
  - `POST /api/v1/clusters` - Create cluster
  - `PUT /api/v1/clusters/:id` - Update cluster
  - `DELETE /api/v1/clusters/:id` - Delete cluster
- **Status**: Structure complete, needs Management service API integration

#### Machine Write Operations

- **File**: `internal/api/handlers/machinewrite.go`
- **Endpoints**:
  - `PATCH /api/v1/machines/:id` - Update machine (labels, extensions, maintenance)
- **Status**: Structure complete, needs Management service API integration

#### ConfigPatch Write Operations

- **File**: `internal/api/handlers/configpatchwrite.go`
- **Endpoints**:
  - `POST /api/v1/configpatches` - Create config patch
  - `PUT /api/v1/configpatches/:id` - Update config patch
  - `DELETE /api/v1/configpatches/:id` - Delete config patch
- **Status**: Structure complete, needs Management service API integration

#### MachineSet Write Operations

- **File**: `internal/api/handlers/machinesetwrite.go`
- **Endpoints**:
  - `POST /api/v1/machinesets` - Create machine set
  - `PUT /api/v1/machinesets/:id` - Update machine set
  - `DELETE /api/v1/machinesets/:id` - Delete machine set
- **Status**: Structure complete, needs Management service API integration

### 3. Action Handlers Created

#### Cluster Actions

- **File**: `internal/api/handlers/clusteractions.go`
- **Endpoints**:
  - `POST /api/v1/clusters/:id/actions/kubernetes-upgrade` - Trigger Kubernetes upgrade
  - `POST /api/v1/clusters/:id/actions/talos-upgrade` - Trigger Talos upgrade
  - `POST /api/v1/clusters/:id/actions/bootstrap` - Trigger bootstrap
  - `POST /api/v1/clusters/:id/actions/destroy` - Trigger cluster destruction
- **Status**: Structure complete, needs Management/Talos service API integration

#### Machine Actions

- **File**: `internal/api/handlers/machineactions.go`
- **Endpoints**:
  - `POST /api/v1/machines/:id/actions/reboot` - Reboot machine
  - `POST /api/v1/machines/:id/actions/shutdown` - Shutdown machine
  - `POST /api/v1/machines/:id/actions/reset` - Reset machine
  - `POST /api/v1/machines/:id/actions/maintenance` - Toggle maintenance mode
- **Status**: Structure complete, needs Talos/Management service API integration

#### EtcdBackup Actions

- **File**: `internal/api/handlers/etcdbackupactions.go`
- **Endpoints**:
  - `POST /api/v1/etcdbackups` - Trigger manual etcd backup
- **Status**: Structure complete, needs Management service API integration

#### MachineSet Actions

- **File**: `internal/api/handlers/machinesetactions.go`
- **Endpoints**:
  - `POST /api/v1/machinesets/:id/actions/destroy` - Trigger machine set destruction
- **Status**: Structure complete, needs Management service API integration

### 4. Auth and OIDC Service Handlers Created

#### Auth Service Endpoints

- **File**: `internal/api/handlers/auth.go`
- **Endpoints**:
  - `GET /api/v1/auth/service-accounts` - List service accounts
  - `GET /api/v1/auth/service-accounts/:id` - Get service account
  - `POST /api/v1/auth/service-accounts` - Create service account
  - `DELETE /api/v1/auth/service-accounts/:id` - Delete service account
- **Status**: Structure complete, needs Auth service API integration

#### OIDC Service Endpoints

- **File**: `internal/api/handlers/oidc.go`
- **Endpoints**:
  - `GET /api/v1/oidc/providers` - List OIDC providers
  - `GET /api/v1/oidc/providers/:id` - Get OIDC provider
  - `POST /api/v1/oidc/providers` - Create OIDC provider
  - `PUT /api/v1/oidc/providers/:id` - Update OIDC provider
  - `DELETE /api/v1/oidc/providers/:id` - Delete OIDC provider
- **Status**: Structure complete, needs OIDC service API integration

### 5. Routes Registered

- **File**: `main.go`
- **Status**: All new endpoints registered in the router
- **Note**: Endpoints are ready but return placeholder responses

---

## ⚠️ Pending: Service API Integration

### Management Service Integration

All write operation handlers currently have placeholder implementations. They need to be integrated with the actual Management service API.

**Files needing integration:**

- `internal/api/handlers/clusterwrite.go`
- `internal/api/handlers/machinewrite.go`
- `internal/api/handlers/configpatchwrite.go`
- `internal/api/handlers/clusteractions.go`
- `internal/api/handlers/etcdbackupactions.go`

**Required steps:**

1. Investigate Management service API methods
2. Replace placeholder calls with actual API calls
3. Add proper error handling
4. Add request/response validation

**Example pattern to follow:**

```go
// Current placeholder
log.Printf("Creating cluster %s (Management service integration needed)", req.ID)

// Should become (example - actual API may differ)
mgmtClient := h.management.(*management.Client)
err := mgmtClient.CreateCluster(ctx, &management.CreateClusterRequest{
    ID: req.ID,
    KubernetesVersion: req.KubernetesVersion,
    // ...
})
```

### Talos Service Integration

Machine action handlers need Talos service integration for direct machine operations.

**Files needing integration:**

- `internal/api/handlers/machineactions.go` (reboot, shutdown, reset)
- `internal/api/handlers/clusteractions.go` (some operations may use Talos)

**Required steps:**

1. Investigate Talos service API methods
2. Replace placeholder calls with actual API calls
3. Add proper error handling

---

## 📋 Not Yet Implemented

### Additional Write Operations

The following resources still need write operation handlers:

- ClusterMachines (create, update, delete)
- MachineClasses (create, update, delete)
- Schematics (create, update, delete)
- SchematicConfigurations (create, update, delete)
- ExtensionsConfigurations (create, update, delete)
- KernelArgs (create, update, delete)
- LoadBalancerConfigs (create, update, delete)
- ExposedServices (create, update, delete)
- MachineRequestSets (create, update, delete)
- ImagePullRequests (create, update, delete)
- InstallationMedias (create, update, delete)
- InfraMachineConfigs (create, update, delete)

**Status**: Not started
**Priority**: Medium (can be added incrementally)

---

## 🔧 Next Steps

### Immediate (High Priority)

1. **Investigate Management Service API**
   - Review Omni client documentation
   - Identify method signatures for Create, Update, Delete, Teardown operations
   - Document API patterns

2. **Integrate Management Service**
   - Update cluster write handlers
   - Update machine write handlers
   - Update configpatch write handlers
   - Update action handlers

3. **Investigate Talos Service API**
   - Review Talos client documentation
   - Identify method signatures for machine operations
   - Document API patterns

4. **Integrate Talos Service**
   - Update machine action handlers
   - Test machine operations

### Short Term (Medium Priority)

1. **Add Tests**
   - Unit tests for write operation handlers
   - Integration tests for write operations
   - Mock Management/Talos services for testing

2. **Add Additional Write Handlers**
   - MachineSets write operations
   - Other high-priority resources

### Long Term (Low Priority)

1. **Auth Service Endpoints** (if needed)
2. **OIDC Service Endpoints** (if needed)
3. **Complete remaining resource write operations**

---

## 📝 Notes

### Current Implementation Pattern

All write operation handlers follow this pattern:

1. **Request Validation**: Validate incoming JSON request
2. **Resource Verification**: Check if resource exists (for updates/deletes)
3. **Service Call**: Call Management/Talos service (currently placeholder)
4. **Response**: Return appropriate HTTP status and response body

### Placeholder Responses

Currently, all write operations return responses with a `note` field indicating that Management service integration is required. This allows:

- API structure to be complete
- Swagger documentation to be generated
- Endpoints to be tested for routing
- Actual integration to be added incrementally

### Error Handling

Error handling structure is in place:

- 400 Bad Request for invalid input
- 404 Not Found for missing resources
- 500 Internal Server Error for service failures

Actual error responses from Management/Talos services need to be mapped appropriately.

---

## 🎯 Summary

**Completed:**

- ✅ Service wrapper for all services
- ✅ Write operation handler structure (4 resources: clusters, machines, machinesets, configpatches)
- ✅ Action handler structure (clusters, machines, machinesets, etcd backups)
- ✅ Auth service endpoints structure
- ✅ OIDC service endpoints structure
- ✅ All routes registered in main.go

**Pending:**

- ⚠️ Management service API integration
- ⚠️ Talos service API integration
- ⚠️ Auth service API integration
- ⚠️ OIDC service API integration
- ⚠️ Tests for new endpoints

**Not Started:**

- ❌ Additional resource write operations (ClusterMachines, MachineClasses, etc.)

**Estimated Completion:**

- Core write operations: ~85% (structure done for 4 resources, needs API integration)
- Action operations: ~85% (structure done, needs API integration)
- Auth/OIDC services: ~80% (structure done, needs API integration)
- Overall remaining services: ~60% (structure in place, needs implementation)

---

*Last Updated: 2025-01-27*
