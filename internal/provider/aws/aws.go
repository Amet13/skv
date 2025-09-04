// Package aws implements AWS Secrets Manager provider.
package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"skv/internal/provider"
)

type awsProvider struct{}

// New returns a new AWS Secrets Manager provider.
func New() provider.Provider { return &awsProvider{} }

// seam interfaces/funcs for testing
type smClient interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

var loadAWSConfig = awsconfig.LoadDefaultConfig
var newSMClient = func(cfg aws.Config) smClient { return secretsmanager.NewFromConfig(cfg) }

func (a *awsProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
	// Region precedence: spec.Extras["region"] > default chain
	var opts []func(*awsconfig.LoadOptions) error
	if r, ok := spec.Extras["region"]; ok && r != "" {
		opts = append(opts, awsconfig.WithRegion(r))
	}
	cfg, err := loadAWSConfig(ctx, opts...)
	if err != nil {
		return "", fmt.Errorf("aws config: %w", err)
	}

	sm := newSMClient(cfg)
	var versionStage *string
	if s, ok := spec.Extras["version_stage"]; ok && s != "" {
		versionStage = &s
	}
	var versionID *string
	if vid, ok := spec.Extras["version_id"]; ok && vid != "" {
		versionID = &vid
	}
	out, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     &spec.Name,
		VersionId:    versionID,
		VersionStage: versionStage,
	})
	if err != nil {
		var rnfe *types.ResourceNotFoundException
		if errors.As(err, &rnfe) {
			return "", provider.ErrNotFound
		}
		return "", fmt.Errorf("aws get secret: %w", err)
	}
	if out.SecretString != nil {
		return *out.SecretString, nil
	}
	if out.SecretBinary != nil {
		return string(out.SecretBinary), nil
	}
	return "", fmt.Errorf("aws secret has no SecretString or SecretBinary: %s", spec.Name)
}

