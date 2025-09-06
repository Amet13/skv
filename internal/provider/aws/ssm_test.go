package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
	"skv/internal/provider"
)

type fakeSSMClient struct {
	out *ssm.GetParameterOutput
	err error
}

func (f *fakeSSMClient) GetParameter(_ context.Context, _ *ssm.GetParameterInput, _ ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	return f.out, f.err
}

func TestSSMNotFoundMapsToErrNotFound(t *testing.T) {
	oldNew := newSSMClient
	defer func() { newSSMClient = oldNew }()
	newSSMClient = func(_ aws.Config) ssmClient {
		return &fakeSSMClient{err: &smithy.GenericAPIError{Code: "ParameterNotFound"}}
	}
	p := NewSSM()
	_, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "/path/name"})
	if !errors.Is(err, provider.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSSMSuccess(t *testing.T) {
	oldNew := newSSMClient
	defer func() { newSSMClient = oldNew }()
	val := "ok"
	newSSMClient = func(_ aws.Config) ssmClient {
		return &fakeSSMClient{out: &ssm.GetParameterOutput{Parameter: &types.Parameter{Value: &val}}}
	}
	p := NewSSM()
	out, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "/path/name"})
	if err != nil || out != "ok" {
		t.Fatalf("got %q err=%v", out, err)
	}
}

