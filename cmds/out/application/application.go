package application

import "github.com/nickwei84/sms-resource/cmds/out/models"

//go:generate counterfeiter . SMSService
type SMSService interface {
	CreateTopic(topic string) (string, error)
	GetExistingSubscribers(topicID string) ([]string, error)
	CreateNewSubscriptions(topicID string, newSubscribers []string) error
	PublishMessage(topicID string, message string) error
}

type Application struct {
	client SMSService
	config models.SMSConfig
}

func NewApplication(client SMSService, config models.SMSConfig) Application {
	return Application{
		client: client,
		config: config,
	}
}

func (a Application) Run() error {
	topicArn, err := a.client.CreateTopic(a.config.Source.Topic)
	if err != nil {
		return err
	}

	existingSubscribers, err := a.client.GetExistingSubscribers(topicArn)
	if err != nil {
		return err
	}

	newSubscribers := findNewSubscribers(existingSubscribers, a.config.Params.Subscribers)

	err = a.client.CreateNewSubscriptions(topicArn, newSubscribers)
	if err != nil {
		return err
	}

	err = a.client.PublishMessage(topicArn, a.config.Params.Message)
	if err != nil {
		return err
	}

	return nil
}

func findNewSubscribers(existingSubscribers []string, subscribersFromInput []string) []string {
	if len(existingSubscribers) == 0 {
		return subscribersFromInput
	}

	newSubscribers := []string{}
	for _, subscriberFromInput := range subscribersFromInput {
		new := true
		for _, existingSubscriber := range existingSubscribers {
			if subscriberFromInput == existingSubscriber {
				new = false
				break
			}
		}
		if new {
			newSubscribers = append(newSubscribers, subscriberFromInput)
		}
	}

	return newSubscribers
}
