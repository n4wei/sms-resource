package models_test

import (
	"github.com/nickwei84/sms-resource/cmds/out/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SMSConfig", func() {
	Describe("CheckInput", func() {
		var config models.SMSConfig

		BeforeEach(func() {
			config = models.SMSConfig{
				Source: models.Source{
					AWSAccessKeyID:     "key123",
					AWSSecretAccessKey: "secretabc",
					Topic:              "my-topic",
				},
				Params: models.Params{
					Subscribers: []string{
						"subscriber1",
						"subscriber2",
					},
					Message: "hello",
				},
			}
		})

		It("should return an error if AWS key ID is missing", func() {
			config.Source.AWSAccessKeyID = ""
			err := config.CheckInput()
			Expect(err).Should(MatchError("source.aws_access_key_id from stdin is either empty or missing"))
		})

		It("should return an error if AWS secret key is missing", func() {
			config.Source.AWSSecretAccessKey = ""
			err := config.CheckInput()
			Expect(err).Should(MatchError("source.aws_secret_access_key from stdin is either empty or missing"))
		})

		It("should return an error if topic is missing", func() {
			config.Source.Topic = ""
			err := config.CheckInput()
			Expect(err).Should(MatchError("source.topic from stdin is either empty or missing"))
		})

		It("should return an error if topic exceeds max character limit", func() {
			config.Source.Topic = "very-long-topic-1234567890abcdefg"
			err := config.CheckInput()
			Expect(err).Should(MatchError("source.topic from stdin cannot exceed 10 characters"))
		})

		It("should return an error if no subscribers are provided", func() {
			config.Params.Subscribers = []string{}
			err := config.CheckInput()
			Expect(err).Should(MatchError("params.subscribers from stdin is either empty or missing"))
		})

		It("should return an error if message is missing", func() {
			config.Params.Message = ""
			err := config.CheckInput()
			Expect(err).Should(MatchError("params.message from stdin is either empty or missing"))
		})

		It("should not return an error if all fields are valid", func() {
			err := config.CheckInput()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
