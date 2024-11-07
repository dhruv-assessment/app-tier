package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
)

func GetInstanceID() (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", nil
	}
	client := imds.NewFromConfig(cfg)
	output, err := client.GetMetadata(context.TODO(), &imds.GetMetadataInput{
		Path: "instance-id",
	})
	if err != nil {
		return "", err
	}
	defer output.Content.Close()
	bytes, err := io.ReadAll(output.Content)
	if err != nil {
		return "", err
	}
	resp := string(bytes)
	print(resp)
	return resp, err
}

func TerminateEC2(instanceID string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := ec2.NewFromConfig(cfg)
	_, err = client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	})
	return err
}

func DownloadFromS3(fileName string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)
	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_IN_BUCKET_NAME")),
		Key:    aws.String(fileName),
	})
	if err != nil {
		log.Printf("Couldn't get the object from S3: %v", err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v: %v\n", fileName, err)
	}
	_, err = file.Write(body)
	return err
}

func ReadMessagesInRequestQueue() (string, string, string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", "", "", err
	}
	client := sqs.NewFromConfig(cfg)
	log.Print("Reading message from queue")
	result, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(os.Getenv("AWS_REQ_URL")),
		VisibilityTimeout:   45,
		WaitTimeSeconds:     0,
		MaxNumberOfMessages: 1,
	})
	if err != nil {
		return "", "", "", err
	}
	if len(result.Messages) > 0 {
		log.Printf("Message ID: %v ---> Body: %v\n", result.Messages[0].MessageId, *result.Messages[0].Body)
		return *result.Messages[0].Body, *result.Messages[0].MessageId, *result.Messages[0].ReceiptHandle, nil
	}
	return "", "", "", fmt.Errorf("no meesage found")
}

func SendMessageToSQS(prediction string, reqQueueMessageID string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}
	client := sqs.NewFromConfig(cfg)

	if sender, err := client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(prediction),
		QueueUrl:    aws.String(os.Getenv("AWS_RESP_URL")),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Request-Queue-Message-ID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(reqQueueMessageID),
			},
		},
		DelaySeconds: 0,
	}); err != nil {
		return "", err
	} else {
		fmt.Println(*sender.MessageId)
		return *sender.MD5OfMessageBody, nil
	}
}

func DeleteMessageFromSQS(reqMessageReceiptID string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := sqs.NewFromConfig(cfg)

	_, err = client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(os.Getenv("AWS_REQ_URL")),
		ReceiptHandle: aws.String(reqMessageReceiptID),
	})
	if err != nil {
		return err
	}
	return nil
}

func GetNoOfMessagesInRequestQueue() (int, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return -1, err
	}
	client := sqs.NewFromConfig(cfg)

	result, err := client.GetQueueAttributes(context.TODO(), &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(os.Getenv("AWS_REQ_URL")),
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameApproximateNumberOfMessages,
		},
	})
	if err != nil {
		return -1, err
	}
	noOfMessages, err := strconv.Atoi(result.Attributes[string(types.QueueAttributeNameApproximateNumberOfMessages)])
	return noOfMessages, err
}
