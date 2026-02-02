package main

import (
	"auth-service/internal/application"
	"auth-service/internal/infrastructure"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

	// Crear cliente de DynamoDB
	dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if awsEndpoint != "" {
			o.BaseEndpoint = aws.String(awsEndpoint)
		}
	})

	// Obtener variables de entorno
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		tableName = "employees"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my-secret-key-change-in-production"
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	jwtExpirationStr := os.Getenv("JWT_EXPIRATION_MINUTES")
	jwtExpiration := 60 // Default: 60 minutos
	if jwtExpirationStr != "" {
		if exp, err := strconv.Atoi(jwtExpirationStr); err == nil {
			jwtExpiration = exp
		}
	}

	// Crear instancias de infraestructura (adaptadores)
	repository := infrastructure.NewDynamoDBUserRepository(dynamoClient, tableName)
	passwordHasher := infrastructure.NewBcryptPasswordHasher()
	tokenGenerator := infrastructure.NewJWTTokenGenerator(jwtSecret, jwtExpiration)

	// Crear servicio de aplicación con inyección de dependencias
	service := application.NewAuthService(repository, passwordHasher, tokenGenerator)

	// Crear manejador HTTP
	handler := infrastructure.NewHTTPHandler(service)
	router := handler.SetupRoutes()

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Auth service starting on port %s...", port)
	log.Printf("JWT expiration: %d minutes", jwtExpiration)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
