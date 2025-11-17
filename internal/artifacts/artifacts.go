package artifacts

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/koobie777/ark-android-forge/internal/config"
)

// Manifest captures high-level release metadata.
type Manifest struct {
	GeneratedAt time.Time             `yaml:"generatedAt"`
	Commander   string                `yaml:"commander"`
	Version     string                `yaml:"version"`
	Devices     []config.FleetDevice  `yaml:"devices"`
	Notes       map[string]string     `yaml:"notes,omitempty"`
}

// Generate builds a manifest from the current configuration.
func Generate(cfg *config.Config) Manifest {
	return Manifest{
		GeneratedAt: time.Now().UTC(),
		Commander:   cfg.Commander,
		Version:     cfg.Version,
		Devices:     cfg.Fleet,
	}
}

// Write stores the manifest on disk.
func Write(path string, manifest Manifest) error {
	if path == "" {
		path = "artifacts/manifest.yaml"
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create manifest dir: %w", err)
		}
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}
