#!/bin/bash

# Script para inicializar los recursos de AWS en LocalStack

echo "Esperando a que LocalStack esté listo..."
sleep 10

echo "Creando cola SQS..."
aws --endpoint-url=http://localhost:4566 sqs create-queue \
    --queue-name employee-queue \
    --region us-east-1

echo "Creando tabla DynamoDB para empleados..."
aws --endpoint-url=http://localhost:4566 dynamodb create-table \
    --table-name employees \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region us-east-1

echo "Creando tabla DynamoDB para logs..."
aws --endpoint-url=http://localhost:4566 dynamodb create-table \
    --table-name employee-logs \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region us-east-1

echo "¡Recursos AWS creados exitosamente!"
echo ""
echo "Verificando recursos..."
echo ""
echo "Colas SQS:"
aws --endpoint-url=http://localhost:4566 sqs list-queues --region us-east-1
echo ""
echo "Tablas DynamoDB:"
aws --endpoint-url=http://localhost:4566 dynamodb list-tables --region us-east-1
