package awsclient

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type AWSClient struct {
	snsService *sns.SNS
}

func NewAWSClient(awsAccessKeyID string, awsSecretAccessKey string) AWSClient {
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")
	return AWSClient{
		snsService: sns.New(session.New(), aws.NewConfig().WithCredentials(creds).WithRegion("us-east-1")),
	}
}

func (s AWSClient) CreateTopic(topic string) (string, error) {
	createTopicResp, err := s.snsService.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(topic),
	})
	if err != nil {
		return "", fmt.Errorf("error creating topic: %v", err)
	}

	topicArn := *createTopicResp.TopicArn

	_, err = s.snsService.SetTopicAttributes(&sns.SetTopicAttributesInput{
		TopicArn:       aws.String(topicArn),
		AttributeName:  aws.String("DisplayName"),
		AttributeValue: aws.String(topic),
	})
	if err != nil {
		return "", fmt.Errorf("error creating SMS display name for topic: %v", err)
	}

	return topicArn, nil
}

func (s AWSClient) GetExistingSubscribers(topicArn string) ([]string, error) {
	existingSubscribers := []string{}

	// TODO: this returns a max of 100 subscriptions, additional calls will be needed to retrieve paginated results
	listSubscriptionsByTopicResp, err := s.snsService.ListSubscriptionsByTopic(&sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		return existingSubscribers, fmt.Errorf("error getting list of existing subscribers: %v", err)
	}

	subscriptions := listSubscriptionsByTopicResp.Subscriptions

	if subscriptions == nil {
		return existingSubscribers, nil
	}

	for _, subscription := range subscriptions {
		existingSubscribers = append(existingSubscribers, *subscription.Endpoint)
	}

	return existingSubscribers, nil
}

func (s AWSClient) CreateNewSubscriptions(topicArn string, newSubscribers []string) error {
	for _, subscriber := range newSubscribers {
		_, err := s.snsService.Subscribe(&sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("sms"),
			Endpoint: aws.String(subscriber),
		})
		if err != nil {
			return fmt.Errorf("error subscribing %s: %v", subscriber, err)
		}
	}

	return nil
}

func (s AWSClient) PublishMessage(topicArn string, message string) error {
	_, err := s.snsService.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("error publishing message: %v", err)
	}

	return nil
}
