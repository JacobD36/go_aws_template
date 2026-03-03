package main

import (
	"context"
	"log"
	"messaging-service/internal/application"
	"messaging-service/internal/infrastructure"
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
		tableName = "messages"
	}

	employeeEventsQueueURL := os.Getenv("EMPLOYEE_EVENTS_QUEUE_URL")
	if employeeEventsQueueURL == "" {
		log.Fatal("EMPLOYEE_EVENTS_QUEUE_URL environment variable is required")
	}

	logQueueURL := os.Getenv("LOG_QUEUE_URL")
	if logQueueURL == "" {
		log.Fatal("LOG_QUEUE_URL environment variable is required")
	}

	// Crear instancias de infraestructura (Dependency Injection)
	repository := infrastructure.NewDynamoDBRepository(dynamoClient, tableName)
	sender := infrastructure.NewSimulatedMessageSender()
	publisher := infrastructure.NewSQSEventPublisher(sqsClient, logQueueURL)
	consumer := infrastructure.NewSQSEventConsumer(sqsClient, employeeEventsQueueURL)

	// Crear servicio de aplicación
	service := application.NewMessagingService(repository, sender, publisher)

	// Configurar manejo de señales para graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping service...")
		cancel()
	}()

	// Iniciar consumidor de eventos
	log.Println("Messaging service starting...")
	log.Printf("Consuming events from: %s", employeeEventsQueueURL)
	log.Printf("Publishing logs to: %s", logQueueURL)

	if err := consumer.ConsumeEvents(ctx, service.HandleEmployeeEvent); err != nil {
		if err != context.Canceled {
			log.Fatalf("Error consuming events: %v", err)
		}
	}

	log.Println("Messaging service stopped gracefully")
}
