package infrastructure

import (
	"auth-service/internal/domain"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBUserRepository implementa el repositorio de usuarios usando DynamoDB
type DynamoDBUserRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBUserRepository crea una nueva instancia del repositorio
func NewDynamoDBUserRepository(client *dynamodb.Client, tableName string) *DynamoDBUserRepository {
	return &DynamoDBUserRepository{
		client:    client,
		tableName: tableName,
	}
}

// FindByEmail busca un usuario por su email
func (r *DynamoDBUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	// DynamoDB Scan para buscar por email (en producción usar GSI)
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("Email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	})

	if err != nil {
		log.Printf("Error scanning DynamoDB for email %s: %v", email, err)
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, domain.ErrUserNotFound
	}

	var user domain.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		log.Printf("Error unmarshaling user: %v", err)
		return nil, err
	}

	return &user, nil
}
