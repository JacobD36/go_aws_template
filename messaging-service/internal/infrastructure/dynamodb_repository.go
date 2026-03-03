package infrastructure

import (
	"context"
	"log"
	"messaging-service/internal/domain"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBRepository implementa el repositorio usando DynamoDB
type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBRepository crea una nueva instancia del repositorio
func NewDynamoDBRepository(client *dynamodb.Client, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}
}

// Save guarda un mensaje en DynamoDB
func (r *DynamoDBRepository) Save(ctx context.Context, message *domain.Message) error {
	// Actualizar timestamp de envío
	now := time.Now()
	message.SentAt = &now
	message.Status = "sent"

	item, err := attributevalue.MarshalMap(message)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		log.Printf("Error saving message to DynamoDB: %v", err)
		return err
	}

	log.Printf("Message saved to DynamoDB: %s", message.ID)
	return nil
}

// FindByID busca un mensaje por ID en DynamoDB
func (r *DynamoDBRepository) FindByID(ctx context.Context, id string) (*domain.Message, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	})

	if err != nil {
		log.Printf("Error getting message from DynamoDB: %v", err)
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var message domain.Message
	err = attributevalue.UnmarshalMap(result.Item, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
