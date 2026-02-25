# Sistema de Registro de Empleados con Microservicios

Sistema distribuido de registro de empleados utilizando microservicios en Go con arquitectura hexagonal, AWS SQS, DynamoDB, LocalStack y frontend en React/Next.js.

## 📋 Descripción

Este proyecto implementa un sistema de registro de empleados usando:

- **Frontend Basic**: Portal administrativo web con React, Next.js y TypeScript
- **API Gateway**: Punto de entrada HTTP con endpoints REST
- **Employee Service**: Microservicio que gestiona el registro de empleados
- **Auth Service**: Microservicio de autenticación que genera tokens JWT
- **Logger Service**: Microservicio que procesa eventos y genera logs
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

auth-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/          # Entidades (User, LoginCredentials, AuthToken)
│   ├── application/     # Lógica de autenticación
│   ├── ports/           # Interfaces (puertos)
│   │   ├── repository.go
│   │   ├── password_hasher.go    # 🔒 Puerto para comparar passwords
│   │   └── token_generator.go    # 🔐 Puerto para generar JWT
│   └── infrastructure/  # Adaptadores (DynamoDB, Bcrypt, JWT, HTTP)
│       ├── dynamodb_repository.go
│       ├── bcrypt_hasher.go       # 🔒 Comparación con bcrypt
│       ├── jwt_token_generator.go # 🔐 Generación de JWT
│       └── http_handler.go
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

# Terminal 2 - Logger Service
cd logger-service
export AWS_ENDPOINT=http://localhost:4566
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export SQS_QUEUE_URL=http://localhost:4566/000000000000/employee-queue
export DYNAMODB_TABLE=employee-logs
go run cmd/main.go

# Terminal 3 - Auth Service
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

# Terminal 4 - API Gateway
cd api-gateway
export EMPLOYEE_SERVICE_URL=http://localhost:8081
go run main.go
```

## 📊 Flujo de Datos

### Registro de Empleados
1. Cliente envía POST a `/api/employees` en API Gateway
2. API Gateway reenvía la petición a Employee Service
3. Employee Service:
   - Valida los datos y complejidad del password
   - Hashea el password con bcrypt
   - Guarda el empleado en DynamoDB (tabla `employees`)
   - Publica evento en SQS
4. Logger Service:
   - Consume el evento de SQS
   - Guarda log en DynamoDB (tabla `employee-logs`)
   - Muestra información en consola

### Autenticación (Login)
1. Cliente envía POST a `/auth/login` en Auth Service con email y password
2. Auth Service:
   - Valida que email y password no estén vacíos
   - Busca el usuario por email en DynamoDB (tabla `employees`)
   - Compara el password ingresado con el hash almacenado usando bcrypt
   - Genera un token JWT que contiene el ID del usuario
   - Retorna el token con tiempo de expiración (60 minutos por defecto)

## 🧪 Verificar tablas DynamoDB
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

El **Auth Service** es un microservicio independiente responsable de la autenticación de usuarios. Implementa:

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

El Auth Service sigue arquitectura hexagonal con 3 puertos principales:

1. **UserRepository** (puerto): Interfaz para buscar usuarios
   - **Adaptador**: `DynamoDBUserRepository` - Busca por email en DynamoDB

2. **PasswordHasher** (puerto): Interfaz para comparar passwords
   - **Adaptador**: `BcryptPasswordHasher` - Compara usando bcrypt

3. **TokenGenerator** (puerto): Interfaz para generar/validar JWT
   - **Adaptador**: `JWTTokenGenerator` - Usa golang-jwt/jwt/v5

### Flujo de Autenticación

```
1. Usuario → POST /auth/login {email, password}
2. Auth Service valida que email y password no estén vacíos
3. Busca usuario en DynamoDB por email
4. Compara password con hash almacenado (bcrypt)
5. Si coincide: Genera JWT con user_id
6. Retorna {token, user_id, expires_at}
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
# 1. Crear un usuario
curl -X POST http://localhost:8080/api/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "password": "SecurePass123!"
  }'

# 2. Autenticarse
curl -X POST http://localhost:8082/auth/login \
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
```

### Error en autenticación: "Invalid email or password"
Verifica que:
1. El usuario fue creado correctamente en employee-service
2. El email es exacto (case-sensitive)
3. El password es correcto
4. El auth-service tiene acceso a la misma tabla DynamoDB (employees)

## 📄 Licencia

Este proyecto es de código abierto y está disponible bajo la licencia MIT.
