package main

import (
	"context"
	"log"
	"logger-service/internal/application"
	"logger-service/internal/infrastructure"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		tableName = "employee-logs"
	}

	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL environment variable is required")
	}

	// Crear instancias de infraestructura
	repository := infrastructure.NewDynamoDBLogRepository(dynamoClient, tableName)
	consumer := infrastructure.NewSQSEventConsumer(sqsClient, queueURL)

	// Crear servicio de aplicación
	service := application.NewLoggerService(repository, consumer)

	// Manejar señales de interrupción
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down logger service...")
		cancel()
	}()

	// Iniciar consumo de eventos
	log.Println("Logger service starting...")
	if err := service.StartConsuming(ctx); err != nil {
		if err != context.Canceled {
			log.Fatalf("Error consuming events: %v", err)
		}
	}

	log.Println("Logger service stopped")
}
