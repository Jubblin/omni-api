# Phase 1 Investigation Guide

## Overview

This guide provides step-by-step instructions for investigating the Omni client service APIs to document available methods, signatures, and usage patterns.

**Estimated Time**: 8-14 hours  
**Status**: 🔄 In Progress

---

## Investigation Commands

### Prerequisites

1. Ensure Go modules are downloaded:
   ```bash
   go mod download
   ```

2. Verify Omni client version:
   ```bash
   go list -m github.com/siderolabs/omni/client
   ```

### Step 1: Management Service Investigation

#### 1.1 View Package Documentation
```bash
# Basic package info
go doc github.com/siderolabs/omni/client/pkg/management

# All exported symbols
go doc -all github.com/siderolabs/omni/client/pkg/management

# Specific type
go doc github.com/siderolabs/omni/client/pkg/management.Client
```

#### 1.2 List All Methods
```bash
# Get all methods on Client type
go doc github.com/siderolabs/omni/client/pkg/management.Client | grep -E "func|method"
```

#### 1.3 Inspect Source Code (if available)
```bash
# Find source location
go list -f '{{.Dir}}' github.com/siderolabs/omni/client/pkg/management

# View source files
ls -la $(go list -f '{{.Dir}}' github.com/siderolabs/omni/client/pkg/management)
```

#### 1.4 Document Findings

**What to document:**
- [ ] Client struct definition
- [ ] All public methods (Create*, Update*, Delete*, Teardown*)
- [ ] Method signatures (parameters, return types)
- [ ] Request types (CreateClusterRequest, UpdateClusterRequest, etc.)
- [ ] Response types
- [ ] Error types and patterns
- [ ] Context usage
- [ ] Example usage patterns

**Template:**
```go
// Management Service Client
type Client struct {
    // Fields (if visible)
}

// Methods
func (c *Client) CreateCluster(ctx context.Context, req *CreateClusterRequest) (*CreateClusterResponse, error)
func (c *Client) UpdateCluster(ctx context.Context, req *UpdateClusterRequest) (*UpdateClusterResponse, error)
func (c *Client) DeleteCluster(ctx context.Context, req *DeleteClusterRequest) error
// ... etc
```

---

### Step 2: Talos Service Investigation

#### 2.1 View Package Documentation
```bash
go doc github.com/siderolabs/omni/client/pkg/talos
go doc -all github.com/siderolabs/omni/client/pkg/talos
go doc github.com/siderolabs/omni/client/pkg/talos.Client
```

#### 2.2 Document Findings

**What to document:**
- [ ] Client struct definition
- [ ] Machine control methods (Reboot, Shutdown, Reset)
- [ ] Configuration methods
- [ ] Command execution methods
- [ ] File operation methods
- [ ] Request/Response types
- [ ] Error patterns
- [ ] Machine connection requirements

**Template:**
```go
// Talos Service Client
type Client struct {
    // Fields
}

// Methods
func (c *Client) Reboot(ctx context.Context, machineID string) error
func (c *Client) Shutdown(ctx context.Context, machineID string) error
func (c *Client) Reset(ctx context.Context, machineID string) error
// ... etc
```

---

### Step 3: Auth Service Investigation

#### 3.1 View Package Documentation
```bash
go doc github.com/siderolabs/omni/client/pkg/auth
go doc -all github.com/siderolabs/omni/client/pkg/auth
go doc github.com/siderolabs/omni/client/pkg/auth.Client
```

#### 3.2 Document Findings

**What to document:**
- [ ] Client struct definition
- [ ] Service account methods
- [ ] API key methods
- [ ] User authentication methods
- [ ] Permission management methods
- [ ] Request/Response types
- [ ] Security considerations
- [ ] Authentication patterns

---

### Step 4: OIDC Service Investigation

#### 4.1 View Package Documentation
```bash
go doc github.com/siderolabs/omni/client/pkg/oidc
go doc -all github.com/siderolabs/omni/client/pkg/oidc
go doc github.com/siderolabs/omni/client/pkg/oidc.Client
```

#### 4.2 Document Findings

**What to document:**
- [ ] Client struct definition
- [ ] OIDC provider methods
- [ ] Configuration methods
- [ ] Authentication flow methods
- [ ] Request/Response types
- [ ] OIDC-specific patterns

---

## Alternative Investigation Methods

### Method 1: Using Reflection (Runtime)

Create a test program that uses reflection to inspect services:

