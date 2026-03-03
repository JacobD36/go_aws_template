package infrastructure

import (
	"context"
	"encoding/json"
	"log"
	"messaging-service/internal/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSEventPublisher implementa el publicador de eventos usando SQS
type SQSEventPublisher struct {
	client   *sqs.Client
	queueURL string
}

// NewSQSEventPublisher crea una nueva instancia del publicador
func NewSQSEventPublisher(client *sqs.Client, queueURL string) *SQSEventPublisher {
	return &SQSEventPublisher{
		client:   client,
		queueURL: queueURL,
	}
}

// PublishLogEvent publica un evento de log
func (p *SQSEventPublisher) PublishLogEvent(ctx context.Context, event *domain.LogEvent) error {
	messageBody, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(messageBody)),
	})

	if err != nil {
		log.Printf("Error publishing log event to SQS: %v", err)
		return err
	}

	log.Printf("Log event published successfully: %s", event.EventType)
	return nil
}
