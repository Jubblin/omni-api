#!/bin/bash
# Version management script for semantic versioning

set -e

VERSION_FILE="VERSION"

# Get current version
get_version() {
    if [ -f "$VERSION_FILE" ]; then
        cat "$VERSION_FILE" | tr -d '[:space:]'
    else
        echo "0.0.1"
    fi
}

# Increment patch version
increment_patch() {
    local version=$(get_version)
    IFS='.' read -r major minor patch <<< "$version"
    patch=$((patch + 1))
    echo "${major}.${minor}.${patch}"
}

# Increment minor version
increment_minor() {
    local version=$(get_version)
    IFS='.' read -r major minor patch <<< "$version"
    minor=$((minor + 1))
    patch=0
    echo "${major}.${minor}.${patch}"
}

# Increment major version
increment_major() {
    local version=$(get_version)
    IFS='.' read -r major minor patch <<< "$version"
    major=$((major + 1))
    minor=0
    patch=0
    echo "${major}.${minor}.${patch}"
}

# Set version
set_version() {
    local new_version=$1
    echo "$new_version" > "$VERSION_FILE"
    
    # Update main.go
    if [ -f "main.go" ]; then
        sed -i.bak "s/@version[[:space:]]*[0-9]\+\.[0-9]\+\.[0-9]\+/@version         $new_version/" main.go
        rm -f main.go.bak
    fi
    
    # Update health.go
    if [ -f "internal/api/handlers/health.go" ]; then
        sed -i.bak "s/Version:[[:space:]]*\"[0-9]\+\.[0-9]\+\.[0-9]\+\"/Version:   \"$new_version\"/" internal/api/handlers/health.go
        rm -f internal/api/handlers/health.go.bak
    fi
    
    # Regenerate Swagger docs if swag is available
    if command -v swag &> /dev/null; then
        echo "Regenerating Swagger documentation..."
        swag init
        echo "Swagger documentation regenerated"
    else
        echo "Note: swag not found. Run 'make swagger' to regenerate documentation."
    fi
    
    echo "Version updated to $new_version"
}

# Main command handler
case "${1:-}" in
    get)
        get_version
        ;;
    patch)
        set_version $(increment_patch)
        ;;
    minor)
        set_version $(increment_minor)
        ;;
    major)
        set_version $(increment_major)
        ;;
    set)
        if [ -z "$2" ]; then
            echo "Usage: $0 set <version>"
            exit 1
        fi
        set_version "$2"
        ;;
    *)
        echo "Usage: $0 {get|patch|minor|major|set <version>}"
        exit 1
        ;;
esac
