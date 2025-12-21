# SLSA Compliance

> 📖 [Back to README](../README.md)

This document describes the SLSA (Supply-chain Levels for Software Artifacts) compliance implementation for the Omni API project.

## SLSA Overview

SLSA is a security framework that provides a way to ensure the integrity of software artifacts throughout the supply chain. It defines four levels (0-3), with Level 3 being the highest assurance.

## Current SLSA Level

**SLSA Build Level 2+** ✅

The project implements SLSA Build Level 2 requirements with some Level 3 features:

### SLSA Build Level 2 Requirements ✅

- ✅ **Version Controlled**: Source code in version control (Git)
- ✅ **Verified History**: Git history preserved and verified
- ✅ **Retained Logs**: Build logs retained in GitHub Actions
- ✅ **Non-Falsifiable Provenance**: Provenance generated using GitHub OIDC
- ✅ **Dependencies Complete**: All dependencies declared in `go.mod`
- ✅ **Ephemeral Environment**: Builds run in ephemeral GitHub Actions runners
- ✅ **Isolated**: Builds run in isolated containers/environments
- ✅ **Provenance Available**: Provenance attestations generated and uploaded

### SLSA Build Level 3 Features ✅

- ✅ **Hardened Build Environment**: Uses GitHub-hosted runners with isolation
- ✅ **Provenance Generation**: Full provenance with `mode=max` in Docker Buildx
- ✅ **Artifact Attestations**: GitHub Artifact Attestations API integration
- ✅ **Cryptographic Signing**: Container images signed with cosign
- ✅ **SBOM Generation**: Software Bill of Materials for each architecture

## Implementation Details

### Provenance Generation

The build process generates SLSA-compliant provenance in multiple ways:

1. **Docker Buildx Provenance** (`mode=max`):
   - Full build provenance including all build steps
   - Source code information
   - Build environment details
   - Dependency information

2. **GitHub Artifact Attestations**:
   - Uses `actions/attest-build-provenance@v2`
   - Links artifacts to workflow runs
   - Includes repository, commit, and build information
   - Cryptographically signed using GitHub OIDC

3. **Container Signing**:
   - Images signed with cosign (keyless signing)
   - Verifiable signatures in registry
   - Links to provenance

### Provenance Contents

The generated provenance includes:

- **Build Information**:
  - Workflow run ID and URL
  - Build timestamp
  - Build environment details
  
- **Source Information**:
  - Repository URL
  - Commit SHA
  - Branch/tag information
  - Source code digest
  
- **Artifact Information**:
  - Container image digest
  - Image tags
  - Architecture (AMD64, ARM64)
  
- **Dependencies**:
  - Go module dependencies
  - Base image information
  - Build tool versions

### Verification

Provenance can be verified using:

```bash
# Verify GitHub Artifact Attestations
gh attestation verify \
  --owner OWNER \
  --repo REPO \
  --subject-name ghcr.io/OWNER/omni-api:latest

# Verify container signature
cosign verify --registry ghcr.io \
  ghcr.io/OWNER/omni-api:latest

# Verify Docker Buildx provenance
docker buildx imagetools inspect \
  ghcr.io/OWNER/omni-api:latest \
  --format '{{ json .Provenance }}'
```

## SLSA Requirements Checklist

### Build Requirements

- [x] **Scripted Build**: Build process defined in workflow files
- [x] **Version Controlled**: All build scripts in Git
- [x] **Build Service**: Uses GitHub Actions (hosted service)
- [x] **Build as Code**: Workflow defined in `.github/workflows/`
- [x] **Ephemeral Environment**: Each build runs in fresh environment
- [x] **Isolated**: Builds isolated from external influences
- [x] **Provenance**: Generated for all artifacts
- [x] **Non-Falsifiable**: Uses OIDC for authentication

### Source Requirements

- [x] **Version Controlled**: Source in Git repository
- [x] **Verified History**: Git history preserved
- [x] **Retained Indefinitely**: Repository history retained
- [x] **Two-Person Review**: PRs require review (if configured)
- [x] **Superhuman Access**: No direct push to main (if protected)

### Dependencies Requirements

- [x] **Defined**: All dependencies in `go.mod`
- [x] **Locked**: Dependencies locked in `go.sum`
- [x] **Available**: Dependencies available from public sources
- [x] **Verified**: Dependencies verified during build

## Achieving SLSA Build Level 3

To achieve full SLSA Build Level 3, the following would be required:

1. **Reusable Workflows**: Use reusable workflows for build isolation
2. **Build Service**: Use a dedicated build service (e.g., SLSA GitHub Generator)
3. **Hardened Builders**: Use SLSA-compliant build tools
4. **Additional Attestations**: Generate additional attestation types

### Current Status

The project implements most Level 3 requirements:
- ✅ Hardened build environment (GitHub Actions)
- ✅ Non-falsifiable provenance (OIDC)
- ✅ Full provenance generation
- ✅ Artifact attestations
- ⚠️ Reusable workflows (can be added for full Level 3)

## Viewing Attestations

### GitHub UI

1. Go to repository
2. Navigate to **Packages** → Container image
3. View **Attestations** tab
4. See provenance and other attestations

### Command Line

```bash
# List attestations
gh attestation list \
  --owner OWNER \
  --repo REPO \
  --subject-name ghcr.io/OWNER/omni-api:latest

# Verify attestation
gh attestation verify \
  --owner OWNER \
  --repo REPO \
  --subject-name ghcr.io/OWNER/omni-api:latest \
  --bundle attestation.bundle
```

## References

- [SLSA Specification](https://slsa.dev/spec/v1.0/)
- [GitHub Artifact Attestations](https://docs.github.com/en/actions/security-guides/using-artifact-attestations)
- [SLSA GitHub Generator](https://github.com/slsa-framework/slsa-github-generator)
- [Docker Buildx Provenance](https://docs.docker.com/build/attestations/)

## Related Documentation

- 📖 [README.md](../README.md) - Main project documentation
- 🔒 [.github/workflows/README.md](../.github/workflows/README.md) - CI/CD workflow documentation
- 🐳 [README.Docker.md](../README.Docker.md) - Docker build and deployment guide
- 📋 [RESOURCES.md](../RESOURCES.md) - Available Omni resources
- 📊 [TEST_COVERAGE.md](../TEST_COVERAGE.md) - Test coverage report
