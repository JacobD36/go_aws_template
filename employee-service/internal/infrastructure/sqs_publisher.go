package infrastructure

import (
	"context"
	"employee-service/internal/domain"
	"encoding/json"
	"log"

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

// PublishEmployeeCreated publica un evento de empleado creado
func (p *SQSEventPublisher) PublishEmployeeCreated(ctx context.Context, event *domain.EmployeeEvent) error {
	messageBody, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(messageBody)),
	})

	if err != nil {
		log.Printf("Error publishing event to SQS: %v", err)
		return err
	}

	log.Printf("Event published successfully: %s", event.EventType)
	return nil
}
