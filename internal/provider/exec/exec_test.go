package execprovider

import (
	"context"
	"os/exec"
	"runtime"
	"testing"

	"skv/internal/provider"
)

func TestExecProviderSuccess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	p := New()
	out, err := p.FetchSecret(context.Background(), provider.SecretSpec{
		Alias:    "a",
		Name:     "world",
		Provider: "exec",
		Extras: map[string]string{
			"cmd":  "echo",
			"args": "hello",
			"trim": "true",
		},
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if out != "hello world" {
		t.Fatalf("got %q", out)
	}
}

func TestExecProviderError(t *testing.T) {
	p := New()
	_, err := p.FetchSecret(context.Background(), provider.SecretSpec{
		Alias:    "a",
		Name:     "arg",
		Provider: "exec",
		Extras: map[string]string{
			"cmd": "does-not-exist-hopefully",
		},
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	var ee *exec.Error
	_ = ee // acceptable: exact error type varies by platform
}

func TestExecArgsAndTrim(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	p := New()
	out, err := p.FetchSecret(context.Background(), provider.SecretSpec{
		Alias: "a",
		Name:  "/bin/echo",
		Extras: map[string]string{
			"args": "hello\n",
			"trim": "true",
		},
	})
	if err != nil || out != "hello" {
		t.Fatalf("got %q err=%v", out, err)
	}
}

