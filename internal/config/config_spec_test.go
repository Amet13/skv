package config

import (
	"reflect"
	"testing"
)

func TestToSpecBuildsEnvAndExtras(t *testing.T) {
	v := 7
	s := Secret{
		Alias:    "db-password",
		Provider: "aws-secrets-manager",
		Name:     "app/prod/db",
		Region:   "us-east-1",
		Address:  "https://vault.local",
		Token:    "t",
		Path:     "/p",
		Version:  &v,
		Metadata: map[string]string{"custom": "x", "region": "ignored"},
	}
	spec := s.ToSpec()
	if spec.EnvName != "DB_PASSWORD" {
		t.Fatalf("EnvName = %q want DB_PASSWORD", spec.EnvName)
	}
	if spec.Provider != s.Provider || spec.Name != s.Name || spec.Alias != s.Alias {
		t.Fatalf("provider/name/alias not preserved")
	}
	want := map[string]string{
		"region":  "us-east-1",
		"address": "https://vault.local",
		"token":   "t",
		"path":    "/p",
		"version": "7",
		"custom":  "x",
	}
	if !reflect.DeepEqual(spec.Extras, want) {
		t.Fatalf("extras = %#v want %#v", spec.Extras, want)
	}
}
