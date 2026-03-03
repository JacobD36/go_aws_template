# Sistema de Registro de Empleados con Microservicios

Sistema distribuido de registro de empleados utilizando microservicios en Go con arquitectura hexagonal, AWS SQS, DynamoDB, LocalStack y frontend en React/Next.js.

## 📋 Descripción

Este proyecto implementa un sistema de registro de empleados usando:

- **Frontend Basic**: Portal administrativo web con React, Next.js y TypeScript
- **API Gateway**: Punto de entrada HTTP con endpoints REST
- **Employee Service**: Microservicio que gestiona el registro de empleados y publica eventos de creación
- **Auth Service**: Microservicio de autenticación que genera tokens JWT
- **Messaging Service**: Microservicio que consume eventos de empleados creados y simula envío de emails de bienvenida
- **Logger Service**: Microservicio que consume eventos de mensajería y genera logs auditables
- **LocalStack**: Emulador local de servicios AWS (SQS y DynamoDB)

## 🏗️ Arquitectura

El sistema sigue los principios de arquitectura hexagonal (puertos y adaptadores) y SOLID:

```
frontend_basic/
├── app/                 # Next.js App Router
│   ├── login/          # Página de login
│   └── employees/      # Gestión de empleados
├── components/         # Componentes React
│   ├── ui/            # Componentes base
│   ├── layout/        # Layout del dashboard
│   └── employees/     # Componentes de empleados
├── lib/               # Utilidades y servicios
│   ├── api.ts        # Cliente API
│   ├── auth.ts       # Autenticación
│   └── constants.ts  # Constantes
└── types/            # Tipos TypeScript

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
│   │   ├── repository.go
│   │   ├── event_publisher.go
│   │   └── password_hasher.go    # 🔒 Puerto para hash de passwords
│   └── infrastructure/  # Adaptadores (DynamoDB, SQS, HTTP, Bcrypt)
│       ├── dynamodb_repository.go
│       ├── sqs_publisher.go
│       ├── http_handler.go
│       └── bcrypt_hasher.go       # 🔒 Implementación con bcrypt
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

messaging-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/          # Entidades (Message, EmployeeEvent, LogEvent)
│   ├── application/     # Lógica de mensajería (ProcessEmployeeCreatedEvent)
│   ├── ports/           # Interfaces (puertos)
│   │   ├── repository.go
│   │   ├── event_consumer.go     # 📨 Puerto para consumir eventos
│   │   ├── event_publisher.go    # 📤 Puerto para publicar eventos
│   │   └── message_sender.go     # ✉️ Puerto para enviar mensajes
│   └── infrastructure/  # Adaptadores (SQS, DynamoDB, Simulador)
│       ├── sqs_consumer.go        # 📨 Consumo de employee-events-queue
│       ├── sqs_publisher.go       # 📤 Publicación a employee-queue
│       ├── simulated_sender.go    # ✉️ Simulador de email/SMS
│       └── dynamodb_repository.go
├── go.mod
└── Dockerfile

auth-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/          # Entidades (User, LoginCredentials, AuthToken)
│   ├── application/     # Lógica de autenticación (Login, ValidateToken)
│   ├── ports/           # Interfaces (puertos)
│   │   ├── repository.go          # 🔍 Puerto para buscar usuarios
│   │   ├── password_hasher.go    # 🔒 Puerto para comparar passwords
│   │   └── token_generator.go    # 🔐 Puerto para generar JWT
│   └── infrastructure/  # Adaptadores (DynamoDB, Bcrypt, JWT, HTTP)
│       ├── dynamodb_repository.go # 🔍 Búsqueda en DynamoDB
│       ├── bcrypt_hasher.go       # 🔒 Comparación con bcrypt
│       ├── jwt_token_generator.go # 🔐 Generación de JWT
│       └── http_handler.go        # 🌐 Endpoints: /auth/login, /health
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

# Auth Service
cd auth-service
go mod init auth-service
go mod tidy
cd ..
go mod tidy
cd ..
```

