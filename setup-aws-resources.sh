#!/bin/bash

echo "Esperando a que LocalStack esté completamente listo..."

# Esperar a que LocalStack esté healthy
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if curl -sf http://localhost:4566/_localstack/health > /dev/null 2>&1; then
        echo "✓ LocalStack está listo"
        break
    fi
    attempt=$((attempt + 1))
    echo "Esperando... intento $attempt/$max_attempts"
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "✗ Error: LocalStack no está disponible después de esperar"
    exit 1
fi

# Esperar 2 segundos adicionales para asegurar que todos los servicios internos estén listos
sleep 2

echo ""
echo "Creando cola SQS..."
aws --endpoint-url=http://localhost:4566 sqs create-queue \
    --queue-name employee-queue \
    --region us-east-1 \
    --no-cli-pager 2>/dev/null || echo "Cola SQS ya existe o error al crear"

echo ""
echo "Creando tabla DynamoDB para empleados..."
aws --endpoint-url=http://localhost:4566 dynamodb create-table \
    --table-name employees \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region us-east-1 \
    --no-cli-pager 2>/dev/null || echo "Tabla employees ya existe o error al crear"

echo ""
echo "Creando tabla DynamoDB para logs..."
aws --endpoint-url=http://localhost:4566 dynamodb create-table \
    --table-name employee-logs \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region us-east-1 \
    --no-cli-pager 2>/dev/null || echo "Tabla employee-logs ya existe o error al crear"

echo ""
echo "=========================================="
echo "✓ Recursos AWS creados exitosamente!"
echo "=========================================="
echo ""
echo "Verificando recursos creados..."
echo ""
echo "Colas SQS:"
aws --endpoint-url=http://localhost:4566 sqs list-queues \
    --region us-east-1 \
    --no-cli-pager 2>/dev/null || echo "Error al listar colas"

echo ""
echo "Tablas DynamoDB:"
aws --endpoint-url=http://localhost:4566 dynamodb list-tables \
    --region us-east-1 \
    --no-cli-pager 2>/dev/null || echo "Error al listar tablas"

echo ""
echo "=========================================="
echo "Sistema listo para usar!"
echo "=========================================="
