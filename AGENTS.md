# Repository Guidelines

## Project Structure & Module Organization

* **Repository**: ARK-Android-Forge (Go orchestrator)
* `cmd/ark-android-forge/` — CLI entry (cobra commands: preflight, sync, build, release)
* `internal/android/` — wrappers for `repo`, `envsetup`, `lunch`, `m/mka`
* `internal/preflight/` — checks (java, repo, ulimit, ccache, disk)
* `internal/execx/` — robust exec (timeouts/retries, log tee, ANSI strip)
* `internal/artifacts/` — paths, hashing, manifests
* `internal/config/` — profile, manifests/templates
* `scripts/` — thin POSIX shims & examples
* `tests/` — unit + e2e harness
* `ci/` — GitHub Actions
* `.docs/` — setup, guides, release notes

## Build, Test, and Development Commands

* Go toolchain: **Go 1.22+**
* Build (native): `go build ./cmd/ark-android-forge`
* Build (Windows exe from Linux): `GOOS=windows GOARCH=amd64 go build -o bin/ark-android-forge.exe ./cmd/ark-android-forge`
* Test: `go test ./...`
* Lint/format: `gofmt -s -w .` • `go vet ./...`
* Release (optional): `goreleaser release --clean`
* Run (WSL autodetect): `./ark-android-forge preflight` • `./ark-android-forge sync` • `./ark-android-forge build --device <codename>`

## Coding Style & Naming Conventions

* Idiomatic Go; modules enabled; errors wrapped with context (`fmt.Errorf("…: %w", err)`).
* Packages: lower_snake; exported names PascalCase; unexported camelCase.
* Logs: structured with zerolog; support `--json`.
* Config: viper (env + file), default file `forge.yaml` at repo root.

## Testing Guidelines

* Unit tests for exec, preflight, config, path translation.
* e2e smoke via container/WSL: preflight→sync→envsetup→lunch (no full build by default).
* Attach logs + checksums for any artifact-producing test.

## Commit & Pull Request Guidelines

* **Conventional Commits** (e.g., `feat(build): add parallel log tee (#123)`).
* One logical change per PR; include description, linked issue, repro steps, risk/rollback.
* PR must pass CI (lint, unit, e2e smoke) and update docs when behavior changes.

## Security & Configuration Tips

* Do **not** commit proprietary blobs/keys/credentials. Use extraction scripts locally.
* No elevation required. Never write outside the repo tree by default.
* Redact personal paths and tokens in logs.

## Agent‑Specific Instructions (for AI/Copilot)

Provide in this order: (1) complete file(s) or minimal diff, (2) one‑liner patch, (3) paths + exact apply steps, (4) `UPDATE.md` snippet (semver + bullets), (5) brief reasoning + rollback. Prefer `--dry-run` where destructive; produce reproducible commands for WSL and native Linux/Windows.
