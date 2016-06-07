package application_test

import (
	"github.com/nickwei84/sms-resource/out/application"
	"github.com/nickwei84/sms-resource/out/application/applicationfakes"
	"github.com/nickwei84/sms-resource/out/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application", func() {
	var (
		client *applicationfakes.FakeSMSService
		config models.SMSConfig
		app    application.Application
	)

	BeforeSuite(func() {
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

	Describe("Run", func() {
		var runAppErr error

		BeforeEach(func() {
			client = new(applicationfakes.FakeSMSService)
			client.CreateTopicReturns("my-topic-arn", nil)
			client.GetExistingSubscribersReturns([]string{}, nil)
			client.CreateNewSubscriptionsReturns(nil)
			client.PublishMessageReturns(nil)
			app = application.NewApplication(client, config)
		})

		JustBeforeEach(func() {
			runAppErr = app.Run()
		})

		It("should create the SMS topic from configuration", func() {
			Expect(runAppErr).NotTo(HaveOccurred())
			Expect(client.CreateTopicCallCount()).To(Equal(1))
			Expect(client.CreateTopicArgsForCall(0)).To(Equal("my-topic"))
		})

		It("should get existing subscribers of the topic", func() {
			Expect(runAppErr).NotTo(HaveOccurred())
			Expect(client.GetExistingSubscribersCallCount()).To(Equal(1))
			Expect(client.GetExistingSubscribersArgsForCall(0)).To(Equal("my-topic-arn"))
		})

		Context("when there are no existing subscribers to the topic", func() {
			BeforeEach(func() {
				client.GetExistingSubscribersReturns([]string{}, nil)
			})

			It("should subscribe all subscribers from configuration", func() {
				Expect(runAppErr).NotTo(HaveOccurred())
				Expect(client.CreateNewSubscriptionsCallCount()).To(Equal(1))
				arg1, arg2 := client.CreateNewSubscriptionsArgsForCall(0)
				Expect(arg1).To(Equal("my-topic-arn"))
				Expect(arg2).To(Equal([]string{
					"subscriber1",
					"subscriber2",
				}))
			})
		})

		Context("when there are existing subscribers to the topic", func() {
			BeforeEach(func() {
				client.GetExistingSubscribersReturns([]string{
					"subscriber1",
				}, nil)
			})

			It("should subscribe only those subscribers from configuration that are new", func() {
				Expect(runAppErr).NotTo(HaveOccurred())
				Expect(client.CreateNewSubscriptionsCallCount()).To(Equal(1))
				arg1, arg2 := client.CreateNewSubscriptionsArgsForCall(0)
				Expect(arg1).To(Equal("my-topic-arn"))
				Expect(arg2).To(Equal([]string{
					"subscriber2",
				}))
			})
		})

		It("should publish the message from configuration", func() {
			Expect(runAppErr).NotTo(HaveOccurred())
			Expect(client.PublishMessageCallCount()).To(Equal(1))
			arg1, arg2 := client.PublishMessageArgsForCall(0)
			Expect(arg1).To(Equal("my-topic-arn"))
			Expect(arg2).To(Equal("hello"))
		})
	})
})
