# Phase 1: Investigation Tasks

## Task List

### Task 1: Management Service Investigation ⏳
**Status**: In Progress  
**Estimated Time**: 2-3 hours  
**Priority**: High

**Steps**:
1. Run investigation commands:
   ```bash
   go doc github.com/siderolabs/omni/client/pkg/management
   go doc -all github.com/siderolabs/omni/client/pkg/management
   go doc github.com/siderolabs/omni/client/pkg/management.Client
   ```

2. Document findings in `PHASE1_INVESTIGATION.md`:
   - [ ] Client struct definition
   - [ ] All public methods
   - [ ] Method signatures
   - [ ] Request/Response types
   - [ ] Error patterns

3. Identify methods needed for our handlers:
   - [ ] CreateCluster
   - [ ] UpdateCluster
   - [ ] DeleteCluster/TeardownCluster
   - [ ] CreateMachineSet
   - [ ] UpdateMachineSet
   - [ ] DeleteMachineSet
   - [ ] CreateConfigPatch
   - [ ] UpdateConfigPatch
   - [ ] DeleteConfigPatch
   - [ ] UpdateMachineLabels
   - [ ] UpdateMachineExtensions
   - [ ] SetMachineMaintenance
   - [ ] UpgradeKubernetes
   - [ ] UpgradeTalos
   - [ ] BootstrapCluster
   - [ ] TeardownCluster
   - [ ] CreateEtcdManualBackup

**Deliverable**: Complete Management service API documentation

---

### Task 2: Talos Service Investigation ⏳
**Status**: Pending  
**Estimated Time**: 1-2 hours  
**Priority**: High

**Steps**:
1. Run investigation commands:
   ```bash
   go doc github.com/siderolabs/omni/client/pkg/talos
   go doc -all github.com/siderolabs/omni/client/pkg/talos
   go doc github.com/siderolabs/omni/client/pkg/talos.Client
   ```

2. Document findings in `PHASE1_INVESTIGATION.md`:
   - [ ] Client struct definition
   - [ ] Machine control methods
   - [ ] Method signatures
   - [ ] Request/Response types
   - [ ] Error patterns
   - [ ] Machine connection requirements

3. Identify methods needed for our handlers:
   - [ ] RebootMachine
   - [ ] ShutdownMachine
   - [ ] ResetMachine

**Deliverable**: Complete Talos service API documentation

---

### Task 3: Auth Service Investigation ⏳
**Status**: Pending  
**Estimated Time**: 1-2 hours  
**Priority**: Medium (if needed)

**Steps**:
1. Run investigation commands:
   ```bash
   go doc github.com/siderolabs/omni/client/pkg/auth
   go doc -all github.com/siderolabs/omni/client/pkg/auth
   go doc github.com/siderolabs/omni/client/pkg/auth.Client
   ```

2. Document findings in `PHASE1_INVESTIGATION.md`:
   - [ ] Client struct definition
   - [ ] Service account methods
   - [ ] Method signatures
   - [ ] Request/Response types
   - [ ] Security considerations

3. Identify methods needed for our handlers:
   - [ ] ListServiceAccounts
   - [ ] GetServiceAccount
   - [ ] CreateServiceAccount
   - [ ] DeleteServiceAccount

**Deliverable**: Complete Auth service API documentation

---

### Task 4: OIDC Service Investigation ⏳
**Status**: Pending  
**Estimated Time**: 1-2 hours  
**Priority**: Low (if needed)

**Steps**:
1. Run investigation commands:
   ```bash
   go doc github.com/siderolabs/omni/client/pkg/oidc
   go doc -all github.com/siderolabs/omni/client/pkg/oidc
   go doc github.com/siderolabs/omni/client/pkg/oidc.Client
   ```

2. Document findings in `PHASE1_INVESTIGATION.md`:
   - [ ] Client struct definition
   - [ ] OIDC provider methods
   - [ ] Method signatures
   - [ ] Request/Response types
   - [ ] OIDC-specific patterns

3. Identify methods needed for our handlers:
   - [ ] ListOIDCProviders
   - [ ] GetOIDCProvider
   - [ ] CreateOIDCProvider
   - [ ] UpdateOIDCProvider
   - [ ] DeleteOIDCProvider

**Deliverable**: Complete OIDC service API documentation

---

### Task 5: Create API Reference Document ⏳
**Status**: Pending  
**Estimated Time**: 2-3 hours  
**Priority**: High

**Steps**:
1. Consolidate all findings from Tasks 1-4
2. Create comprehensive API reference document
3. Include:
   - [ ] All method signatures
   - [ ] Request/Response type definitions
   - [ ] Error handling guide
   - [ ] Usage examples
   - [ ] Common patterns

**Deliverable**: `API_REFERENCE.md`

---

### Task 6: Create Integration Examples ⏳
**Status**: Pending  
**Estimated Time**: 2-3 hours  
**Priority**: High

**Steps**:
1. Create code examples for each operation type
2. Include error handling examples
3. Document best practices
4. Create integration patterns guide

**Deliverable**: `INTEGRATION_EXAMPLES.md`

---

## Quick Start Commands

### Run Automated Investigation
```bash
./scripts/run_investigation.sh
```

### Manual Investigation
```bash
# Management Service
go doc -all github.com/siderolabs/omni/client/pkg/management > investigation_results/management_all.txt

# Talos Service
go doc -all github.com/siderolabs/omni/client/pkg/talos > investigation_results/talos_all.txt

# Auth Service
go doc -all github.com/siderolabs/omni/client/pkg/auth > investigation_results/auth_all.txt

# OIDC Service
go doc -all github.com/siderolabs/omni/client/pkg/oidc > investigation_results/oidc_all.txt
```

---

## Progress Tracking

- [ ] Task 1: Management Service (In Progress)
- [ ] Task 2: Talos Service
- [ ] Task 3: Auth Service
- [ ] Task 4: OIDC Service
- [ ] Task 5: API Reference Document
- [ ] Task 6: Integration Examples

**Overall Progress**: 0/6 tasks complete

---

*Update this file as tasks are completed*
