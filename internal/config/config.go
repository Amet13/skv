// Package config handles YAML configuration loading and validation.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
	"skv/internal/provider"
)

// Config is the top-level configuration.
type Config struct {
	Defaults Defaults `yaml:"defaults"` // Global default parameters
	Secrets  []Secret `yaml:"secrets"`  // List of secrets to manage
}

// Defaults holds global default parameters merged into each secret unless overridden.
type Defaults struct {
	Region  string            `yaml:"region"`  // Default AWS region or similar
	Address string            `yaml:"address"` // Default server address (e.g., Vault URL)
	Token   string            `yaml:"token"`   // Default authentication token
	Extras  map[string]string `yaml:"extras"`  // Provider-specific defaults
}

// Secret represents a single secret to fetch and where to place it.
type Secret struct {
	Alias    string            `yaml:"alias"`    // Human-readable identifier
	Provider string            `yaml:"provider"` // Provider type (aws, gcp, etc.)
	Name     string            `yaml:"name"`     // Provider-specific secret path/name
	Env      string            `yaml:"env"`      // Environment variable name
	Region   string            `yaml:"region"`   // Provider region (AWS, GCP zones, etc.)
	Address  string            `yaml:"address"`  // Provider address (Vault URL, etc.)
	Token    string            `yaml:"token"`    // Authentication token
	Path     string            `yaml:"path"`     // Secret path (for Vault-like providers)
	Version  *int              `yaml:"version"`  // Secret version (if supported)
	Metadata map[string]string `yaml:"metadata"` // Additional metadata
	Extras   map[string]string `yaml:"extras"`   // Provider-specific options
}

