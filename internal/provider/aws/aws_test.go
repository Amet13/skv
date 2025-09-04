package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"skv/internal/provider"
)

type fakeSM struct {
	out *secretsmanager.GetSecretValueOutput
	err error
}

func (f fakeSM) GetSecretValue(_ context.Context, _ *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return f.out, f.err
}

func TestAWSNotFoundMapsToErrNotFound(t *testing.T) {
	oldNew := newSMClient
	defer func() { newSMClient = oldNew }()
	newSMClient = func(_ aws.Config) smClient {
		return fakeSM{nil, &types.ResourceNotFoundException{}}
	}
	a := &awsProvider{}
	_, err := a.FetchSecret(context.Background(), provider.SecretSpec{Name: "n", Extras: map[string]string{"region": "us-east-1"}})
	if !errors.Is(err, provider.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

