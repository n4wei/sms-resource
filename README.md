# SMS Concourse Resource

Sends SMS messages using AWS SNS (Simple Notification Service).

## Source Configuration

- `aws_access_key_id`: *Required.* The AWS credential for accessing the SNS service.
- `aws_secret_access_key`: *Required.* The AWS credential for accessing the SNS service.
- `topic`: *Required.* The topic of the SMS messages. Phone numbers are subscribed to the topic and messages are published to the topic.

### Example

The SMS resource is available on Dockerhub at [`nwei/sms-concourse-resource`](https://hub.docker.com/r/nwei/sms-concourse-resource/)

To use this in your concourse pipeline:

```yaml
resource_types:
- name: sms-resource
  type: docker-image
  source:
    repository: nwei/sms-concourse-resource

resources:
- name: sms
  type: sms-resource
  source:
    aws_access_key_id: abc123
    aws_secret_access_key: secret
    topic: concourse
```

```yaml
- put: sms
  params:
    subscribers: ["14151234567", "16501234567"]
    message: "hello"
```

## Behavior

### `check`: No-Op

### `in`: No-Op

### `out`: Send SMS message

Subscribe phone number(s) provided in `subscribers` to a topic, and publishes the message provided in `message` to that topic.

#### Parameters

- `subscribers`: *Required.* A list of phone numbers to subscribe to the topic.
- `message`: *Required.* The message to publish to the topic.
