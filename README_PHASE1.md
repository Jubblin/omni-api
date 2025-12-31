# Phase 1: Investigation - Getting Started

## Quick Start

To begin Phase 1 investigation, follow these steps:

### 1. Run Investigation Commands

```bash
# Investigate Management service
go doc github.com/siderolabs/omni/client/pkg/management
go doc -all github.com/siderolabs/omni/client/pkg/management

# Investigate Talos service
go doc github.com/siderolabs/omni/client/pkg/talos
go doc -all github.com/siderolabs/omni/client/pkg/talos

# Investigate Auth service
go doc github.com/siderolabs/omni/client/pkg/auth
go doc -all github.com/siderolabs/omni/client/pkg/auth

# Investigate OIDC service
go doc github.com/siderolabs/omni/client/pkg/oidc
go doc -all github.com/siderolabs/omni/client/pkg/oidc
```

### 2. Use Investigation Scripts

```bash
# Run the shell script
./scripts/inspect_omni_client.sh

# Or use the Go investigation tool (requires client setup)
go run scripts/investigate_services.go
```

### 3. Document Findings

Update `PHASE1_INVESTIGATION.md` with your findings as you investigate each service.

## Files Created

1. **PHASE1_INVESTIGATION.md** - Investigation results (to be filled in)
2. **PHASE1_INVESTIGATION_GUIDE.md** - Detailed investigation guide
3. **scripts/inspect_omni_client.sh** - Shell script for investigation
4. **scripts/investigate_services.go** - Go tool for investigation

## Priority Order

1. **Management Service** (Start here - highest priority)
2. **Talos Service** (Second priority)
3. **Auth Service** (If needed)
4. **OIDC Service** (If needed)

## What to Document

For each service, document:
- Client type and structure
- All public methods
- Method signatures (parameters, return types)
- Request/Response types
- Error handling patterns
- Usage examples

## Next Steps

Once investigation is complete:
1. Create API reference document
2. Create integration examples
3. Begin Phase 2: Management Service Integration

---

*See PHASE1_INVESTIGATION_GUIDE.md for complete instructions*