```go
package main

import (
    "context"
    "fmt"
    "reflect"
    
    "github.com/siderolabs/omni/client/pkg/client"
)

func main() {
    // Create client (requires credentials)
    c, _ := client.New("grpc://localhost:443", client.WithServiceAccount("..."))
    defer c.Close()
    
    // Inspect Management service
    mgmt := c.Management()
    inspectService(mgmt, "Management")
    
    // Inspect Talos service
    talos := c.Talos()
    inspectService(talos, "Talos")
    
    // Inspect Auth service
    auth := c.Auth()
    inspectService(auth, "Auth")
    
    // Inspect OIDC service
    oidc := c.OIDC()
    inspectService(oidc, "OIDC")
}

func inspectService(service interface{}, name string) {
    fmt.Printf("\n=== %s Service ===\n", name)
    t := reflect.TypeOf(service)
    
    for i := 0; i < t.NumMethod(); i++ {
        m := t.Method(i)
        fmt.Printf("  %s%s\n", m.Name, formatSignature(m.Type))
    }
}

func formatSignature(ft reflect.Type) string {
    // Format method signature
    // Implementation here
    return ""
}
```

### Method 2: Source Code Analysis

If you have access to the Omni client source code:

1. Clone the repository:
   ```bash
   git clone https://github.com/siderolabs/omni.git
   cd omni/client/pkg
   ```

2. Examine each service package:
   ```bash
   find management -name "*.go" -exec grep -l "func.*Client" {} \;
   find talos -name "*.go" -exec grep -l "func.*Client" {} \;
   find auth -name "*.go" -exec grep -l "func.*Client" {} \;
   find oidc -name "*.go" -exec grep -l "func.*Client" {} \;
   ```

3. Extract method signatures:
   ```bash
   grep -E "func.*\(.*Client\)" management/*.go
   ```

### Method 3: Using Go Tools

```bash
# List all types in package
go list -f '{{range .Types}}{{.Name}}{{end}}' github.com/siderolabs/omni/client/pkg/management

# Get detailed type information
go doc -all github.com/siderolabs/omni/client/pkg/management | grep -A 10 "type Client"
```

---

## Documentation Template

For each service, create a section in `PHASE1_INVESTIGATION.md`:

```markdown
## Management Service API

### Client Type
```go
type Client struct {
    // Fields
}
```

### Methods

#### CreateCluster
```go
func (c *Client) CreateCluster(ctx context.Context, req *CreateClusterRequest) (*CreateClusterResponse, error)
```

**Request Type:**
```go
type CreateClusterRequest struct {
    ID                string
    KubernetesVersion string
    TalosVersion      string
    Features          *ClusterFeatures
    // ... other fields
}
```

**Response Type:**
```go
type CreateClusterResponse struct {
    Cluster *Cluster
    // ... other fields
}
```

**Error Handling:**
- `codes.AlreadyExists` - Cluster already exists
- `codes.InvalidArgument` - Invalid request parameters
- `codes.Internal` - Internal server error

**Usage Example:**
```go
req := &management.CreateClusterRequest{
    ID: "my-cluster",
    KubernetesVersion: "v1.28.0",
}
resp, err := mgmtClient.CreateCluster(ctx, req)
```

#### UpdateCluster
[Similar structure...]

#### DeleteCluster
[Similar structure...]

### Common Patterns
- All methods take `context.Context` as first parameter
- Errors are gRPC status codes
- Request/Response types are protobuf-generated
```

---

## Investigation Checklist

### For Each Service:

- [ ] Package location identified
- [ ] Client type documented
- [ ] All public methods listed
- [ ] Method signatures documented
- [ ] Request types documented
- [ ] Response types documented
- [ ] Error patterns documented
- [ ] Context usage documented
- [ ] Example usage created
- [ ] Common patterns identified

### Overall:

- [ ] All services investigated
- [ ] API reference document created
- [ ] Integration patterns documented
- [ ] Error mapping guide created
- [ ] Examples provided for each service

---

## Next Steps After Investigation

1. **Create API Reference Document**
   - Comprehensive method reference
   - Request/Response type definitions
   - Error handling guide

2. **Create Integration Examples**
   - Code examples for each operation
   - Error handling examples
   - Best practices

3. **Update Handler Implementations**
   - Replace placeholders with actual API calls
   - Implement error handling
   - Add validation

---

## Running the Investigation

### Option 1: Manual Investigation
Run the commands above in your terminal and document findings in `PHASE1_INVESTIGATION.md`.

### Option 2: Automated Script
Use the provided scripts:
```bash
# Run investigation script
./scripts/inspect_omni_client.sh

# Or use Go investigation tool
go run scripts/investigate_services.go
```

### Option 3: Source Code Analysis
If you have access to source code, examine the package files directly.

---

## Expected Deliverables

1. **PHASE1_INVESTIGATION.md** - Complete investigation results
2. **API_REFERENCE.md** - Comprehensive API reference
3. **INTEGRATION_EXAMPLES.md** - Code examples for integration
4. **ERROR_HANDLING_GUIDE.md** - Error mapping and handling guide

---

*Start with Management service (highest priority), then proceed to others.*
