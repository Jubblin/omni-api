#!/bin/bash
# Script to investigate Omni client service APIs
# Run this script to document available methods

set -e

echo "=== Omni Client Service API Investigation ==="
echo ""

# Check if we can access the Go module
if [ -z "$GOPATH" ] && [ -z "$HOME/go" ]; then
    echo "Error: GOPATH not set"
    exit 1
fi

MODULE_PATH="${GOPATH:-$HOME/go}/pkg/mod/github.com/siderolabs/omni/client@v1.4.6"

echo "Looking for Omni client module at: $MODULE_PATH"
echo ""

if [ ! -d "$MODULE_PATH" ]; then
    echo "Module not found. Try running: go mod download"
    echo ""
    echo "Alternative: Use go doc to inspect packages:"
    echo "  go doc github.com/siderolabs/omni/client/pkg/management"
    echo "  go doc github.com/siderolabs/omni/client/pkg/talos"
    echo "  go doc github.com/siderolabs/omni/client/pkg/auth"
    echo "  go doc github.com/siderolabs/omni/client/pkg/oidc"
    exit 1
fi

echo "=== Management Service ==="
if [ -d "$MODULE_PATH/pkg/management" ]; then
    echo "Package found at: $MODULE_PATH/pkg/management"
    echo ""
    echo "Files:"
    ls -1 "$MODULE_PATH/pkg/management"/*.go 2>/dev/null | head -10
    echo ""
    echo "Use 'go doc' to see exported symbols:"
    echo "  go doc github.com/siderolabs/omni/client/pkg/management"
else
    echo "Package not found"
fi

echo ""
echo "=== Talos Service ==="
if [ -d "$MODULE_PATH/pkg/talos" ]; then
    echo "Package found at: $MODULE_PATH/pkg/talos"
    echo ""
    echo "Files:"
    ls -1 "$MODULE_PATH/pkg/talos"/*.go 2>/dev/null | head -10
    echo ""
    echo "Use 'go doc' to see exported symbols:"
    echo "  go doc github.com/siderolabs/omni/client/pkg/talos"
else
    echo "Package not found"
fi

echo ""
echo "=== Auth Service ==="
if [ -d "$MODULE_PATH/pkg/auth" ]; then
    echo "Package found at: $MODULE_PATH/pkg/auth"
    echo ""
    echo "Files:"
    ls -1 "$MODULE_PATH/pkg/auth"/*.go 2>/dev/null | head -10
    echo ""
    echo "Use 'go doc' to see exported symbols:"
    echo "  go doc github.com/siderolabs/omni/client/pkg/auth"
else
    echo "Package not found"
fi

echo ""
echo "=== OIDC Service ==="
if [ -d "$MODULE_PATH/pkg/oidc" ]; then
    echo "Package found at: $MODULE_PATH/pkg/oidc"
    echo ""
    echo "Files:"
    ls -1 "$MODULE_PATH/pkg/oidc"/*.go 2>/dev/null | head -10
    echo ""
    echo "Use 'go doc' to see exported symbols:"
    echo "  go doc github.com/siderolabs/omni/client/pkg/oidc"
else
    echo "Package not found"
fi

echo ""
echo "=== Investigation Commands ==="
echo ""
echo "To investigate each service, run:"
echo ""
echo "# Management Service"
echo "go doc github.com/siderolabs/omni/client/pkg/management"
echo "go doc -all github.com/siderolabs/omni/client/pkg/management"
echo ""
echo "# Talos Service"
echo "go doc github.com/siderolabs/omni/client/pkg/talos"
echo "go doc -all github.com/siderolabs/omni/client/pkg/talos"
echo ""
echo "# Auth Service"
echo "go doc github.com/siderolabs/omni/client/pkg/auth"
echo "go doc -all github.com/siderolabs/omni/client/pkg/auth"
echo ""
echo "# OIDC Service"
echo "go doc github.com/siderolabs/omni/client/pkg/oidc"
echo "go doc -all github.com/siderolabs/omni/client/pkg/oidc"
