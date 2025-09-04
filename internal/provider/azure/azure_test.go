package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestAzureResponseErrorType(t *testing.T) {
	// Simple type sanity check to avoid unused imports while keeping mapping logic covered elsewhere
	var err error = &azcore.ResponseError{StatusCode: 403}
	if _, ok := err.(*azcore.ResponseError); !ok {
		t.Fatalf("expected azcore.ResponseError")
	}
}