- Frontend Basic (puerto 3000)
### 2. Iniciar servicios con Docker Compose

Desde la raíz del proyecto:

```bash
docker-compose up -d
```

Esto iniciará:
- LocalStack (puerto 4566)
- API Gateway (puerto 8080)
- Employee Service (puerto 8081)
- Auth Service (puerto 8082)
- Messaging Service (proceso en background)
- Logger Service (proceso en background)
- Frontend Basic (puerto 3000)

### 3. Configurar recursos AWS en LocalStack

Espera unos segundos para que LocalStack esté listo, luego ejecuta:

```bash
# Crear colas SQS
aws --endpoint-url=http://localhost:4566 sqs create-queue \
    --queue-name employee-queue \
    --region us-east-1

aws --endpoint-url=http://localhost:4566 sqs create-queue \
    --queue-name employee-events-queue \
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

# Crear tabla DynamoDB para mensajes
aws --endpoint-url=http://localhost:4566 dynamodb create-table \
    --table-name messages \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region us-east-1
```

**Nota:** También puedes usar el script automatizado:
```bash
# Desde la raíz del proyecto
./init-aws-resources.sh
# o
./setup-aws-resources.sh
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
    "email": "juan.perez@example.com",
    "password": "SecurePass123!"
  }'
```

**Requisitos del Password:**
- Mínimo 8 caracteres
- Al menos una letra mayúscula
- Al menos un número
- Al menos un caracter especial (!@#$%^&*()_+-=[]{};\':"|,.<>/?~)

Respuesta esperada:
```json
{
  "id": "uuid-generated",
  "name": "Juan Pérez",
  "email": "juan.perez@example.com",
  "created_at": "2026-01-27T10:30:00Z"
}
```

**Nota de Seguridad:** El password nunca se devuelve en las respuestas ni aparece en los logs.

### Obtener todos los empleados (GET)

```bash
curl http://localhost:8080/api/employees
```

### Autenticación (Login)

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan.perez@example.com",
    "password": "SecurePass123!"
  }'
```

Respuesta esperada:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "uuid-generated",
  "expires_at": 1738384800
}
```

**Nota:** El token JWT contiene únicamente el ID del usuario y se puede usar para autenticar peticiones HTTP. El endpoint `/api/auth/login` en el API Gateway reenvía las peticiones al Auth Service.

### 🌐 Usar el Frontend (Interfaz Web)

El portal administrativo está disponible en: **http://localhost:3000**

#### Características del Frontend:
- **Login**: Página de autenticación con formulario validado
- **Dashboard**: Portal con menú lateral
- **Gestión de Empleados**:
  - Lista de empleados en formato tabla
  - Botón "Nuevo Empleado" que abre un modal
  - Validación de formularios en tiempo real
  - Estados de carga y error
  - Mensajes de éxito y error

#### Flujo de uso:
1. Abre http://localhost:3000 en tu navegador
2. Inicia sesión con credenciales de un empleado registrado
3. Serás redirigido al dashboard con la lista de empleados
4. Usa el botón "Nuevo Empleado" para registrar nuevos empleados
5. Los botones editar/eliminar están implementados solo visualmente

**Nota:** El frontend se comunica con el API Gateway en el puerto 8080. Asegúrate de tener todos los servicios corriendo.

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

# Terminal 2 - Messaging Service
cd messaging-service
export AWS_ENDPOINT=http://localhost:4566
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export EMPLOYEE_EVENTS_QUEUE_URL=http://localhost:4566/000000000000/employee-events-queue
export LOG_QUEUE_URL=http://localhost:4566/000000000000/employee-queue
export DYNAMODB_TABLE=messages
go run cmd/main.go

# Terminal 3 - Logger Service
cd logger-service
export AWS_ENDPOINT=http://localhost:4566
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export SQS_QUEUE_URL=http://localhost:4566/000000000000/employee-queue
export DYNAMODB_TABLE=employee-logs
go run cmd/main.go

