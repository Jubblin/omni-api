# gRPC API Audit - REST Implementation Gaps

## Executive Summary

This document provides a comprehensive audit of the gRPC API available through the Sidero Omni client (`github.com/siderolabs/omni/client v1.4.6`) and identifies gaps in the current REST implementation.

**Key Finding**: The current REST API implementation is **read-only** and only exposes a subset of the available gRPC services. All write operations and several service interfaces are not exposed via REST.

---

## Available gRPC Services

Based on the Omni client structure, the following services are available:

### 1. **Omni Service** (`client.Omni()`)

- **Status**: ✅ Partially Implemented
- **Current Usage**: Only `State()` interface is used (read-only resource access)
- **Available Methods**:
  - `State()` - Resource state management (read operations)
  - Potentially other methods for resource management

### 2. **Management Service** (`client.Management()`)

- **Status**: ❌ **NOT IMPLEMENTED**
- **Purpose**: Write operations (Create, Update, Delete resources)
- **Gap**: All resource modification operations are missing from REST API

### 3. **Talos Service** (`client.Talos()`)

- **Status**: ❌ **NOT IMPLEMENTED**
- **Purpose**: Direct Talos API operations
- **Gap**: Talos-specific operations not exposed

### 4. **Auth Service** (`client.Auth()`)

- **Status**: ❌ **NOT IMPLEMENTED**
- **Purpose**: Authentication operations
- **Gap**: Authentication management not exposed

### 5. **OIDC Service** (`client.OIDC()`)

- **Status**: ❌ **NOT IMPLEMENTED**
- **Purpose**: OIDC authentication operations
- **Gap**: OIDC operations not exposed

---

## Current REST Implementation Status

### Implemented: Read Operations (GET only)

The REST API currently implements **70+ GET endpoints** covering:

#### Resource Types Exposed (38+ types)

1. ✅ Clusters - List, Get, Status, Metrics, Bootstrap, Kubeconfig, Upgrades, Endpoints, Kubernetes Status, Nodes, Control Plane Status, Diagnostics, Destroy Status, Workload Proxy Status
2. ✅ Machines - List, Get, Status, Labels, Extensions, Upgrade Status, Metrics, Config Diff
3. ✅ MachineSets - List, Get, Status, Destroy Status
4. ✅ MachineSetNodes - List, Get
5. ✅ ConfigPatches - List, Get
6. ✅ ClusterMachines - List, Get, Status, Config Status, Talos Version, Config
7. ✅ MachineClass - List, Get
8. ✅ EtcdBackups - List, Get, Status
9. ✅ EtcdManualBackups - List, Get
10. ✅ Schematics - List, Get
11. ✅ SchematicConfigurations - List, Get
12. ✅ OngoingTasks - List, Get
13. ✅ KubernetesVersions - List, Get
14. ✅ ExtensionsConfigurations - List, Get
15. ✅ KernelArgs - List, Get
16. ✅ LoadBalancerConfigs - List, Get
17. ✅ LoadBalancerStatus - Get
18. ✅ ExposedServices - List, Get
19. ✅ MachineRequestSets - List, Get
20. ✅ ImagePullRequests - List, Get, Status
21. ✅ ImagePullStatus - Get
22. ✅ InstallationMedias - List, Get
23. ✅ InfraMachineConfigs - List, Get
24. ✅ MachineConfigDiff - Get
25. ✅ MachineStatus - Get
26. ✅ MachineStatusMetrics - Get
27. ✅ MachineLabels - Get
28. ✅ MachineExtensions - Get
29. ✅ MachineUpgradeStatus - Get
30. ✅ MachineSetStatus - Get
31. ✅ ClusterMachineStatus - Get
32. ✅ ClusterMachineConfigStatus - Get
33. ✅ ClusterMachineTalosVersion - Get
34. ✅ ClusterMachineConfig - Get
35. ✅ Kubeconfigs - Get
36. ✅ KubernetesUpgradeStatus - Get
37. ✅ TalosUpgradeStatus - Get
38. ✅ KubernetesStatus - Get
39. ✅ ClusterKubernetesNodes - List, Get
40. ✅ ControlPlaneStatus - Get
41. ✅ ClusterDiagnostics - Get
42. ✅ ClusterDestroyStatus - Get
43. ✅ ClusterWorkloadProxyStatus - Get
44. ✅ ClusterEndpoints - Get

---

