# Prometheus SNS Webhook

Send prometheus alerts to an SNS Topic

## Configure

the host machine must have credentials to access aws, you can do this by:
 - adding a default profile to .aws/credentials
 - setting AWS_PROFILE to a profile at .aws/credentials
 - setting AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY

place a prometheus-sns-webhook.yaml file in /etc/prometheus-sns-webhook or the current directory

```
# prometheus-sns-webhook.yaml

sns:
  topicarn: YOUR_ARN
  region: eu-west-1
fakeMessage: false
```

## Send a test alert
```
./tools/send_test.sh
```