# Terminal 4 - Auth Service
cd auth-service
export AWS_ENDPOINT=http://localhost:4566
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export DYNAMODB_TABLE=employees
export JWT_SECRET=my-super-secret-jwt-key
export JWT_EXPIRATION_MINUTES=60
export PORT=8082
go run cmd/main.go

# Terminal 5 - API Gateway
cd api-gateway
export EMPLOYEE_SERVICE_URL=http://localhost:8081
export AUTH_SERVICE_URL=http://localhost:8082
go run main.go
```

## 📊 Flujo de Datos

### Registro de Empleados (Event-Driven Architecture)

```
┌─────────────────┐
│   API Gateway   │  1. POST /api/employees
│  (puerto 8080)  │
└────────┬────────┘
         │
         ▼
┌─────────────────────┐
│ Employee Service    │  2. Valida, hashea password, guarda en DB
│   (puerto 8081)     │  3. Publica evento employee.created
└──────────┬──────────┘
           │
           ▼
    [employee-events-queue]  ← SQS
           │
           ▼
┌─────────────────────┐
│ Messaging Service   │  4. Consume evento
│   (background)      │  5. Simula envío de email
└──────────┬──────────┘  6. Publica evento message.sent
           │
           ▼
    [employee-queue]  ← SQS
           │
           ▼
┌─────────────────────┐
│  Logger Service     │  7. Consume evento
│   (background)      │  8. Guarda log en DynamoDB
└─────────────────────┘
```

**Paso a paso:**

1. **Cliente → API Gateway**: POST a `/api/employees` con datos del empleado
2. **API Gateway → Employee Service**: Reenvía la petición
3. **Employee Service**:
   - Valida los datos y complejidad del password
   - Hashea el password con bcrypt
   - Guarda el empleado en DynamoDB (tabla `employees`)
   - Publica evento `employee.created` a `employee-events-queue` (SQS)
4. **Messaging Service** (consumidor asíncrono):
   - Consume el evento desde `employee-events-queue`
   - Crea un mensaje de bienvenida
   - **Simula el envío del email** (log en consola)
   - Guarda el mensaje en DynamoDB (tabla `messages`)
   - Publica evento `message.sent` a `employee-queue` (SQS)
5. **Logger Service** (consumidor asíncrono):
   - Consume el evento desde `employee-queue`
   - Guarda log auditable en DynamoDB (tabla `employee-logs`)
   - Muestra información en consola

**Estructura del evento `employee.created`:**
```json
{
  "event_type": "employee.created",
  "employee": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "created_at": "2026-03-02T19:00:00Z"
  },
  "timestamp": "2026-03-02T19:00:00Z"
}
```
⚠️ **Nota de Seguridad**: El password hasheado NO se incluye en el evento.

**Estructura del evento `message.sent`:**
```json
{
  "event_type": "message.sent",
  "action": "email_sent",
  "entity_id": "msg-uuid",
  "details": "Email de bienvenida enviado a juan@example.com",
  "timestamp": "2026-03-02T19:00:01Z"
}
```

### Autenticación (Login)
1. Cliente envía POST a `/api/auth/login` con email y password
2. API Gateway reenvía la petición a Auth Service
3. Auth Service:
   - Valida que email y password no estén vacíos
   - Busca el usuario por email en DynamoDB (tabla `employees`)
   - Compara el password ingresado con el hash almacenado usando bcrypt
   - Genera un token JWT que contiene el ID del usuario
   - Retorna el token con tiempo de expiración (60 minutos por defecto)

**Respuesta exitosa:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "expires_at": 1738384800
}
```

### Ventajas de la Arquitectura Event-Driven
- ✅ **Asíncrona**: La respuesta al cliente no espera al envío del email
- ✅ **Desacoplada**: Los servicios se comunican solo por eventos
- ✅ **Escalable**: Cada servicio puede escalar independientemente
- ✅ **Resiliente**: Si un servicio falla, los mensajes permanecen en SQS
- ✅ **Auditable**: Todo el flujo queda registrado en logs
- ✅ **Extensible**: Se pueden agregar nuevos consumidores sin modificar servicios existentes

