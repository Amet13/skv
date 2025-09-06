package vault

import (
	"testing"

	"skv/internal/provider"
)

func TestNew(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("New() returned nil")
	}

	// Ensure it implements the Provider interface
	_ = p
}

func TestKv2MountAndPath(t *testing.T) {
	m, p, ok := kv2MountAndPath(provider.SecretSpec{Name: "kv/data/foo/bar"})
	if !ok || m != "kv" || p != "foo/bar" {
		t.Fatalf("got %v %v ok=%v", m, p, ok)
	}
	m, p, ok = kv2MountAndPath(provider.SecretSpec{Name: "/kv/data/x"})
	if !ok || m != "kv" || p != "x" {
		t.Fatalf("got %v %v ok=%v", m, p, ok)
	}
	_, _, ok = kv2MountAndPath(provider.SecretSpec{Name: "n/a"})
	if ok {
		t.Fatalf("expected not ok")
	}
}

func TestPickValue(t *testing.T) {
	// explicit key
	v, ok := pickValue(map[string]interface{}{"a": "1", "b": "2"}, provider.SecretSpec{Extras: map[string]string{"key": "b"}})
	if !ok || v != "2" {
		t.Fatalf("want 2 got %q ok=%v", v, ok)
	}
	// default 'value'
	v, ok = pickValue(map[string]interface{}{"value": "x"}, provider.SecretSpec{})
	if !ok || v != "x" {
		t.Fatalf("want x got %q ok=%v", v, ok)
	}
	// single string field
	v, ok = pickValue(map[string]interface{}{"only": "x", "n": 1}, provider.SecretSpec{})
	if !ok || v != "x" {
		t.Fatalf("want x got %q ok=%v", v, ok)
	}
	// none
	if _, ok := pickValue(map[string]interface{}{"n": 1}, provider.SecretSpec{}); ok {
		t.Fatalf("expected not ok")
	}
}

