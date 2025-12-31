# Phase 1: Service API Investigation

## Overview

This document tracks the investigation of Omni client service APIs to identify available methods, their signatures, and usage patterns.

**Goal**: Document all available methods for Management, Talos, Auth, and OIDC services to enable integration.

**Status**: 🔄 In Progress

**See**: `PHASE1_INVESTIGATION_GUIDE.md` for detailed investigation instructions and commands.

---

## Investigation Methods

1. **Go Documentation**: Using `go doc` to inspect package APIs
2. **Source Code Analysis**: Examining package structure
3. **Type Inspection**: Using reflection to understand types
4. **Example Code**: Looking for usage patterns

---

## Management Service API

### Package Location
`github.com/siderolabs/omni/client/pkg/management`

### Investigation Steps
- [ ] List all exported types and functions
- [ ] Identify Client struct and methods
- [ ] Document Create/Update/Delete methods
- [ ] Document action methods (upgrades, backups, etc.)
- [ ] Document request/response types
- [ ] Document error handling patterns

### Findings

#### Client Structure
```go
// To be filled in
```

#### Available Methods
```go
// To be filled in
```

#### Request/Response Types
```go
// To be filled in
```

#### Error Patterns
```go
// To be filled in
```

### Notes
_Investigation in progress..._

---

## Talos Service API

### Package Location
`github.com/siderolabs/omni/client/pkg/talos`

### Investigation Steps
- [ ] List all exported types and functions
- [ ] Identify Client struct and methods
- [ ] Document machine control methods (reboot, shutdown, reset)
- [ ] Document configuration methods
- [ ] Document request/response types
- [ ] Document error handling patterns

### Findings

#### Client Structure
```go
// To be filled in
```

#### Available Methods
```go
// To be filled in
```

#### Request/Response Types
```go
// To be filled in
```

#### Error Patterns
```go
// To be filled in
```

### Notes
_Investigation in progress..._

---

## Auth Service API

### Package Location
`github.com/siderolabs/omni/client/pkg/auth`

### Investigation Steps
- [ ] List all exported types and functions
- [ ] Identify Client struct and methods
- [ ] Document service account methods
- [ ] Document API key methods
- [ ] Document request/response types
- [ ] Document error handling patterns

### Findings

#### Client Structure
```go
// To be filled in
```

#### Available Methods
```go
// To be filled in
```

#### Request/Response Types
```go
// To be filled in
```

#### Error Patterns
```go
// To be filled in
```

### Notes
_Investigation in progress..._

---

## OIDC Service API

### Package Location
`github.com/siderolabs/omni/client/pkg/oidc`

### Investigation Steps
- [ ] List all exported types and functions
- [ ] Identify Client struct and methods
- [ ] Document OIDC provider methods
- [ ] Document configuration methods
- [ ] Document request/response types
- [ ] Document error handling patterns

### Findings

#### Client Structure
```go
// To be filled in
```

#### Available Methods
```go
// To be filled in
```

#### Request/Response Types
```go
// To be filled in
```

#### Error Patterns
```go
// To be filled in
```

### Notes
_Investigation in progress..._

---

## Common Patterns Identified

### Error Handling
_To be documented..._

### Request/Response Patterns
_To be documented..._

### Context Usage
_To be documented..._

---

## Next Steps

1. Complete Management service investigation
2. Complete Talos service investigation
3. Complete Auth service investigation (if needed)
4. Complete OIDC service investigation (if needed)
5. Create API reference document
6. Create integration guide with examples

---

*Last Updated: 2025-01-27*