## 🧪 Verificar tablas DynamoDB

```bash
# Ver empleados/usuarios
aws --endpoint-url=http://localhost:4566 dynamodb scan \
    --table-name employees \
    --region us-east-1

# Ver logs
aws --endpoint-url=http://localhost:4566 dynamodb scan \
    --table-name employee-logs \
    --region us-east-1

# Ver mensajes enviados
aws --endpoint-url=http://localhost:4566 dynamodb scan \
    --table-name messages \
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

### Patrones de Diseño Implementados

#### Strategy Pattern + Dependency Inversion (Hash de Passwords & JWT)
El sistema implementa hash de passwords y generación de JWT aplicando arquitectura hexagonal:

**Estructura en Employee Service:**
```
ports/
  └── password_hasher.go      # Puerto (interfaz para hash)
infrastructure/
  └── bcrypt_hasher.go        # Adaptador (implementación con bcrypt)
application/
  └── employee_service.go     # Inyección de dependencia
```

**Estructura en Auth Service:**
```
ports/
  ├── password_hasher.go      # Puerto (interfaz para comparar passwords)
  ├── repository.go           # Puerto (búsqueda de usuarios)
  └── token_generator.go      # Puerto (interfaz para generar JWT)
infrastructure/
  ├── bcrypt_hasher.go        # Adaptador (comparación con bcrypt)
  ├── dynamodb_repository.go  # Adaptador (DynamoDB)
  └── jwt_token_generator.go  # Adaptador (JWT con golang-jwt/jwt)
application/
  └── auth_service.go         # Inyección de dependencias
```

**Principios aplicados:**
1. **Dependency Inversion Principle (DIP)**: 
   - Los servicios dependen de abstracciones (`PasswordHasher`, `TokenGenerator`, `UserRepository`)
   - No dependen de implementaciones concretas
   - El dominio permanece puro sin conocer bcrypt ni JWT

2. **Strategy Pattern**:
   - Los algoritmos (hash, JWT) están encapsulados en estrategias intercambiables
   - Se puede cambiar de bcrypt a argon2, o de JWT a OAuth sin modificar servicios
   - Solo se crea un nuevo adaptador que implemente el puerto

3. **Ports & Adapters (Hexagonal)**:
   - Los puertos son interfaces en la capa de dominio/aplicación
   - Los adaptadores son implementaciones en infraestructura
   - Inyección de dependencias en constructores de servicios

**Beneficios:**
- ✅ Fácil de testear (mock de hasher, token generator, repository)
- ✅ Extensible (nuevos algoritmos sin cambiar código existente)
- ✅ Dominio independiente de librerías externas
- ✅ Cumple Open/Closed Principle
- ✅ Auth Service reutiliza la misma arquitectura que Employee Service

## 🔒 Seguridad

### Gestión de Passwords
El sistema implementa las siguientes medidas de seguridad para los passwords:

- **Validación de Complejidad**: Los passwords deben cumplir requisitos estrictos:
  - Mínimo 8 caracteres
  - Al menos una letra mayúscula (A-Z)
  - Al menos un número (0-9)
  - Al menos un caracter especial (!@#$%^&* etc.)

- **Hash con Bcrypt**: 
  - Los passwords se hashean usando bcrypt (cost factor 10)
  - Implementado mediante el patrón Strategy y arquitectura hexagonal
  - Los passwords nunca se almacenan en texto plano
  - El hash es irreversible y único por cada password (salt automático)
  - Se guarda solo el hash en DynamoDB

- **Protección en Respuestas**: 
  - El password (hasheado) nunca se serializa en JSON (tag `json:"-"`)
  - No aparece en logs del sistema
  - El endpoint de creación devuelve un objeto `EmployeePublic` sin password
  - El endpoint de listado devuelve arrays de `EmployeePublic` sin passwords

### Autenticación con JWT
El sistema de autenticación implementa las siguientes medidas:

- **Tokens JWT (JSON Web Tokens)**:
  - Generados por el Auth Service tras login exitoso
  - Contienen únicamente el ID del usuario (sin datos sensibles)
  - Firmados con HS256 (HMAC-SHA256)
  - Expiración configurable (60 minutos por defecto)
  - Secret key configurable via variable de entorno `JWT_SECRET`

- **Validaciones en Login**:
  - Email y password son obligatorios
  - Búsqueda de usuario en DynamoDB por email
  - Comparación de password con hash usando bcrypt
  - Retorna 401 Unauthorized si las credenciales son inválidas
  - Retorna 400 Bad Request si faltan datos

- **Estructura del Token**:
  ```json
  {
    "user_id": "uuid-del-usuario",
    "exp": 1738384800,
    "iat": 1738381200
  }
  ```

- **Uso del Token**:
  - El token se puede incluir en headers HTTP: `Authorization: Bearer <token>`
  - Se puede validar usando el método `ValidateToken` del Auth Service
  - Retorna el ID del usuario si el token es válido

**Ejemplo de validación:**
```bash
# Password inválido (falta mayúscula)
curl -X POST http://localhost:8080/api/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "weak123!"
  }'