// Load reads the configuration from file, applying env interpolation and validation.
func Load(overridePath string) (*Config, error) {
	path := locateConfigPath(overridePath)
	if path == "" {
		return nil, errors.New("no config file found; set --config or SKV_CONFIG or create ~/.skv.yaml")
	}

	// #nosec G304: path is sourced from flags/env/home and is expected to be a user-provided file path
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Interpolate environment variables in defaults
	cfg.Defaults.Region = interpolateEnv(cfg.Defaults.Region)
	cfg.Defaults.Address = interpolateEnv(cfg.Defaults.Address)
	cfg.Defaults.Token = interpolateEnv(cfg.Defaults.Token)
	if cfg.Defaults.Extras != nil {
		for k, v := range cfg.Defaults.Extras {
			cfg.Defaults.Extras[k] = interpolateEnv(v)
		}
	}

	// Interpolate environment variables in all string fields.
	for i := range cfg.Secrets {
		s := &cfg.Secrets[i]
		s.Alias = interpolateEnv(s.Alias)
		s.Provider = interpolateEnv(s.Provider)
		s.Name = interpolateEnv(s.Name)
		s.Env = interpolateEnv(s.Env)
		s.Region = interpolateEnv(s.Region)
		s.Address = interpolateEnv(s.Address)
		s.Token = interpolateEnv(s.Token)
		s.Path = interpolateEnv(s.Path)
		// Metadata values
		if s.Metadata != nil {
			for k, v := range s.Metadata {
				s.Metadata[k] = interpolateEnv(v)
			}
		}
		// Extras values
		if s.Extras != nil {
			for k, v := range s.Extras {
				s.Extras[k] = interpolateEnv(v)
			}
		}
	}

	// After interpolation, merge defaults into secrets
	for i := range cfg.Secrets {
		s := &cfg.Secrets[i]
		if s.Region == "" && cfg.Defaults.Region != "" {
			s.Region = cfg.Defaults.Region
		}
		if s.Address == "" && cfg.Defaults.Address != "" {
			s.Address = cfg.Defaults.Address
		}
		if s.Token == "" && cfg.Defaults.Token != "" {
			s.Token = cfg.Defaults.Token
		}
		if cfg.Defaults.Extras != nil {
			if s.Extras == nil {
				s.Extras = map[string]string{}
			}
			for k, v := range cfg.Defaults.Extras {
				if _, exists := s.Extras[k]; !exists {
					s.Extras[k] = v
				}
			}
		}
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func locateConfigPath(overridePath string) string {
	if overridePath != "" {
		return overridePath
	}
	if env := os.Getenv("SKV_CONFIG"); env != "" {
		return env
	}
	// Home fallbacks
	home, err := os.UserHomeDir()
	if err == nil {
		candidates := []string{
			filepath.Join(home, ".skv.yaml"),
			filepath.Join(home, ".skv.yml"),
		}
		for _, p := range candidates {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	return ""
}

func (c *Config) validate() error {
	if len(c.Secrets) == 0 {
		return errors.New("config.secrets is empty")
	}
	aliases := map[string]struct{}{}
	for _, s := range c.Secrets {
		if s.Alias == "" {
			return errors.New("secret.alias is required")
		}
		if _, dup := aliases[s.Alias]; dup {
			return fmt.Errorf("duplicate alias: %s", s.Alias)
		}
		aliases[s.Alias] = struct{}{}
		if s.Provider == "" {
			return fmt.Errorf("provider is required for alias %s", s.Alias)
		}
		if s.Name == "" {
			return fmt.Errorf("name is required for alias %s", s.Alias)
		}
		// Fail fast if interpolation left missing env tokens
		if containsMissingEnvToken(s.Alias, s.Provider, s.Name, s.Env, s.Region, s.Address, s.Token, s.Path) {
			return fmt.Errorf("missing environment variable in configuration for alias %s", s.Alias)
		}
		// Do not enforce provider registration here to keep config loading
		// decoupled from runtime registrations. Unknown providers will be
		// handled at command execution time.
	}
	return nil
}

// FindByAlias returns the secret with the given alias, or nil if not found.
func (c *Config) FindByAlias(alias string) (*Secret, bool) {
	for i := range c.Secrets {
		if c.Secrets[i].Alias == alias {
			return &c.Secrets[i], true
		}
	}
	return nil, false
}

// ToSpec converts a Secret to a provider.SecretSpec for use with providers.
func (s Secret) ToSpec() provider.SecretSpec {
	envName := s.Env
	if strings.TrimSpace(envName) == "" {
		envName = deriveEnvName(s.Alias)
	}
	// Start with built-in mapped fields
	extras := map[string]string{}
	if s.Region != "" {
		extras["region"] = s.Region
	}
	if s.Address != "" {
		extras["address"] = s.Address
	}
	if s.Token != "" {
		extras["token"] = s.Token
	}
	if s.Path != "" {
		extras["path"] = s.Path
	}
	if s.Version != nil {
		extras["version"] = fmt.Sprintf("%d", *s.Version)
	}
	// Merge metadata (back-compat) without overriding built-ins
	for k, v := range s.Metadata {
		if _, exists := extras[k]; !exists {
			extras[k] = v
		}
	}
	// Merge explicit extras with highest precedence
	for k, v := range s.Extras {
		extras[k] = v
	}
	return provider.SecretSpec{
		Alias:    s.Alias,
		Name:     s.Name,
		Provider: s.Provider,
		EnvName:  envName,
		Extras:   extras,
	}
}

var envTpl = regexp.MustCompile(`\{\{\s*([A-Za-z_][A-Za-z0-9_]*)\s*\}\}`)

func interpolateEnv(s string) string {
	if s == "" {
		return s
	}
	return envTpl.ReplaceAllStringFunc(s, func(m string) string {
		name := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(m, "{{"), "}}"))
		name = strings.TrimSpace(name)
		if name == "" {
			return m
		}
		val, ok := os.LookupEnv(name)
		if !ok {
			// For MVP, missing env is an error. We mark it with a sentinel to fail on validation.
			// Use a distinctive token to simplify debugging.
			return fmt.Sprintf("__MISSING_ENV_%s__", name)
		}
		return val
	})
}

func deriveEnvName(alias string) string {
	b := strings.Builder{}
	prevUnderscore := false
	for _, r := range alias {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			ch := r
			if ch >= 'a' && ch <= 'z' {
				ch = ch - 'a' + 'A'
			}
			b.WriteRune(ch)
			prevUnderscore = false
		} else {
			if !prevUnderscore {
				b.WriteByte('_')
				prevUnderscore = true
			}
		}
	}
	res := b.String()
	res = strings.Trim(res, "_")
	res = strings.ReplaceAll(res, "__", "_")
	if res == "" {
		return "SECRET"
	}
	return res
}

func containsMissingEnvToken(values ...string) bool {
	for _, v := range values {
		if strings.Contains(v, "__MISSING_ENV_") {
			return true
		}
	}
	return false
}
