# Continuous Integration & Delivery

## Overview

To guarantee code quality and prevent the infamous "It works on my machine" syndrome, this project leverages a highly deterministic CI pipeline powered by **GitHub Actions** and **Nix**.

Every push to `main` and every Pull Request automatically triggers this pipeline to ensure the application compiles, lints, and passes all integration tests against a real database.

---

## The "Works on My Machine" Problem

Traditionally, CI pipelines run directly on an Ubuntu runner, installing Go via `actions/setup-go`. While this works, it can cause subtle environment discrepancies if the developer is building on macOS with a slightly different Go version or different native dependencies.

### The Nix Solution

Instead of relying on the host system's packages, we use a **Nix Flake** (`flake.nix`). 
Nix is a purely functional package manager. It provides an iron-clad guarantee that the development environment is exactly the same down to the byte, regardless of the host OS.

Our `flake.nix` specifies:
- The exact Go version (`go`)
- The exact migration tool (`go-migrate-pg`)
- Multi-architecture support (`x86_64-linux`, `aarch64-linux`, `x86_64-darwin`, `aarch64-darwin`)

---

## GitHub Actions Pipeline

Our pipeline (`.github/workflows/ci.yml`) executes the following steps:

1. **Checkout**: Retrieves the repository source code.
2. **Install Nix**: Uses `cachix/install-nix-action` to install the Nix package manager on the Ubuntu runner.
3. **Lint (`go vet`)**: Drops into the `nix develop` shell and runs standard Go static analysis to catch suspicious constructs.
4. **Build (`go build`)**: Verifies that the `cmd/api` application can successfully compile into an executable binary within the Nix environment.
5. **Integration Tests (`go test`)**: Executes the full Testcontainers test suite without caching (`-count=1`).

### Testcontainers on GitHub Actions

Because GitHub Actions Linux runners come with the Docker Daemon running natively, `testcontainers-go` works flawlessly without any additional configuration. The pipeline will automatically spin up the `postgres:16-alpine` container, run migrations, and execute the tests exactly as it does on a developer's local machine.
