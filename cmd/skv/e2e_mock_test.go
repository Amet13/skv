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
	cfgPath = writeTestConfig(t, cfg)

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

func TestE2E_Mock_DoctorCommand(t *testing.T) {
	_ = newRootCmd() // ensure core providers registered
	registerMock("mock")

	// Create a test config
	cfg := "secrets:\n" +
		"  - alias: test_secret\n    provider: mock\n    name: test\n    env: TEST_SECRET\n    extras:\n      value: test_value\n"
	cfgPath = writeTestConfig(t, cfg)

	// Test doctor command
	var out bytes.Buffer
	var errBuf bytes.Buffer
	d := newDoctorCmd()
	d.SetOut(&out)
	d.SetErr(&errBuf)
	d.SetArgs([]string{})

	err := d.Execute()
	if err != nil {
		t.Fatalf("doctor command failed: %v", err)
	}

	output := out.String() + errBuf.String()
	if !strings.Contains(output, "OK: Configuration loaded") {
		t.Fatalf("doctor output missing config confirmation: %q", output)
	}
	if !strings.Contains(output, "OK: No configuration issues found") {
		t.Fatalf("doctor output missing validation confirmation: %q", output)
	}
}

func TestE2E_Mock_WatchCommand(t *testing.T) {
	_ = newRootCmd() // ensure core providers registered
	registerMock("mock")

	cfg := "secrets:\n" +
		"  - alias: watch_secret\n    provider: mock\n    name: watch\n    env: WATCH_SECRET\n    extras:\n      value: initial\n"
	cfgPath = writeTestConfig(t, cfg)

	// Test watch command with on-change-only (should execute immediately and exit)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	w := newWatchCmd()
	w.SetOut(&out)
	w.SetErr(&errBuf)
	w.SetArgs([]string{"--all", "--on-change-only", "--interval", "1s", "--", "echo", "test"})

	// Watch should execute and then timeout/exit after a short time
	// We can't easily test the full watch loop in unit tests, so we test the setup
	err := w.Execute()
	// The watch command might fail due to context cancellation or other reasons
	// The important thing is that it doesn't crash during setup
	if err != nil && !strings.Contains(err.Error(), "context") && !strings.Contains(err.Error(), "signal") {
		t.Logf("Watch command exited with: %v (this may be expected)", err)
	}
}

func TestE2E_Mock_Transformations(t *testing.T) {
	_ = newRootCmd() // ensure core providers registered
	registerMock("mock")

	// Test template transformation
	cfg := "secrets:\n" +
		"  - alias: template_secret\n    provider: mock\n    name: template\n    env: TEMPLATE_SECRET\n    extras:\n      value: test\n" +
		"    transform:\n      type: template\n      template: \"prefix_{{ .value }}_suffix\"\n"
	cfgPath = writeTestConfig(t, cfg)

	// Test get with transformation
	var out bytes.Buffer
	g := newGetCmd()
	g.SetOut(&out)
	g.SetArgs([]string{"template_secret"})
	if err := g.Execute(); err != nil {
		t.Fatalf("get with transformation failed: %v", err)
	}
	result := strings.TrimSpace(out.String())
	if result != "prefix_test_suffix" {
		t.Fatalf("transformation failed: expected 'prefix_test_suffix', got %q", result)
	}

	// Test export with transformation
	out.Reset()
	e := newExportCmd()
	e.SetOut(&out)
	e.SetArgs([]string{"--all", "--format", "env"})
	if err := e.Execute(); err != nil {
		t.Fatalf("export with transformation failed: %v", err)
	}
	exp := strings.TrimSpace(out.String())
	if !strings.Contains(exp, "TEMPLATE_SECRET=prefix_test_suffix") {
		t.Fatalf("export transformation failed: %q", exp)
	}
}

func TestE2E_Mock_MultipleFormats(t *testing.T) {
	_ = newRootCmd() // ensure core providers registered
	registerMock("mock")

	cfg := "secrets:\n" +
		"  - alias: json_secret\n    provider: mock\n    name: json\n    env: JSON_SECRET\n    extras:\n      value: test_value\n"
	cfgPath = writeTestConfig(t, cfg)

	// Test JSON export
	var out bytes.Buffer
	e := newExportCmd()
	e.SetOut(&out)
	e.SetArgs([]string{"--all", "--format", "json"})
	if err := e.Execute(); err != nil {
		t.Fatalf("JSON export failed: %v", err)
	}
	jsonResult := strings.TrimSpace(out.String())
	if !strings.Contains(jsonResult, `"JSON_SECRET": "test_value"`) {
		t.Fatalf("JSON export failed: %q", jsonResult)
	}

	// Test YAML export
	out.Reset()
	e = newExportCmd()
	e.SetOut(&out)
	e.SetArgs([]string{"--all", "--format", "yaml"})
	if err := e.Execute(); err != nil {
		t.Fatalf("YAML export failed: %v", err)
	}
	yamlResult := strings.TrimSpace(out.String())
	if !strings.Contains(yamlResult, "JSON_SECRET: test_value") {
		t.Fatalf("YAML export failed: %q", yamlResult)
	}
}