# Error: invalid password: must be at least 8 characters with at least one uppercase letter, one number, and one special character

# Password válido
curl -X POST http://localhost:8080/api/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
# Success: devuelve empleado sin password
```
## 🔐 Microservicio Auth Service

### Descripción

El **Auth Service** es un microservicio independiente responsable de la autenticación y registro de usuarios. Implementa:

- ✅ Registro de nuevos usuarios con validación de password
- ✅ Hash de passwords con bcrypt antes de almacenarlos
- ✅ Publicación de eventos `user.created` a la cola SQS
- ✅ Validación de credenciales (email y password obligatorios)
- ✅ Búsqueda de usuarios en DynamoDB por email
- ✅ Comparación segura de passwords usando bcrypt
- ✅ Generación de tokens JWT con el ID del usuario
- ✅ Validación de tokens JWT
- ✅ Arquitectura hexagonal con puertos y adaptadores

### Endpoints Disponibles

#### POST /auth/login
Autentica un usuario y genera un token JWT.

**Request:**
```json
{
  "email": "usuario@example.com",
  "password": "SecurePass123!"
}
```

**Response exitosa (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "uuid-del-usuario",
  "expires_at": 1738384800
}
```

**Errores posibles:**
- `400 Bad Request`: Email o password faltante
- `401 Unauthorized`: Credenciales inválidas
- `500 Internal Server Error`: Error del servidor

#### POST /auth/register
Registra un nuevo usuario y publica un evento para enviar mensaje de bienvenida.

**Request:**
```json
{
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "password": "SecurePass123!"
}
```

**Validaciones:**
- Email y password son obligatorios
- Password debe cumplir requisitos de seguridad:
  - Mínimo 8 caracteres
  - Al menos una letra mayúscula
  - Al menos un número
  - Al menos un caracter especial

**Response exitosa (201):**
```json
{
  "id": "uuid-del-usuario",
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "created_at": "2024-01-20T10:30:00Z"
}
```

**Eventos publicados:**
- Publica evento `user.created` a la cola `user-queue`
- El Messaging Service procesa este evento y envía un email de bienvenida (simulado)

**Errores posibles:**
- `400 Bad Request`: Datos inválidos o password débil
- `409 Conflict`: El email ya está registrado
- `500 Internal Server Error`: Error del servidor

#### GET /health
Verifica el estado del servicio.

**Response (200):**
```json
{
  "status": "ok"
}
```

### Configuración

El servicio se configura mediante variables de entorno:

```bash
# DynamoDB
AWS_REGION=us-east-1
AWS_ENDPOINT=http://localstack:4566
DYNAMODB_TABLE=employees

# JWT
JWT_SECRET=my-super-secret-jwt-key-change-in-production
JWT_EXPIRATION_MINUTES=60

# Servidor
PORT=8082
```

### Arquitectura Interna

El Auth Service sigue arquitectura hexagonal con 4 puertos principales:

