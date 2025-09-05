package main

import (
	"bytes"
	"strings"
	"testing"

	"skv/internal/provider"
	mockprovider "skv/internal/provider/mock"
)

// registerMock registers a mock provider under a unique name for e2e tests.
func registerMock(name string) {
	provider.Register(name, mockprovider.New())
}

// writeTempConfig provided by commands_test.go in same package

func TestE2E_Mock_ListGetExportRun(t *testing.T) {
	_ = newRootCmd() // ensure core providers registered
	registerMock("mock")

	cfg := "secrets:\n" +
		"  - alias: a\n    provider: mock\n    name: secretA\n    env: A\n    extras:\n      value: va\n" +
		"  - alias: b\n    provider: mock\n    name: secretB\n    env: B\n    extras:\n      value: vb\n"
	cfgPath = writeTempConfig(t, cfg)

	// list
	var out bytes.Buffer
	c := newListCmd()
	c.SetOut(&out)
	c.SetArgs([]string{"-v"})
	if err := c.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
	s := out.String()
	if !strings.Contains(s, "a") || !strings.Contains(s, "b") {
		t.Fatalf("unexpected list: %q", s)
	}

	// get
	out.Reset()
	g := newGetCmd()
	g.SetOut(&out)
	g.SetArgs([]string{"a"})
	if err := g.Execute(); err != nil {
		t.Fatalf("get: %v", err)
	}
	if strings.TrimSpace(out.String()) != "va" {
		t.Fatalf("unexpected get: %q", out.String())
	}

	// export env format
	out.Reset()
	e := newExportCmd()
	e.SetOut(&out)
	e.SetArgs([]string{"--all", "--format", "env"})
	if err := e.Execute(); err != nil {
		t.Fatalf("export: %v", err)
	}
	exp := strings.TrimSpace(out.String())
	if !strings.Contains(exp, "A=va") || !strings.Contains(exp, "B=vb") {
		t.Fatalf("unexpected export: %q", exp)
	}

	// run dry-run
	var errBuf bytes.Buffer
	r := newRunCmd()
	r.SetErr(&errBuf)
	r.SetArgs([]string{"--all", "--dry-run", "--", "/usr/bin/env"})
	if err := r.Execute(); err != nil {
		t.Fatalf("run: %v", err)
	}
	dr := errBuf.String()
	if !strings.Contains(dr, "A=") || !strings.Contains(dr, "B=") {
		t.Fatalf("unexpected dry-run: %q", dr)
	}
}

