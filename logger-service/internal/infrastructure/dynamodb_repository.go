package infrastructure

import (
	"context"
	"log"
	"logger-service/internal/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBLogRepository implementa el repositorio de logs usando DynamoDB
type DynamoDBLogRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBLogRepository crea una nueva instancia del repositorio
func NewDynamoDBLogRepository(client *dynamodb.Client, tableName string) *DynamoDBLogRepository {
	return &DynamoDBLogRepository{
		client:    client,
		tableName: tableName,
	}
}

// Save guarda una entrada de log en DynamoDB
func (r *DynamoDBLogRepository) Save(ctx context.Context, entry *domain.LogEntry) error {
	item, err := attributevalue.MarshalMap(entry)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		log.Printf("Error saving log entry to DynamoDB: %v", err)
		return err
	}

	log.Printf("Log entry saved successfully: %s", entry.ID)
	return nil
}

// FindAll obtiene todas las entradas de log
func (r *DynamoDBLogRepository) FindAll(ctx context.Context) ([]*domain.LogEntry, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})

	if err != nil {
		return nil, err
	}

	var entries []*domain.LogEntry
	for _, item := range result.Items {
		var entry domain.LogEntry
		err := attributevalue.UnmarshalMap(item, &entry)
		if err != nil {
			log.Printf("Error unmarshaling log entry: %v", err)
			continue
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}
