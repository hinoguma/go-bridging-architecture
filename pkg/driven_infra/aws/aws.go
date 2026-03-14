package driven_infra_aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// how to create aws service client with config
// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/getting-started.html
func NewAWSConfig(ctx context.Context) (aws.Config, error) {
	// load aws config from environment variables
	return config.LoadDefaultConfig(ctx)
}
