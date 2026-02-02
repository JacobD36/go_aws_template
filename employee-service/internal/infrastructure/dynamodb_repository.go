package infrastructure

import (
	"context"
	"employee-service/internal/domain"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

// Save guarda un empleado en DynamoDB
func (r *DynamoDBRepository) Save(ctx context.Context, employee *domain.Employee) error {
	item, err := attributevalue.MarshalMap(employee)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		log.Printf("Error saving employee to DynamoDB: %v", err)
		return err
	}

	log.Printf("Employee saved successfully: %s", employee.ID)
	return nil
}

// FindByID busca un empleado por su ID
func (r *DynamoDBRepository) FindByID(ctx context.Context, id string) (*domain.Employee, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, domain.ErrNotFound
	}

	var employee domain.Employee
	err = attributevalue.UnmarshalMap(result.Item, &employee)
	if err != nil {
		return nil, err
	}

	return &employee, nil
}

// FindAll obtiene todos los empleados
func (r *DynamoDBRepository) FindAll(ctx context.Context) ([]*domain.Employee, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})

	if err != nil {
		return nil, err
	}

	var employees []*domain.Employee
	for _, item := range result.Items {
		var employee domain.Employee
		err := attributevalue.UnmarshalMap(item, &employee)
		if err != nil {
			log.Printf("Error unmarshaling employee: %v", err)
			continue
		}
		employees = append(employees, &employee)
	}

	return employees, nil
}
