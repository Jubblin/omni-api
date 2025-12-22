# GitHub Actions Workflows

This directory contains GitHub Actions workflows for building, testing, and securing the Omni API container images.

## Workflows

### build-and-push.yml

**Purpose**: Builds multi-architecture container images (AMD64 and ARM64) and pushes them to GitHub Container Registry (GHCR) with cryptographic signatures.

**Triggers**:

- Push to `main` or `master` branch
- Push of version tags (e.g., `v0.0.1`)
- Pull requests (builds but doesn't push)
- Manual workflow dispatch

**Features**:

- ✅ Dockerfile scanning with Hadolint and Checkov (before build)
- ✅ Multi-architecture builds (linux/amd64, linux/arm64)
- ✅ Container scanning with Trivy (after build, per architecture)
- ✅ SBOM (Software Bill of Materials) generation per architecture
- ✅ Container signing with cosign (keyless signing)
- ✅ Provenance attestation
- ✅ Build caching for faster builds
- ✅ Automatic version tagging
- ✅ Security scanning integration

**Outputs**:

- Container images pushed to `ghcr.io/<owner>/<repo>`
- Signed container images
- SBOM artifacts (separate for AMD64 and ARM64)
- Trivy scan results (SARIF format, uploaded to GitHub Security)
- Build summaries with scan status

**Tags Created**:

- Branch name (e.g., `main`)
- Commit SHA (e.g., `main-abc1234`)
- Semantic version (e.g., `0.0.1`, `0.0`, `0`)
- `latest` (for default branch)

### version-bump.yml

**Purpose**: Automatically increments patch version numbers for pull requests targeting the main branch.

**Triggers**:

- Pull request opened to `main` or `master`
- Pull request synchronized (updated/rebased)
- Pull request reopened

**Features**:

- ✅ Automatic patch version increment
- ✅ Updates VERSION file
- ✅ Updates version in `main.go` (Swagger annotation)
- ✅ Updates version in `internal/api/handlers/health.go` (health endpoint)
- ✅ Regenerates Swagger documentation (`docs/`) with new version
- ✅ Adds version number as label to PR (e.g., `v0.0.2`)
- ✅ Handles rebases (checks if version already incremented)
- ✅ Commits changes back to PR branch
- ✅ Comments on PR with version update status

**How it works**:

1. Reads base version from `main` branch's `VERSION` file
2. Checks if PR branch already has an incremented version
3. If not, increments patch version (e.g., `0.0.1` → `0.0.2`)
4. Updates all version references in code
5. Commits and pushes changes to PR branch
6. Comments on PR with update status

**Version File**:

- The `VERSION` file in the repository root stores the current version
- Format: `MAJOR.MINOR.PATCH` (e.g., `0.0.1`)
- Can be manually updated using `make version-patch`, `make version-minor`, or `make version-major`

### dependabot-recreate.yml

**Purpose**: Automatically recreates Dependabot pull requests when the destination branch (main/master) is updated.

**Triggers**:

- Push to `main` or `master` branch

**Features**:

- ✅ Automatically finds all open Dependabot PRs
- ✅ Recreates PRs by commenting `@dependabot recreate`
- ✅ Prevents duplicate comments (checks for recent recreate comments)
- ✅ Handles rate limiting with delays between requests
- ✅ Logs all actions for debugging

**How it works**:

1. Triggers when code is pushed to `main` or `master`
2. Lists all open pull requests targeting the updated branch
3. Filters for Dependabot PRs (bot user)
4. Comments `@dependabot recreate` on each PR
5. Dependabot recreates the PR with the latest base branch changes

**Note**: This works in conjunction with `rebase-strategy: "auto"` in `dependabot.yml` to ensure PRs stay up-to-date with the destination branch.

### security-scan.yml

**Purpose**: Performs security scanning on container images and codebase.

**Triggers**:

- Push to `main` or `master` branch
- Push of version tags
- Pull requests
- Weekly schedule (Mondays at 00:00 UTC)

**Features**:

- ✅ Trivy container vulnerability scanning
- ✅ Trivy filesystem scanning
- ✅ SARIF upload to GitHub Security

**Scans**:

- Container image vulnerabilities
- Go dependencies
- Filesystem security issues

## Setup

### Required Secrets

No secrets are required for basic functionality. The workflow uses:

- `GITHUB_TOKEN` (automatically provided) for registry authentication
- Keyless signing with cosign (no keys needed)

### Optional Secrets

No optional secrets required.

### Permissions

The workflows require the following permissions:

- `contents: read` - Read repository contents
- `packages: write` - Push to GitHub Container Registry
- `id-token: write` - For keyless signing with cosign and SLSA attestations
- `security-events: write` - Upload security scan results
- `attestations: write` - Upload SLSA provenance attestations

These are configured in the workflow files.

## Usage

### Automatic Builds

1. **Push to main branch**: Automatically builds and pushes `latest` tag
2. **Create a tag**: `git tag v0.0.1 && git push origin v0.0.1`
   - Builds and pushes versioned images
3. **Create a PR**: Builds but doesn't push (for testing)

### Manual Builds

1. Go to Actions tab in GitHub
2. Select "Build and Push Container"
3. Click "Run workflow"
4. Enter version (e.g., `0.0.1`)
5. Click "Run workflow"

### Pulling Images

```bash
# Login to GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Pull latest
docker pull ghcr.io/jubblin/omni-api:latest

# Pull specific version
docker pull ghcr.io/jubblin/omni-api:0.0.1

# Pull for specific architecture
docker pull --platform linux/arm64 ghcr.io/jubblin/omni-api:latest
```

### Verifying Signatures

```bash
# Install cosign
brew install cosign  # macOS
# or download from https://github.com/sigstore/cosign/releases

# Verify signature
cosign verify --registry ghcr.io \
  ghcr.io/jubblin/omni-api:latest
```

## Security Features

### Dockerfile Scanning (Before Build)

Before building containers, Dockerfiles are scanned for security issues:

- **Hadolint**: Dockerfile linting and best practices
  - Scans both `Dockerfile` and `Dockerfile.alpine`
  - Checks for common Dockerfile anti-patterns
  - Results displayed in workflow logs
  
- **Checkov**: Infrastructure as Code security scanning
  - Scans Dockerfiles for security misconfigurations
  - Results uploaded to GitHub Security tab (SARIF format)
  - Build fails if critical issues found

### Container Scanning (After Build)

After building each architecture, containers are scanned:

- **Trivy**: Container vulnerability scanning
  - AMD64 container scanned separately with `--platform linux/amd64`
  - ARM64 container scanned separately with `--platform linux/arm64`
  - Scans for CRITICAL, HIGH, and MEDIUM severity vulnerabilities
  - Results uploaded to GitHub Security tab (SARIF format)
  - SBOMs generated for each architecture (SPDX JSON format)

### Container Signing

All container images are signed using cosign with keyless signing:

- Uses GitHub OIDC for authentication
- No private keys to manage
- Verifiable signatures in registry

### SBOM Generation

Software Bill of Materials is generated for each architecture:

- SPDX JSON format
- Separate SBOMs for AMD64 and ARM64
- Uploaded as workflow artifacts (90-day retention)
- Includes all dependencies and packages

### SLSA Compliance

The workflow implements **SLSA Build Level 2+** compliance:

- **SLSA Provenance Attestations**: Generated using GitHub Artifact Attestations API
  - Links artifacts to workflow runs
  - Includes repository, commit, and build information
  - Cryptographically signed using GitHub OIDC (non-falsifiable)
  
- **Docker Buildx Provenance**: Full provenance with `mode=max`
  - Complete build provenance including all build steps
  - Source code information
  - Build environment details
  - Dependency information
  
- **Build Requirements Met**:
  - ✅ Scripted build (workflow-defined)
  - ✅ Version controlled (in Git)
  - ✅ Ephemeral environment (GitHub Actions runners)
  - ✅ Isolated builds
  - ✅ Non-falsifiable provenance (OIDC-based)

See [.slsa/README.md](../.slsa/README.md) for detailed SLSA compliance documentation.

### Provenance

Build provenance is automatically generated in multiple formats:

- **GitHub Artifact Attestations**: SLSA-compliant provenance via `actions/attest-build-provenance`
- **Docker Buildx Provenance**: Full build provenance (mode=max)
- **Build information**: Workflow run, timestamp, environment
- **Source information**: Repository, commit SHA, branch/tag
- **Artifact information**: Image digest, tags, architectures

### Additional Security Scanning

All results uploaded to GitHub Security tab

## Troubleshooting

### Build Failures

1. Check workflow logs in Actions tab
2. Verify Dockerfile syntax
3. Check for dependency issues
4. Verify build arguments

### Signing Failures

1. Ensure `id-token: write` permission is set
2. Check cosign installation step
3. Verify registry authentication

### Multi-Arch Build Issues

1. Verify Docker Buildx is set up correctly
2. Check platform support in Dockerfile
3. Review build logs for architecture-specific errors

## Best Practices

1. **Always use version tags** for production deployments
2. **Verify signatures** before deploying
3. **Review security scans** regularly
4. **Use semantic versioning** for tags
5. **Keep workflows updated** with latest action versions

## References

- [Docker Buildx](https://docs.docker.com/buildx/)
- [cosign](https://github.com/sigstore/cosign)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Trivy](https://github.com/aquasecurity/trivy)
- [SBOM](https://spdx.dev/)
