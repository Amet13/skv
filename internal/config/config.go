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
	Secrets []Secret `yaml:"secrets"`
}

// Secret represents a single secret to fetch and where to place it.
type Secret struct {
	Alias    string            `yaml:"alias"`
	Provider string            `yaml:"provider"`
	Name     string            `yaml:"name"`
	Env      string            `yaml:"env"`
	Region   string            `yaml:"region"`
	Address  string            `yaml:"address"`
	Token    string            `yaml:"token"`
	Path     string            `yaml:"path"`
	Version  *int              `yaml:"version"`
	Metadata map[string]string `yaml:"metadata"`
}

// Load reads the configuration from file, applying env interpolation and validation.
func Load(overridePath string) (*Config, error) {
	path := overridePath
	if path == "" {
		if env := os.Getenv("SKV_CONFIG"); env != "" {
			path = env
		}
	}
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determine home dir: %w", err)
		}
		path = filepath.Join(home, ".skv.yaml")
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
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
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

// FindByAlias returns the secret with the given alias.
func (c *Config) FindByAlias(alias string) (*Secret, bool) {
	for i := range c.Secrets {
		if c.Secrets[i].Alias == alias {
			return &c.Secrets[i], true
		}
	}
	return nil, false
}

// ToSpec converts the secret into a provider-agnostic spec.
func (s Secret) ToSpec() provider.SecretSpec {
	envName := s.Env
	if strings.TrimSpace(envName) == "" {
		envName = deriveEnvName(s.Alias)
	}
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
	for k, v := range s.Metadata {
		if _, exists := extras[k]; !exists {
			extras[k] = v
		}
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

