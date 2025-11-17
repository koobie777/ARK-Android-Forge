# 🛰️ ARKFORGE - The ARK Ecosystem Android Development Module

<div align="center">
  
  ```
  ╭─────────────────────────────────────────────╮
  │       🛰️ THE ARK ECOSYSTEM SUPREME         │
  │          ARKFORGE Android Builder           │
  │           Commander: koobie777              │
  │              Version: 1.1.4                 │
  ╰─────────────────────────────────────────────╯
  ```

  **The ARK Ecosystem** - A modular homelab command center  
  **ARKFORGE** - Android development and ROM building orchestrator

</div>

## 🌌 Overview

ARKFORGE is the Android development module of The ARK Ecosystem, designed to streamline and automate the Android ROM building process. Built with a space-themed interface and tmux-powered session management, ARKFORGE makes Android development accessible for both expert builders and newcomers (Cadets).

## Go Rewrite Command Deck

The legacy bash orchestrator now has a Go-based successor named `ark-android-forge`.  
This new binary keeps the ARK theming while adding typed configuration (`forge.yaml`), Cobra commands, and structured logging.

### Build & Run

```bash
go build ./cmd/ark-android-forge

# Run the interactive deck
./ark-android-forge menu

# Run sub-commands directly
./ark-android-forge preflight
./ark-android-forge sync --manifest default.xml --force
./ark-android-forge build --device waffle --target recovery
./ark-android-forge release --output artifacts/manifest.yaml
```

### Configuration

`forge.yaml` at the repository root replaces `config/ark-settings.conf` and is loaded via Viper (env vars still override fields):

```yaml
version: "1.1.4"
commander: "koobie777"
mode: "expert"
jobs: 8
build:
  workspace: "./builds"
  defaultType: "recovery"
theme:
  enabled: true
  accent: "cyan"
fleet:
  - name: "OnePlus 12"
    codename: "waffle"
    role: "primary"
    repository: "lineageos"
```

### The ARK Ecosystem Components:
- **ARKFORGE** (This Module) - Android ROM/Recovery building
- **Orbital Command** - System management and monitoring (Coming Soon)
- Additional modules in development

## 🚀 Features

### Core Capabilities
- 🔧 **Smart Build System** - Automated device discovery → repo selection → build
- 📱 **Multi-Device Support** - Manage multiple Android devices in your fleet
- 🔄 **Persistent Sessions** - tmux-powered build sessions that survive disconnections
- 🛠️ **Modular Architecture** - Easy to extend with new modules
- 🎨 **ARK Themed Interface** - Consistent space-themed UI across all modules

### Current Modules
1. **Smart Build** - Intelligent build automation
2. **Recovery Build** - TWRP/OrangeFox builder (In Development)
3. **ROM Build** - Full ROM compilation
4. **Boot/Recovery Images** - Extract and build kernel images
5. **Resume Build** - Continue interrupted builds
6. **Repo Sync Only** - Smart repository synchronization
7. **Device Manager** - Fleet device management (In Development)
8. **Repository Manager** - ROM source management
9. **Directory Manager** - Build directory organization
10. **Configuration Manager** - ARK settings customization
11. **Fleet Status** - View all ARK devices
12. **Tmux Manager** - Session management

## 📋 Requirements

- **OS**: Linux (Tested on Arch Linux)
- **Shell**: Bash 4.0+
- **Dependencies**:
  - `tmux` - Terminal multiplexer for persistent sessions
  - `git` - Version control
  - `repo` - Android repo tool
  - Android build environment prerequisites
- **Storage**: 200GB+ recommended for ROM building
- **RAM**: 16GB+ recommended

## 🛸 Installation

```bash
# Clone The ARK Ecosystem ARKFORGE module
git clone https://github.com/koobie777/Ark-Forge.git ~/the-ark-ecosystem/ark-forge
cd ~/the-ark-ecosystem/ark-forge

# Make scripts executable
chmod +x arkforge-launcher arkforge.sh
chmod +x modules/*.sh

# Create config directory if needed
mkdir -p config

# Launch The ARK
./arkforge-launcher
```

## 🎮 Usage

### Quick Start
```bash
# Launch with tmux session management (Recommended)
./arkforge-launcher

# Or run directly without tmux
./arkforge.sh
```

### ARK Fleet Devices (Example)
- **Primary**: OnePlus 12 "waffle"
- **Secondary**: OnePlus 10 Pro "op515dl1"

### tmux Session Management
When using `arkforge-launcher`, builds run in persistent tmux sessions:
- **Detach**: `Ctrl-B`, then `D`
- **Reattach**: `./arkforge-launcher` and select option 1
- **List Windows**: `Ctrl-B`, then `W`
- **Switch Windows**: `Ctrl-B`, then window number

### Build Directory Structure
```
/home/koobie777/android/
├── lineageos-waffle/       # LineageOS for OnePlus 12
├── evolution-op515dl1/     # Evolution X for OnePlus 10 Pro
└── twrp-waffle/           # TWRP recovery for OnePlus 12
```

## 🛰️ ARK Modes

### Expert Mode (Default)
- Minimal prompts
- Fast operations
- Assumes familiarity with Android building

### Cadet Mode (Coming Soon)
- Guided experience
- Detailed explanations
- Step-by-step assistance

## 🔧 Configuration

ARK settings are stored in `config/ark-settings.conf`:
```bash
# The ARK Configuration
ARK_THEME_ENABLED=true
ARK_MODE="expert"  # or "cadet"
ARK_BUILD_DIR="/home/koobie777/android"
ARK_DEFAULT_JOBS=$(nproc --all)
```

## 📚 Module Development

Create new ARK modules:
```bash
# Module template
cat > modules/ark-new-module.sh << 'EOF'
#!/usr/bin/env bash
# ╭─────────────────────────────────────────────╮
# │         ARK MODULE NAME                     │
# │              Description                    │
# │           Commander: koobie777              │
# ╰─────────────────────────────────────────────╯

source "$ARK_CONFIG_DIR/ark-settings.conf"

# Module code here
EOF
```

## 🐛 Troubleshooting

### Common Issues

**tmux sessions persisting:**
```bash
# List all tmux sessions
tmux ls

# Kill specific session
tmux kill-session -t arkforge
```

**Module not loading:**
- Check file permissions: `chmod +x modules/*.sh`
- Verify path in arkforge.sh

**Build failures:**
- Check tmux windows: `tmux attach -t arkforge`
- Review logs in build directory

## 🚀 Roadmap

### Version 1.2.0 (Planned)
- [ ] Complete Recovery Build module
- [ ] Implement Cadet Mode
- [ ] Add build statistics/logging
- [ ] Device tree automation

### Version 2.0.0 (Future)
- [ ] Integration with Orbital Command
- [ ] Fleet-wide build orchestration
- [ ] Cloud backup integration
- [ ] AI-assisted troubleshooting

## 🤝 Contributing

The ARK Ecosystem welcomes contributions! Please ensure:
- Maintain ARK theming and style
- Test modules thoroughly
- Update documentation
- Follow the modular architecture

## 📜 License

This project is part of The ARK Ecosystem  
Commander: koobie777  
System: arksupreme-mk1

---

<div align="center">
  
  **🛰️ May The ARK be with you! 🛰️**
  
  Built with ❤️ for the Android development community
  
</div>
