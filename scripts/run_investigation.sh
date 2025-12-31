#!/bin/bash
# Phase 1 Investigation Script
# This script runs investigation commands and saves output to files

set -e

OUTPUT_DIR="investigation_results"
mkdir -p "$OUTPUT_DIR"

echo "=== Phase 1: Omni Client Service API Investigation ==="
echo ""

# Check if go is available
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

echo "Investigating Omni client services..."
echo "Results will be saved to: $OUTPUT_DIR/"
echo ""

# Management Service
echo "[1/4] Investigating Management Service..."
go doc github.com/siderolabs/omni/client/pkg/management > "$OUTPUT_DIR/management_basic.txt" 2>&1 || echo "Could not access Management service docs"
go doc -all github.com/siderolabs/omni/client/pkg/management > "$OUTPUT_DIR/management_all.txt" 2>&1 || echo "Could not access Management service docs (all)"
go doc github.com/siderolabs/omni/client/pkg/management.Client > "$OUTPUT_DIR/management_client.txt" 2>&1 || echo "Could not access Management Client docs"
echo "  ✓ Management service investigation complete"

# Talos Service
echo "[2/4] Investigating Talos Service..."
go doc github.com/siderolabs/omni/client/pkg/talos > "$OUTPUT_DIR/talos_basic.txt" 2>&1 || echo "Could not access Talos service docs"
go doc -all github.com/siderolabs/omni/client/pkg/talos > "$OUTPUT_DIR/talos_all.txt" 2>&1 || echo "Could not access Talos service docs (all)"
go doc github.com/siderolabs/omni/client/pkg/talos.Client > "$OUTPUT_DIR/talos_client.txt" 2>&1 || echo "Could not access Talos Client docs"
echo "  ✓ Talos service investigation complete"

# Auth Service
echo "[3/4] Investigating Auth Service..."
go doc github.com/siderolabs/omni/client/pkg/auth > "$OUTPUT_DIR/auth_basic.txt" 2>&1 || echo "Could not access Auth service docs"
go doc -all github.com/siderolabs/omni/client/pkg/auth > "$OUTPUT_DIR/auth_all.txt" 2>&1 || echo "Could not access Auth service docs (all)"
go doc github.com/siderolabs/omni/client/pkg/auth.Client > "$OUTPUT_DIR/auth_client.txt" 2>&1 || echo "Could not access Auth Client docs"
echo "  ✓ Auth service investigation complete"

# OIDC Service
echo "[4/4] Investigating OIDC Service..."
go doc github.com/siderolabs/omni/client/pkg/oidc > "$OUTPUT_DIR/oidc_basic.txt" 2>&1 || echo "Could not access OIDC service docs"
go doc -all github.com/siderolabs/omni/client/pkg/oidc > "$OUTPUT_DIR/oidc_all.txt" 2>&1 || echo "Could not access OIDC service docs (all)"
go doc github.com/siderolabs/omni/client/pkg/oidc.Client > "$OUTPUT_DIR/oidc_client.txt" 2>&1 || echo "Could not access OIDC Client docs"
echo "  ✓ OIDC service investigation complete"

echo ""
echo "=== Investigation Complete ==="
echo ""
echo "Results saved to: $OUTPUT_DIR/"
echo ""
echo "Next steps:"
echo "1. Review the generated documentation files"
echo "2. Update PHASE1_INVESTIGATION.md with findings"
echo "3. Create API reference document"
echo ""
echo "Files generated:"
ls -lh "$OUTPUT_DIR/" | tail -n +2
