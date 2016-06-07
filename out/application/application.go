package application

import "github.com/nickwei84/sms-resource/out/models"

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

	return a.client.PublishMessage(topicArn, a.config.Params.Message)
}

func findNewSubscribers(existingSubscribers []string, subscribersFromInput []string) []string {
	if len(existingSubscribers) == 0 {
		return subscribersFromInput
	}

	existingSubscribersMap := map[string]string{}
	for _, existingSubscriber := range existingSubscribers {
		existingSubscribersMap[existingSubscriber] = ""
	}

	newSubscribers := []string{}
	for _, subscriberFromInput := range subscribersFromInput {
		_, exist := existingSubscribersMap[subscriberFromInput]
		if !exist {
			newSubscribers = append(newSubscribers, subscriberFromInput)
		}
	}

	return newSubscribers
}
