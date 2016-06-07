package models

import (
	"fmt"
	"time"
)

type OutputJSON struct {
	Version  Version
	Metadata []MetadataItem
}

type Version struct {
	Time time.Time
}

type MetadataItem struct {
	Name  string
	Value string
}

type SMSConfig struct {
	Source Source `json:"source"`
	Params Params `json:"params"`
}

type Source struct {
	AWSAccessKeyID     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
	Topic              string `json:"topic"`
}

type Params struct {
	Subscribers []string `json:"subscribers"`
	Message     string   `json:"message"`
}

func (s SMSConfig) CheckInput() error {
	if s.Source.AWSAccessKeyID == "" {
		return fmt.Errorf("source.aws_access_key_id from stdin is either empty or missing")
	}

	if s.Source.AWSSecretAccessKey == "" {
		return fmt.Errorf("source.aws_secret_access_key from stdin is either empty or missing")
	}

	if s.Source.Topic == "" {
		return fmt.Errorf("source.topic from stdin is either empty or missing")
	}

	if len(s.Source.Topic) > 10 {
		return fmt.Errorf("source.topic from stdin cannot exceed 10 characters")
	}

	if len(s.Params.Subscribers) == 0 {
		return fmt.Errorf("params.subscribers from stdin is either empty or missing")
	}

	if s.Params.Message == "" {
		return fmt.Errorf("params.message from stdin is either empty or missing")
	}

	return nil
}
