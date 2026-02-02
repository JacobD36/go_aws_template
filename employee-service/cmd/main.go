package main

import (
	"context"
	"employee-service/internal/application"
	"employee-service/internal/infrastructure"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx := context.Background()

	// Configurar AWS SDK
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}

	// Crear clientes de AWS
	dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if awsEndpoint != "" {
			o.BaseEndpoint = aws.String(awsEndpoint)
		}
	})

	sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		if awsEndpoint != "" {
			o.BaseEndpoint = aws.String(awsEndpoint)
		}
	})

	// Obtener variables de entorno
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		tableName = "employees"
	}

	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL environment variable is required")
	}

	// Crear instancias de infraestructura
	repository := infrastructure.NewDynamoDBRepository(dynamoClient, tableName)
	publisher := infrastructure.NewSQSEventPublisher(sqsClient, queueURL)

	// Crear servicio de aplicación
	service := application.NewEmployeeService(repository, publisher)

	// Crear manejador HTTP
	handler := infrastructure.NewHTTPHandler(service)
	router := handler.SetupRoutes()

	// Iniciar servidor
	log.Println("Employee service starting on port 8081...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal(err)
	}
}