1. **UserRepository** (puerto): Interfaz para buscar y guardar usuarios
   - **Adaptador**: `DynamoDBUserRepository` - Busca por email y guarda en DynamoDB

2. **PasswordHasher** (puerto): Interfaz para hashear y comparar passwords
   - **Adaptador**: `BcryptPasswordHasher` - Hash y comparación con bcrypt

3. **TokenGenerator** (puerto): Interfaz para generar/validar JWT
   - **Adaptador**: `JWTTokenGenerator` - Usa golang-jwt/jwt/v5

4. **EventPublisher** (puerto): Interfaz para publicar eventos
   - **Adaptador**: `SQSPublisher` - Publica a `user-queue`

### Flujo de Autenticación

```
Login:
1. Usuario → POST /auth/login {email, password}
2. Auth Service valida que email y password no estén vacíos
3. Busca usuario en DynamoDB por email
4. Compara password con hash almacenado (bcrypt)
5. Si coincide: Genera JWT con user_id
6. Retorna {token, user_id, expires_at}
```

### Flujo de Registro

```
Registro:
1. Usuario → POST /auth/register {name, email, password}
2. Auth Service valida datos y complejidad del password
3. Hashea el password con bcrypt
4. Guarda usuario en DynamoDB (tabla employees)
5. Publica evento user.created a user-queue
6. Retorna datos del usuario (sin password)
7. Messaging Service consume el evento
8. Genera y simula envío de email de bienvenida
9. Guarda el mensaje en DynamoDB (tabla messages)
10. Publica evento message.sent para logging
```

### Estructura del Token JWT

El token contiene únicamente el ID del usuario (sin datos sensibles):

```json
{
  "user_id": "uuid-del-usuario",
  "exp": 1738384800,  // Timestamp de expiración
  "iat": 1738381200   // Timestamp de emisión
}
```

### Ejemplo de Uso Completo

```bash
# 1. Registrar un nuevo usuario (con Auth Service)
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "password": "SecurePass123!"
  }'

# Respuesta:
{
  "id": "98bb3d68-db50-4b62-9f07-b8268e015182",
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "created_at": "2024-01-20T10:30:00Z"
}

# Este registro:
# - Crea el usuario en DynamoDB (tabla employees)
# - Publica evento user.created a user-queue
# - Messaging Service consume el evento y envía email de bienvenida (simulado)
# - Messaging Service publica evento a employee-queue para logging

# 2. Autenticarse
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan@example.com",
    "password": "SecurePass123!"
  }'

# Respuesta:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOThiYjNkNjgtZGI1MC00YjYyLTlmMDctYjgyNjhlMDE1MTgyIiwiZXhwIjoxNzM4Mzg0ODAwLCJpYXQiOjE3MzgzODEyMDB9.xyz",
  "user_id": "98bb3d68-db50-4b62-9f07-b8268e015182",
  "expires_at": 1738384800
}

# 3. Usar el token en peticiones autenticadas
curl -X GET http://localhost:8080/api/employees \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Consideraciones de Seguridad

- 🔒 **JWT Secret**: Cambiar `JWT_SECRET` en producción
- 🔒 **HTTPS**: Usar HTTPS en producción para proteger tokens
- 🔒 **Expiración**: Los tokens expiran después de 60 minutos
- 🔒 **Password**: Nunca se transmite ni almacena en texto plano
- 🔒 **Bcrypt**: Los passwords se comparan usando bcrypt (resistente a ataques de fuerza bruta)
## 📨 Microservicio Messaging Service

### Descripción

El **Messaging Service** es un microservicio basado en eventos que gestiona el envío de mensajes de bienvenida a nuevos usuarios. Implementa:

- ✅ Consumo de eventos `user.created` desde la cola SQS `user-queue`
- ✅ Simulación de envío de emails y SMS (logs en consola)
- ✅ Generación automática de mensajes de bienvenida personalizados
- ✅ Persistencia de mensajes enviados en DynamoDB (tabla `messages`)
- ✅ Publicación de eventos `message.sent` a `employee-queue` para logging
- ✅ Arquitectura hexagonal con puertos y adaptadores
- ✅ Principios SOLID y Clean Code

### Arquitectura Interna

El Messaging Service sigue arquitectura hexagonal con 4 puertos principales:

1. **EventConsumer** (puerto): Interfaz para consumir eventos
   - **Adaptador**: `SQSConsumer` - Consume desde `employee-events-queue`

2. **MessageSender** (puerto): Interfaz para enviar mensajes
   - **Adaptador**: `SimulatedSender` - Simula envío de email/SMS

3. **MessageRepository** (puerto): Interfaz para persistir mensajes
   - **Adaptador**: `DynamoDBMessageRepository` - Guarda en tabla `messages`

4. **EventPublisher** (puerto): Interfaz para publicar eventos
   - **Adaptador**: `SQSPublisher` - Publica a `employee-queue`

### Flujo de Procesamiento

```
1. Empleado creado → Employee Service publica evento employee.created a employee-events-queue
2. Messaging Service consume evento desde employee-events-queue
3. Crea mensaje de bienvenida personalizado con el nombre del empleado
4. Simula envío del email (log en consola):
   "📧 Simulando envío de EMAIL a juan@example.com..."
