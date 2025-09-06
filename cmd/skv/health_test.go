package main

import (
	"testing"
)

func TestHealthCmd(t *testing.T) {
	skipIfShort(t)

	tests := getHealthTestCases()
	runTableTests(t, tests, newHealthCmd)
}
