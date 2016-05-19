package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type MetadataItem struct {
	Name  string
	Value string
}

type OutputJSON struct {
	Version struct {
		Time time.Time
	}
	Metadata []MetadataItem
}

type SMSConfig struct {
	Source struct {
		AWSAccessKeyID     string `json:"aws_access_key_id"`
		AWSSecretAccessKey string `json:"aws_secret_access_key"`
		Topic              string `json:"topic"`
	} `json:"source"`
	Params struct {
		Subscribers []string `json:"subscribers"`
		Message     string   `json:"message"`
	} `json:"params"`
}

func (s *SMSConfig) checkInput() error {
	if s.Source.AWSAccessKeyID == "" {
		return fmt.Errorf("Error: source.aws_access_key_id from stdin is either empty or missing")
	}

	if s.Source.AWSSecretAccessKey == "" {
		return fmt.Errorf("Error: source.aws_secret_access_key from stdin is either empty or missing")
	}

	if s.Source.Topic == "" {
		return fmt.Errorf("Error: source.topic from stdin is either empty or missing")
	}

	if len(s.Source.Topic) > 10 {
		return fmt.Errorf("Error: source.topic from stdin cannot exceed 10 characters")
	}

	if len(s.Params.Subscribers) == 0 {
		return fmt.Errorf("Error: params.subscribers from stdin is either empty or missing")
	}

	if s.Params.Message == "" {
		return fmt.Errorf("Error: params.message from stdin is either empty or missing")
	}

	return nil
}

func exitWithErr(err interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func getStdinInput(config *SMSConfig) error {
	stdinData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("Error reading from stdin: %v", err)
	}

	err = json.Unmarshal(stdinData, config)
	if err != nil {
		return fmt.Errorf("Error parsing stdin as JSON: %v", err)
	}

	return nil
}

func createTopic(snsService *sns.SNS, topic string) (string, error) {
	createTopicResp, err := snsService.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(topic),
	})
	if err != nil {
		return "", fmt.Errorf("Error creating topic: %v", err)
	}

	topicArn := *createTopicResp.TopicArn

	_, err = snsService.SetTopicAttributes(&sns.SetTopicAttributesInput{
		TopicArn:       aws.String(topicArn),
		AttributeName:  aws.String("DisplayName"),
		AttributeValue: aws.String(topic),
	})
	if err != nil {
		return "", fmt.Errorf("Error creating SMS display name for topic: %v", err)
	}

	return topicArn, nil
}

func getSubscribers(snsService *sns.SNS, topicArn string) ([]string, error) {
	existingSubscribers := []string{}

	// TODO: this returns a max of 100 subscriptions, additional calls will be needed to retrieve paginated results
	listSubscriptionsByTopicResp, err := snsService.ListSubscriptionsByTopic(&sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		return existingSubscribers, fmt.Errorf("Error getting list of existing subscribers: %v", err)
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

func findNewSubscribers(existingSubscribers []string, subscribersFromInput []string) []string {
	newSubscribers := []string{}

	if len(existingSubscribers) == 0 {
		newSubscribers = subscribersFromInput
	} else {
		for _, existingSubscription := range existingSubscribers {
			for _, subscriber := range subscribersFromInput {
				if subscriber == existingSubscription {
					break
				}
				newSubscribers = append(newSubscribers, subscriber)
			}
		}
	}

	return newSubscribers
}

func createNewSubscriptions(snsService *sns.SNS, topicArn string, newSubscribers []string) error {
	for _, subscriber := range newSubscribers {
		_, err := snsService.Subscribe(&sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("sms"),
			Endpoint: aws.String(subscriber),
		})
		if err != nil {
			return fmt.Errorf("Error subscribing %s: %v", subscriber, err)
		}
	}

	return nil
}

func publishMessage(snsService *sns.SNS, topicArn string, message string) error {
	_, err := snsService.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("Error publishing message: %v", err)
	}

	return nil
}

func generateStdoutOutput() ([]byte, error) {
	var output OutputJSON
	output.Version.Time = time.Now().UTC()
	output.Metadata = []MetadataItem{}

	stdoutOutput, err := json.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling output for stdout: %v", err)
	}

	return stdoutOutput, nil
}

func main() {
	var config SMSConfig

	err := getStdinInput(&config)
	if err != nil {
		exitWithErr(err)
	}

	err = config.checkInput()
	if err != nil {
		exitWithErr(err)
	}

	creds := credentials.NewStaticCredentials(config.Source.AWSAccessKeyID, config.Source.AWSSecretAccessKey, "")
	snsService := sns.New(session.New(), aws.NewConfig().WithCredentials(creds).WithRegion("us-east-1"))

	topicArn, err := createTopic(snsService, config.Source.Topic)
	if err != nil {
		exitWithErr(err)
	}

	existingSubscribers, err := getSubscribers(snsService, topicArn)
	if err != nil {
		exitWithErr(err)
	}

	newSubscribers := findNewSubscribers(existingSubscribers, config.Params.Subscribers)

	err = createNewSubscriptions(snsService, topicArn, newSubscribers)
	if err != nil {
		exitWithErr(err)
	}

	err = publishMessage(snsService, topicArn, config.Params.Message)
	if err != nil {
		exitWithErr(err)
	}

	stdoutOutput, err := generateStdoutOutput()
	if err != nil {
		exitWithErr(err)
	}

	fmt.Println(string(stdoutOutput))
}