5. Guarda el mensaje en DynamoDB (tabla messages)
6. Publica evento message.sent a employee-queue
7. Logger Service registra el evento
```

### Configuración

El servicio se configura mediante variables de entorno:

```bash
# AWS
AWS_REGION=us-east-1
AWS_ENDPOINT=http://localstack:4566

# SQS
EMPLOYEE_EVENTS_QUEUE_URL=http://localstack:4566/000000000000/employee-events-queue
LOG_QUEUE_URL=http://localstack:4566/000000000000/employee-queue

# DynamoDB
DYNAMODB_TABLE=messages
```

### Tipos de Mensajes

El servicio soporta dos tipos de mensajes (preparado para extensión):

- **EMAIL**: Envío de correos electrónicos (actualmente simulado)
- **SMS**: Envío de mensajes de texto (actualmente simulado)

### Ejemplo de Mensaje de Bienvenida

```
Asunto: ¡Bienvenido a nuestro sistema!

Hola Juan Pérez,

¡Bienvenido a nuestro sistema! Nos alegra tenerte aquí.

Tu cuenta ha sido creada exitosamente con el email: juan@example.com

Puedes comenzar a usar nuestros servicios de inmediato.

Saludos,
El equipo
```

### Verificar Mensajes Enviados

```bash
# Ver todos los mensajes en DynamoDB
aws --endpoint-url=http://localhost:4566 dynamodb scan \
    --table-name messages \
    --region us-east-1

# Ver logs del Messaging Service
docker-compose logs -f messaging-service
```
## �🔍 Troubleshooting

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
# Ejemplo: docker-compose logs -f auth-service
# Ejemplo: docker-compose logs -f messaging-service
```

### Error en autenticación: "Invalid email or password"
Verifica que:
1. El usuario fue creado correctamente (con /api/auth/register o /api/employees)
2. El email es exacto (case-sensitive)
3. El password es correcto
4. El auth-service tiene acceso a la misma tabla DynamoDB (employees)

### No se reciben mensajes de bienvenida
Verifica que:
1. La cola `employee-events-queue` fue creada correctamente
2. El Messaging Service está corriendo: `docker-compose ps messaging-service`
3. Employee Service publica eventos a `employee-events-queue` después del registro
4. Revisa logs: `docker-compose logs -f messaging-service`
5. La tabla `messages` fue creada en DynamoDB

## 📑 Historial de Cambios

### Versión 2.0.0 - Corrección de Arquitectura Event-Driven (2 marzo 2026)

#### 🔄 Cambios en Employee-Service
- ✅ Modificada estructura de evento `EmployeeEvent` para usar `EmployeeEventData`
- ✅ Excluido password hasheado del evento publicado (seguridad)
- ✅ Convierte `CreatedAt` de `time.Time` a `string` (RFC3339) en eventos
- ✅ Publica eventos a `employee-events-queue` (sin información sensible)

