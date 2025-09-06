package main

import (
	"testing"
)

func TestValidateCmd(t *testing.T) {
	tests := getValidationTestCases()
	runTableTests(t, tests, newValidateCmd)
}

