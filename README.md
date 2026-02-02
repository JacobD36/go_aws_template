# Sistema de Registro de Empleados con Microservicios

Sistema distribuido de registro de empleados utilizando microservicios en Go con arquitectura hexagonal, AWS SQS, DynamoDB y LocalStack.

## 📋 Descripción

Este proyecto implementa un sistema de registro de empleados usando:

- **API Gateway**: Punto de entrada HTTP con endpoints REST
- **Employee Service**: Microservicio que gestiona el registro de empleados
- **Logger Service**: Microservicio que procesa eventos y genera logs
- **LocalStack**: Emulador local de servicios AWS (SQS y DynamoDB)

## 🏗️ Arquitectura

El sistema sigue los principios de arquitectura hexagonal (puertos y adaptadores) y SOLID:

```
api-gateway/
├── main.go
├── go.mod
└── Dockerfile

employee-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/          # Entidades de negocio
│   ├── application/     # Casos de uso
│   ├── ports/           # Interfaces (puertos)
│   └── infrastructure/  # Adaptadores (DynamoDB, SQS, HTTP)
├── go.mod
└── Dockerfile

logger-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/
│   ├── application/
│   ├── ports/
│   └── infrastructure/
├── go.mod
└── Dockerfile
```

## 🚀 Requisitos Previos

- Docker y Docker Compose
- Go 1.21 o superior (para desarrollo local)
- AWS CLI (opcional, para interactuar con LocalStack)

## 📦 Instalación y Configuración

### 1. Inicializar módulos Go

Navega a cada directorio de servicio y ejecuta:

```bash
# API Gateway
cd api-gateway
go mod init api-gateway
go mod tidy
cd ..

# Employee Service
cd employee-service
go mod init employee-service
go mod tidy
cd ..

# Logger Service
cd logger-service
go mod init logger-service
go mod tidy
cd ..
```

### 2. Iniciar servicios con Docker Compose

Desde la raíz del proyecto:

```bash
docker-compose up -d
```

Esto iniciará:
- LocalStack (puerto 4566)
- API Gateway (puerto 8080)
- Employee Service (puerto 8081)
- Logger Service (proceso en background)

### 3. Configurar recursos AWS en LocalStack

Espera unos segundos para que LocalStack esté listo, luego ejecuta:

```bash
# Crear cola SQS
aws --endpoint-url=http://localhost:4566 sqs create-queue \
    --queue-name employee-queue \
    --region us-east-1

# Crear tabla DynamoDB para empleados
aws --endpoint-url=http://localhost:4566 dynamodb create-table \
    --table-name employees \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region us-east-1

# Crear tabla DynamoDB para logs
aws --endpoint-url=http://localhost:4566 dynamodb create-table \
    --table-name employee-logs \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region us-east-1
```

### 4. Verificar que los servicios están corriendo

```bash
docker-compose ps
```

## 🔧 Uso del Sistema

### Crear un nuevo empleado (POST)

```bash
curl -X POST http://localhost:8080/api/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan.perez@example.com"
  }'
```

Respuesta esperada:
```json
{
  "id": "uuid-generated",
  "name": "Juan Pérez",
  "email": "juan.perez@example.com",
  "created_at": "2026-01-27T10:30:00Z"
}
```

### Obtener todos los empleados (GET)

```bash
curl http://localhost:8080/api/employees
```

### Ver logs del Logger Service

El Logger Service muestra en consola cada evento procesado:

```bash
docker-compose logs -f logger-service
```

Salida esperada:
```
========================================
EVENTO RECIBIDO: employee.created
ID Empleado: uuid-generated
Nombre: Juan Pérez
Email: juan.perez@example.com
Timestamp del evento: 2026-01-27 10:30:00
Procesado el: 2026-01-27 10:30:01
========================================
```

## 🛠️ Desarrollo Local (sin Docker)

### 1. Iniciar LocalStack

```bash
docker run -d \
  --name localstack \
  -p 4566:4566 \
  -e SERVICES=sqs,dynamodb \
  localstack/localstack
```

### 2. Configurar recursos AWS (ver paso 3 de instalación)

### 3. Iniciar servicios manualmente

```bash
# Terminal 1 - Employee Service
cd employee-service
export AWS_ENDPOINT=http://localhost:4566
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export SQS_QUEUE_URL=http://localhost:4566/000000000000/employee-queue
export DYNAMODB_TABLE=employees
go run cmd/main.go

# Terminal 2 - Logger Service
cd logger-service
export AWS_ENDPOINT=http://localhost:4566
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export SQS_QUEUE_URL=http://localhost:4566/000000000000/employee-queue
export DYNAMODB_TABLE=employee-logs
go run cmd/main.go

# Terminal 3 - API Gateway
cd api-gateway
export EMPLOYEE_SERVICE_URL=http://localhost:8081
go run main.go
```

## 📊 Flujo de Datos

1. Cliente envía POST a `/api/employees` en API Gateway
2. API Gateway reenvía la petición a Employee Service
3. Employee Service:
   - Valida los datos
   - Guarda el empleado en DynamoDB (tabla `employees`)
   - Publica evento en SQS
4. Logger Service:
   - Consume el evento de SQS
   - Guarda log en DynamoDB (tabla `employee-logs`)
   - Muestra información en consola

## 🧪 Verificar tablas DynamoDB

```bash
# Ver empleados
aws --endpoint-url=http://localhost:4566 dynamodb scan \
    --table-name employees \
    --region us-east-1

# Ver logs
aws --endpoint-url=http://localhost:4566 dynamodb scan \
    --table-name employee-logs \
    --region us-east-1
```

## 🛑 Detener servicios

```bash
docker-compose down
```

Para limpiar también los volúmenes:

```bash
docker-compose down -v
```

## 📝 Principios de Diseño Aplicados

### Arquitectura Hexagonal
- **Domain**: Entidades y lógica de negocio pura
- **Application**: Casos de uso (orquestación)
- **Ports**: Interfaces que definen contratos
- **Infrastructure**: Implementaciones concretas (adaptadores)

### Principios SOLID
- **Single Responsibility**: Cada capa tiene una responsabilidad única
- **Open/Closed**: Extensible mediante nuevos adaptadores
- **Liskov Substitution**: Las interfaces permiten intercambiar implementaciones
- **Interface Segregation**: Interfaces específicas y cohesivas
- **Dependency Inversion**: Dependencias apuntan hacia abstracciones

## 🔍 Troubleshooting

### Error: "Cannot connect to LocalStack"
Espera 10-15 segundos después de `docker-compose up` antes de crear los recursos.

### Error: "Queue does not exist"
Asegúrate de haber creado la cola SQS con los comandos AWS CLI.

### Error: "Table not found"
Verifica que las tablas DynamoDB se crearon correctamente.

### Ver logs de un servicio específico
```bash
docker-compose logs -f [service-name]
# Ejemplo: docker-compose logs -f employee-service
```

## 📄 Licencia

Este proyecto es de código abierto y está disponible bajo la licencia MIT.