## Identified Gaps

### 🔴 Critical Gaps: Write Operations

All resource modification operations are **completely missing** from the REST API:

#### 1. **Resource Creation (POST)**

- ❌ Create Clusters
- ❌ Create Machines
- ❌ Create MachineSets
- ❌ Create ConfigPatches
- ❌ Create ClusterMachines
- ❌ Create MachineClasses
- ❌ Create EtcdBackups (manual backup requests)
- ❌ Create Schematics
- ❌ Create SchematicConfigurations
- ❌ Create ExtensionsConfigurations
- ❌ Create KernelArgs
- ❌ Create LoadBalancerConfigs
- ❌ Create ExposedServices
- ❌ Create MachineRequestSets
- ❌ Create ImagePullRequests
- ❌ Create InstallationMedias
- ❌ Create InfraMachineConfigs
- ❌ Create ClusterSecrets (if needed)

#### 2. **Resource Updates (PUT/PATCH)**

- ❌ Update Clusters
- ❌ Update Machines (labels, extensions, maintenance mode)
- ❌ Update MachineSets
- ❌ Update ConfigPatches
- ❌ Update ClusterMachines
- ❌ Update MachineClasses
- ❌ Update Schematics
- ❌ Update SchematicConfigurations
- ❌ Update ExtensionsConfigurations
- ❌ Update KernelArgs
- ❌ Update LoadBalancerConfigs
- ❌ Update ExposedServices
- ❌ Update MachineRequestSets
- ❌ Update InstallationMedias
- ❌ Update InfraMachineConfigs

#### 3. **Resource Deletion (DELETE)**

- ❌ Delete Clusters
- ❌ Delete Machines
- ❌ Delete MachineSets
- ❌ Delete ConfigPatches
- ❌ Delete ClusterMachines
- ❌ Delete MachineClasses
- ❌ Delete EtcdBackups
- ❌ Delete Schematics
- ❌ Delete SchematicConfigurations
- ❌ Delete ExtensionsConfigurations
- ❌ Delete KernelArgs
- ❌ Delete LoadBalancerConfigs
- ❌ Delete ExposedServices
- ❌ Delete MachineRequestSets
- ❌ Delete ImagePullRequests
- ❌ Delete InstallationMedias
- ❌ Delete InfraMachineConfigs

#### 4. **Resource Actions (POST)**

- ❌ Trigger Kubernetes Upgrades
- ❌ Trigger Talos Upgrades
- ❌ Trigger Etcd Manual Backups
- ❌ Trigger Cluster Bootstrap
- ❌ Trigger Machine Reboots
- ❌ Trigger Machine Shutdowns
- ❌ Trigger Machine Resets
- ❌ Trigger Image Pulls
- ❌ Trigger Cluster Destruction
- ❌ Trigger MachineSet Destruction
- ❌ Trigger Machine Maintenance Mode Toggle

### 🟡 Medium Priority Gaps: Service-Specific Operations

#### 1. **Management Service Operations**

The `client.Management()` service likely provides:

- ❌ Resource lifecycle management
- ❌ Batch operations
- ❌ Resource validation
- ❌ Resource templating
- ❌ Resource import/export

#### 2. **Talos Service Operations**

The `client.Talos()` service likely provides:

- ❌ Direct Talos API calls
- ❌ Machine configuration management
- ❌ Machine command execution
- ❌ Machine file operations
- ❌ Machine service management

#### 3. **Auth Service Operations**

The `client.Auth()` service likely provides:

- ❌ Service account management
- ❌ API key management
- ❌ User authentication
- ❌ Permission management
- ❌ Token management

#### 4. **OIDC Service Operations**

The `client.OIDC()` service likely provides:

- ❌ OIDC configuration
- ❌ OIDC provider management
- ❌ OIDC authentication flows

### 🟢 Low Priority Gaps: Missing Resource Types

#### Security-Sensitive Resources

- ❌ **ClusterSecrets** - Secret management (marked as low priority in RESOURCES.md due to security concerns)

#### Internal Resources (Not typically exposed)

- ClusterConfigVersion
- ClusterUUID
- MachineStatusLink
- MachineStatusSnapshot
- BackupData
- ClusterMachineIdentity
- ClusterMachineEncryptionKey
- ClusterMachineTemplate

---

## Implementation Recommendations

### Priority 1: Write Operations (Critical)

Implement REST endpoints for resource modification:

