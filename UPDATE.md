## 1.2.0-alpha.0

- replace the bash-only launcher with the new `ark-android-forge` Go binary (cobra commands + interactive menu)
- introduce typed configuration via `forge.yaml` with viper/env overrides and legacy `ark-settings.conf` fallback
- add modular packages (`internal/android`, `internal/preflight`, `internal/execx`, `internal/artifacts`, `internal/ui`)
- wire up repo sync, build scaffolding, release manifest generation, and zerolog-backed logging
