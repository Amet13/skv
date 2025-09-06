// Package vault implements HashiCorp Vault provider.
package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
	"skv/internal/provider"
)

type vaultProvider struct{}

// New returns a new Vault provider.
func New() provider.Provider { return &vaultProvider{} }

func (v *vaultProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
	conf := vaultapi.DefaultConfig()
	if addr, ok := spec.Extras["address"]; ok && addr != "" {
		_ = conf.ReadEnvironment() // ignore
		conf.Address = addr
	}
	client, err := vaultapi.NewClient(conf)
	if err != nil {
		return "", fmt.Errorf("vault client: %w", err)
	}
	if ns, ok := spec.Extras["namespace"]; ok && strings.TrimSpace(ns) != "" {
		client.SetNamespace(ns)
	}
	if tok, ok := spec.Extras["token"]; ok && tok != "" {
		client.SetToken(tok)
	} else if roleID, rok := spec.Extras["role_id"]; rok {
		if secretID, sok := spec.Extras["secret_id"]; sok {
			// AppRole login
			secret, err := client.Logical().WriteWithContext(ctx, "auth/approle/login", map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			})
			if err != nil {
				return "", fmt.Errorf("vault approle login: %w", err)
			}
			if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
				return "", errors.New("vault approle login: empty token")
			}
			client.SetToken(secret.Auth.ClientToken)
		}
	}

	// Try KVv2 if we can infer mount and path from name or extras
	if mount, path, ok := kv2MountAndPath(spec); ok {
		sec, err := client.KVv2(mount).Get(ctx, path)
		if err == nil && sec != nil {
			if val, ok := pickValue(sec.Data, spec); ok {
				return val, nil
			}
			b, _ := json.Marshal(sec.Data)
			return string(b), nil
		}
	}

	// Fallback: logical read with raw path (supports non-KV or already fully qualified paths)
	logical := client.Logical()
	s, err2 := logical.ReadWithContext(ctx, spec.Name)
	if err2 != nil {
		return "", fmt.Errorf("vault read: %w", err2)
	}
	if s == nil {
		return "", errors.New("vault: secret not found")
	}
	// KV v2 typically nests data under "data" key
	if nested, ok := s.Data["data"].(map[string]interface{}); ok {
		if val, ok := pickValue(nested, spec); ok {
			return val, nil
		}
		b, _ := json.Marshal(nested)
		return string(b), nil
	}
	if val, ok := pickValue(s.Data, spec); ok {
		return val, nil
	}
	b, _ := json.Marshal(s.Data)
	return string(b), nil
}

func kv2MountAndPath(spec provider.SecretSpec) (string, string, bool) {
	// Explicit mount in extras
	if m, ok := spec.Extras["mount"]; ok && strings.TrimSpace(m) != "" {
		return m, strings.TrimPrefix(spec.Name, "/"), true
	}
	// Infer mount when name looks like "<mount>/data/<path>"
	n := strings.TrimPrefix(spec.Name, "/")
	if strings.Contains(n, "/data/") {
		parts := strings.SplitN(n, "/data/", 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			return parts[0], parts[1], true
		}
	}
	return "", "", false
}

func pickValue(m map[string]interface{}, spec provider.SecretSpec) (string, bool) {
	// Priority: explicit key in extras, then standard "value" key, then single string field.
	if key, ok := spec.Extras["key"]; ok && key != "" {
		if v, ok := m[key]; ok {
			if s, ok := v.(string); ok {
				return s, true
			}
		}
	}
	if v, ok := m["value"]; ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	// If there is exactly one string field, return it.
	var candidate string
	count := 0
	for _, v := range m {
		if s, ok := v.(string); ok {
			candidate = s
			count++
		}
	}
	if count == 1 {
		return candidate, true
	}
	return "", false
}

