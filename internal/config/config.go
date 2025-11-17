package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config captures the runtime configuration for ARKFORGE.
type Config struct {
	File      string        `mapstructure:"-" yaml:"-"`
	Version   string        `mapstructure:"version" yaml:"version"`
	Commander string        `mapstructure:"commander" yaml:"commander"`
	Mode      string        `mapstructure:"mode" yaml:"mode"`
	Jobs      int           `mapstructure:"jobs" yaml:"jobs"`
	Build     BuildConfig   `mapstructure:"build" yaml:"build"`
	Theme     ThemeConfig   `mapstructure:"theme" yaml:"theme"`
	Fleet     []FleetDevice `mapstructure:"fleet" yaml:"fleet"`
}

// BuildConfig describes build defaults.
type BuildConfig struct {
	Workspace   string `mapstructure:"workspace" yaml:"workspace"`
	DefaultType string `mapstructure:"defaultType" yaml:"defaultType"`
}

// ThemeConfig controls TUI appearance.
type ThemeConfig struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
	Accent  string `mapstructure:"accent" yaml:"accent"`
}

// FleetDevice describes a device that can be built.
type FleetDevice struct {
	Name       string `mapstructure:"name" yaml:"name"`
	Codename   string `mapstructure:"codename" yaml:"codename"`
	Role       string `mapstructure:"role" yaml:"role"`
	Repository string `mapstructure:"repository" yaml:"repository"`
}

// Load returns the parsed configuration or a default if no file exists.
func Load(path string) (*Config, error) {
	if path == "" {
		path = "forge.yaml"
	}

	cfg, err := readModern(path)
	if err == nil {
		return cfg, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	legacyCfg, legacyErr := loadLegacy(filepath.Join("config", "ark-settings.conf"))
	if legacyErr == nil {
		if perr := persistDefault(path, legacyCfg); perr != nil {
			legacyCfg.File = path
			return legacyCfg, nil
		}
		legacyCfg.File = path
		return legacyCfg, nil
	}

	def := Default()
	if perr := persistDefault(path, def); perr != nil {
		return def, nil
	}
	def.File = path
	return def, nil
}

func readModern(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return nil, os.ErrNotExist
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if cfg.Jobs <= 0 {
		cfg.Jobs = Default().Jobs
	}
	if cfg.Build.Workspace == "" {
		cfg.Build.Workspace = Default().Build.Workspace
	}
	if cfg.Build.DefaultType == "" {
		cfg.Build.DefaultType = Default().Build.DefaultType
	}

	cfg.File = path
	return &cfg, nil
}

func persistDefault(path string, cfg *Config) error {
	if path == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create config dir: %w", err)
		}
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal default config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write default config: %w", err)
	}
	return nil
}

func loadLegacy(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	kv := map[string]string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.Trim(parts[0], "\" ")
		value := strings.Trim(parts[1], "\" ")
		value = strings.Trim(value, "\"")
		kv[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	cfg := Default()
	if commander, ok := kv["ARK_COMMANDER"]; ok {
		cfg.Commander = commander
	}
	if primary, ok := kv["ARK_PRIMARY_DEVICE"]; ok {
		cfg.Fleet[0].Codename = primary
	}
	if primaryName, ok := kv["ARK_PRIMARY_DEVICE_NAME"]; ok {
		cfg.Fleet[0].Name = primaryName
	}
	if defaultType, ok := kv["ARK_DEFAULT_BUILD_TYPE"]; ok {
		cfg.Build.DefaultType = defaultType
	}
	if jobsStr, ok := kv["ARK_DEFAULT_JOBS"]; ok {
		if jobs, err := strconv.Atoi(jobsStr); err == nil && jobs > 0 {
			cfg.Jobs = jobs
		}
	}
	return cfg, nil
}

// Default provides a sane baseline configuration.
func Default() *Config {
	return &Config{
		Version:   "1.1.4",
		Commander: "koobie777",
		Mode:      "expert",
		Jobs:      8,
		Build: BuildConfig{
			Workspace:   "./builds",
			DefaultType: "recovery",
		},
		Theme: ThemeConfig{
			Enabled: true,
			Accent:  "cyan",
		},
		Fleet: []FleetDevice{
			{
				Name:       "OnePlus 12",
				Codename:   "waffle",
				Role:       "primary",
				Repository: "lineageos",
			},
			{
				Name:       "OnePlus 10 Pro",
				Codename:   "op515dl1",
				Role:       "secondary",
				Repository: "evolution",
			},
		},
	}
}

func setDefaults(v *viper.Viper) {
	def := Default()
	v.SetDefault("version", def.Version)
	v.SetDefault("commander", def.Commander)
	v.SetDefault("mode", def.Mode)
	v.SetDefault("jobs", def.Jobs)
	v.SetDefault("build.workspace", def.Build.Workspace)
	v.SetDefault("build.defaultType", def.Build.DefaultType)
	v.SetDefault("theme.enabled", def.Theme.Enabled)
	v.SetDefault("theme.accent", def.Theme.Accent)
	v.SetDefault("fleet", def.Fleet)
}

// DeviceByCodename returns the fleet device matching the codename.
func (c *Config) DeviceByCodename(code string) *FleetDevice {
	for i := range c.Fleet {
		if strings.EqualFold(c.Fleet[i].Codename, code) {
			return &c.Fleet[i]
		}
	}
	return nil
}