#### 🔄 Cambios en Messaging-Service
- ✅ Renombrado `UserEvent` → `EmployeeEvent` en dominio
- ✅ Actualizado para consumir desde `employee-events-queue` (antes `user-queue`)
- ✅ Cambiado `ProcessUserCreatedEvent` → `ProcessEmployeeCreatedEvent`
- ✅ Cambiado `HandleUserEvent` → `HandleEmployeeEvent`
- ✅ Variable de entorno: `USER_QUEUE_URL` → `EMPLOYEE_EVENTS_QUEUE_URL`
- ✅ Mantiene publicación de logs a `employee-queue`

#### 🔄 Cambios en Auth-Service
- ✅ Eliminado endpoint `/auth/register` (solo autenticación)
- ✅ Eliminado archivo `internal/domain/user_event.go`
- ✅ Eliminado archivo `internal/ports/event_publisher.go`
- ✅ Eliminado archivo `internal/infrastructure/sqs_publisher.go`
- ✅ Eliminado método `Register` de `auth_service.go`
- ✅ Eliminado método `Save` de repository
- ✅ Eliminado método `Hash` de password_hasher
- ✅ Eliminado cliente SQS de `cmd/main.go`
- ✅ Auth-Service ahora solo maneja `Login` y `ValidateToken`

#### 🔄 Cambios en API-Gateway
- ✅ Eliminado endpoint `/api/auth/register`
- ✅ Endpoints disponibles: `/api/employees` (GET/POST), `/api/auth/login` (POST)

#### 🔄 Cambios en Infraestructura
- ✅ Actualizado `docker-compose.yml` con variables correctas
- ✅ Cola SQS: `user-queue` → `employee-events-queue`
- ✅ Actualizado `init-aws-resources.sh`
- ✅ Actualizado `setup-aws-resources.sh`
- ✅ Corregidos backslashes duplicados en scripts

#### 📊 Flujo Correcto Implementado
```
POST /api/employees
    ↓
Employee Service (crea + publica)
    ↓
employee-events-queue
    ↓
Messaging Service (simula email + publica)
    ↓
employee-queue
    ↓
Logger Service (registra)
```

#### ⚠️ Importante
- El registro de empleados ahora SOLO se realiza vía `/api/employees`
- Auth-Service es exclusivamente para autenticación (login/validación de tokens)
- Los eventos de empleados NO incluyen el password hasheado
- Todo el flujo de mensajería se dispara desde Employee-Service, no desde Auth-Service

## 📝 Notas Adicionales

### Tablas DynamoDB
- `employees`: Almacena empleados (ID, Name, Email, Password hasheado, CreatedAt)
- `employee-logs`: Almacena logs auditables de eventos
- `messages`: Almacena mensajes simulados enviados

### Colas SQS
- `employee-events-queue`: Eventos de empleados creados (Employee → Messaging)
- `employee-queue`: Eventos de logs (Messaging → Logger)

### Servicios y Puertos
- API Gateway: `8080`
- Employee Service: `8081`
- Auth Service: `8082`
- Frontend: `3000`
- LocalStack: `4566`
- Messaging Service: background (sin puerto HTTP)
- Logger Service: background (sin puerto HTTP)

## 💯 Validaciones y Seguridad

### Validación de Passwords
- Mínimo 8 caracteres
- Al menos una mayúscula (A-Z)
- Al menos un número (0-9)
- Al menos un caracter especial (!@#$%^&*)

### Seguridad de Passwords
- Hash con bcrypt (cost factor 10)
- Salt automático por password
- Nunca se devuelve en respuestas JSON
- No se incluye en eventos publicados
- No aparece en logs

### JWT Tokens
- Algoritmo: HS256 (HMAC-SHA256)
- Contiene solo user_id
- Expiración: 60 minutos (configurable)
- Secret key configurable vía `JWT_SECRET`

## 📄 Licencia

Este proyecto es de código abierto y está disponible bajo la licencia MIT.
