// Package aws provides AWS-related providers. This file implements AWS SSM Parameter Store.
package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/smithy-go"
	"skv/internal/provider"
)

type ssmProvider struct{}

// NewSSM returns a new AWS SSM Parameter Store provider.
func NewSSM() provider.Provider { return &ssmProvider{} }

// seam interfaces/funcs for testing
type ssmClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

var loadAWSConfigSSM = awsconfig.LoadDefaultConfig
var newSSMClient = func(cfg aws.Config) ssmClient { return ssm.NewFromConfig(cfg) }

func (p *ssmProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
	var opts []func(*awsconfig.LoadOptions) error
	if prof := strings.TrimSpace(spec.Extras["profile"]); prof != "" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(prof))
	}
	if r := strings.TrimSpace(spec.Extras["region"]); r != "" {
		opts = append(opts, awsconfig.WithRegion(r))
	}
	cfg, err := loadAWSConfigSSM(ctx, opts...)
	if err != nil {
		return "", fmt.Errorf("aws ssm config: %w", err)
	}

	client := newSSMClient(cfg)
	withDecryption := true
	if wd := strings.TrimSpace(spec.Extras["with_decryption"]); strings.EqualFold(wd, "false") || wd == "0" {
		withDecryption = false
	}
	out, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(spec.Name),
		WithDecryption: aws.Bool(withDecryption),
	})
	if err != nil {
		// Map not found errors using smithy API error code
		if apiErr, ok := err.(smithy.APIError); ok {
			if code := apiErr.ErrorCode(); code == "ParameterNotFound" || code == "ParameterVersionNotFound" {
				return "", provider.ErrNotFound
			}
		}
		return "", fmt.Errorf("aws ssm get parameter: %w", err)
	}
	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", fmt.Errorf("aws ssm: empty value for %s", spec.Name)
	}
	return aws.ToString(out.Parameter.Value), nil
}
