package localstack_sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/mitchellh/mapstructure"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/localstack-sqs", new(Sqs))
}

type Sqs struct{}

func (*Sqs) NewClient(endpoint string, region string) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Failed to get AWS default configs: " + err.Error())
	}

	cfg.Region = region
	cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})
	cfg.Credentials = aws.AnonymousCredentials{}

	client := sqs.NewFromConfig(cfg)

	return client
}

func (s *Sqs) SendMessage(sqsClient *sqs.Client, input any) {
	var sqsMessageInput sqs.SendMessageInput
	_ = mapstructure.Decode(input, &sqsMessageInput)
	_, err := sqsClient.SendMessage(context.TODO(), &sqsMessageInput)
	if err != nil {
		panic("Failed when sending message: " + err.Error())
	}
}
