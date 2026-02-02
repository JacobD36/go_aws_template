package infrastructure

import (
	"context"
	"encoding/json"
	"log"
	"logger-service/internal/domain"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSEventConsumer implementa el consumidor de eventos usando SQS
type SQSEventConsumer struct {
	client   *sqs.Client
	queueURL string
}

// NewSQSEventConsumer crea una nueva instancia del consumidor
func NewSQSEventConsumer(client *sqs.Client, queueURL string) *SQSEventConsumer {
	return &SQSEventConsumer{
		client:   client,
		queueURL: queueURL,
	}
}

// ConsumeEvents consume eventos de SQS
func (c *SQSEventConsumer) ConsumeEvents(ctx context.Context, handler func(*domain.EmployeeEvent) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			messages, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(c.queueURL),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
				VisibilityTimeout:   30,
			})

			if err != nil {
				log.Printf("Error receiving messages from SQS: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, message := range messages.Messages {
				if err := c.processMessage(ctx, message, handler); err != nil {
					log.Printf("Error processing message: %v", err)
					continue
				}

				// Eliminar mensaje de la cola
				_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(c.queueURL),
					ReceiptHandle: message.ReceiptHandle,
				})

				if err != nil {
					log.Printf("Error deleting message from SQS: %v", err)
				}
			}
		}
	}
}

func (c *SQSEventConsumer) processMessage(ctx context.Context, message types.Message, handler func(*domain.EmployeeEvent) error) error {
	var event domain.EmployeeEvent
	if err := json.Unmarshal([]byte(*message.Body), &event); err != nil {
		return err
	}

	return handler(&event)
}