1. **POST /api/v1/{resource}** - Create resources
2. **PUT /api/v1/{resource}/:id** - Update resources (full update)
3. **PATCH /api/v1/{resource}/:id** - Partial updates
4. **DELETE /api/v1/{resource}/:id** - Delete resources

**Resources to prioritize:**

- Clusters (create, update, delete)
- Machines (update labels, extensions, maintenance mode)
- MachineSets (create, update, delete)
- ConfigPatches (create, update, delete)
- EtcdBackups (trigger manual backups)

### Priority 2: Action Endpoints (High)

Implement action endpoints for triggering operations:

1. **POST /api/v1/clusters/:id/actions/kubernetes-upgrade** - Trigger Kubernetes upgrade
2. **POST /api/v1/clusters/:id/actions/talos-upgrade** - Trigger Talos upgrade
3. **POST /api/v1/clusters/:id/actions/bootstrap** - Trigger bootstrap
4. **POST /api/v1/clusters/:id/actions/destroy** - Trigger cluster destruction
5. **POST /api/v1/etcdbackups** - Trigger manual backup
6. **POST /api/v1/machines/:id/actions/reboot** - Reboot machine
7. **POST /api/v1/machines/:id/actions/shutdown** - Shutdown machine
8. **POST /api/v1/machines/:id/actions/reset** - Reset machine
9. **POST /api/v1/machines/:id/actions/maintenance** - Toggle maintenance mode

### Priority 3: Service-Specific Operations (Medium)

1. **Investigate Management Service API** - Document available methods
2. **Investigate Talos Service API** - Document available methods
3. **Implement Auth Service endpoints** (if needed for API management)
4. **Implement OIDC Service endpoints** (if needed)

### Priority 4: Additional Resources (Low)

1. **ClusterSecrets** - If security requirements allow

---

## Technical Implementation Notes

### Current Architecture

The REST API currently:

- Uses `client.Omni().State()` for all operations
- Implements read-only GET endpoints
- Uses COSI runtime State interface for resource queries
- Provides hypermedia links in responses

### Required Changes

To implement write operations:

1. **Add Management Service Client**

   ```go
   managementClient := client.Management()
   ```

2. **Implement Write Handlers**
   - Create handlers for POST, PUT, PATCH, DELETE
   - Add request validation
   - Add response formatting
   - Add error handling

3. **Add Action Handlers**
   - Create action-specific endpoints
   - Implement action validation
   - Add async operation support (if needed)

4. **Update Swagger Documentation**
   - Add write operation documentation
   - Add request/response schemas
   - Add error response schemas

5. **Add Tests**
   - Unit tests for write operations
   - Integration tests for write operations
   - Error handling tests

---

## Summary Statistics

| Category | Count | Status |
|----------|-------|--------|
| **Read Operations (GET)** | 70+ | ✅ Implemented |
| **Write Operations (POST/PUT/PATCH/DELETE)** | 0 | ❌ Missing |
| **Action Operations (POST)** | 0 | ❌ Missing |
| **Resource Types (Read)** | 38+ | ✅ Implemented |
| **Resource Types (Write)** | 38+ | ❌ Missing |
| **gRPC Services Used** | 1 of 5 | ⚠️ Partial |
| **gRPC Services Not Used** | 4 of 5 | ❌ Missing |

---

## Conclusion

The current REST API provides comprehensive **read-only** access to Omni resources but is missing all **write operations** and several **service interfaces**. To achieve feature parity with the gRPC API, the following must be implemented:

1. ✅ **Read Operations** - Complete (70+ endpoints)
2. ❌ **Write Operations** - Missing (Create, Update, Delete for all resources)
3. ❌ **Action Operations** - Missing (Upgrades, Backups, Machine actions)
4. ❌ **Service Interfaces** - Missing (Management, Talos, Auth, OIDC)

**Estimated Gap**: ~70% of gRPC API functionality is not exposed via REST (read operations represent ~30% of total API surface area).

---

## Next Steps

1. **Review this audit** with the team
2. **Prioritize gaps** based on business requirements
3. **Design write operation APIs** following RESTful principles
4. **Implement write operations** incrementally
5. **Add comprehensive tests** for write operations
6. **Update documentation** (Swagger/OpenAPI)

---

*Generated: 2025-01-27*
*Omni Client Version: v1.4.6*
*REST API Version: 0.0.10*
