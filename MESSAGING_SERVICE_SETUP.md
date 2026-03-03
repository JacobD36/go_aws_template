# Guía de Configuración - Messaging Service

## Nuevas Funcionalidades Añadidas

### 1. Messaging Service
Nuevo microservicio que:
- Consume eventos de creación de usuarios desde `user-queue`
- Simula el envío de correos electrónicos de bienvenida
- Persiste los mensajes enviados en DynamoDB (`messages`)
- Publica eventos de log para el `logger-service`

### 2. Registro de Usuarios
El `auth-service` ahora incluye:
- Endpoint `/auth/register` para registrar nuevos usuarios
- Publicación automática de eventos `user.created` a `user-queue`
- Integración completa con el flujo de eventos

### 3. API Gateway Actualizado
Nuevo endpoint disponible:
- `POST /api/auth/register` - Registro de nuevos usuarios

## Arquitectura de Eventos

```
User Register (auth-service)
    ↓
  user-queue (SQS)
    ↓
messaging-service
    ↓ (simula envío email)
    ↓
  employee-queue (SQS)
    ↓
logger-service
```

## Recursos AWS Añadidos

### Nuevas Colas SQS:
- `user-queue`: Para eventos de usuario

### Nuevas Tablas DynamoDB:
- `messages`: Almacena los mensajes enviados

## Ejecutar el Sistema

### 1. Iniciar LocalStack y servicios
```bash
# Eliminar contenedores anteriores (si existen)
docker-compose down -v

# Construir e iniciar todos los servicios
docker-compose up --build
```

### 2. Inicializar recursos AWS
```bash
# En otra terminal
./setup-aws-resources.sh
```

## Probar el Nuevo Flujo

### 1. Registrar un nuevo usuario
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "password": "mipassword123"
  }'
```

**Respuesta esperada:**
```json
{
  "user_id": "uuid-generado",
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "created_at": "2026-03-02T10:30:00Z"
}
```

### 2. Verificar logs del messaging-service
```bash
docker logs messaging-service
```

Deberías ver:
```
=== SIMULACIÓN DE ENVÍO DE EMAIL ===
Para: juan@example.com
Asunto: ¡Bienvenido a nuestro sistema!
Cuerpo:
Hola Juan Pérez,

¡Bienvenido a nuestro sistema! Estamos encantados de tenerte con nosotros.
...
===================================
✓ Email enviado exitosamente a juan@example.com
```

### 3. Verificar logs del logger-service
```bash
docker logs logger-service
```

Deberías ver el evento `message.sent` procesado.

### 4. Login con el usuario creado
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan@example.com",
    "password": "mipassword123"
  }'
```

**Respuesta esperada:**
```json
{
  "token": "jwt-token-generado",
  "user_id": "uuid-del-usuario",
  "expires_at": 1709380200
}
```

## Verificar Datos en DynamoDB

### Ver mensajes enviados
```bash
aws --endpoint-url=http://localhost:4566 dynamodb scan \
  --table-name messages \
  --region us-east-1
```

### Ver usuarios registrados
```bash
aws --endpoint-url=http://localhost:4566 dynamodb scan \
  --table-name employees \
  --region us-east-1
```

### Ver logs de eventos
```bash
aws --endpoint-url=http://localhost:4566 dynamodb scan \
  --table-name employee-logs \
  --region us-east-1
```

## Estructura de Servicios

```
├── auth-service (Puerto 8082)
│   ├── POST /auth/register  (NUEVO)
│   └── POST /auth/login
│
├── messaging-service (sin puerto HTTP - solo eventos)
│   ├── Consume: user-queue
│   └── Publica: employee-queue
│
├── logger-service (sin puerto HTTP - solo eventos)
│   └── Consume: employee-queue
│
├── employee-service (Puerto 8081)
│   ├── POST /employees
│   ├── GET /employees
│   └── Publica: employee-queue
│
└── api-gateway (Puerto 8080)
    ├── POST /api/auth/register  (NUEVO)
    ├── POST /api/auth/login
    ├── POST /api/employees
    └── GET /api/employees
```

## Principios Aplicados

### Arquitectura Hexagonal
Todos los servicios siguen el patrón de puertos y adaptadores:
- **Domain**: Lógica de negocio pura
- **Ports**: Interfaces (contratos)
- **Infrastructure**: Implementaciones concretas (AWS, simuladores)
- **Application**: Orquestación y casos de uso

### SOLID
- **S**: Cada clase/función tiene una única responsabilidad
- **O**: Abierto para extensión, cerrado para modificación
- **L**: Los adaptadores son intercambiables
- **I**: Interfaces específicas y segregadas
- **D**: Dependencias mediante abstracciones (interfaces)

### Event-Driven Architecture
- Comunicación asíncrona mediante eventos
- Bajo acoplamiento entre servicios
- Alta escalabilidad y resiliencia

## Troubleshooting

### Los servicios no se comunican
```bash
# Verificar que LocalStack está funcionando
curl http://localhost:4566/_localstack/health

# Verificar que las colas existen
aws --endpoint-url=http://localhost:4566 sqs list-queues --region us-east-1
```

### El messaging-service no procesa eventos
```bash
# Verificar logs del servicio
docker logs messaging-service -f

# Verificar mensajes en la cola
aws --endpoint-url=http://localhost:4566 sqs receive-message \
  --queue-url http://localhost:4566/000000000000/user-queue \
  --region us-east-1
```

### Reiniciar todo desde cero
```bash
# Detener y eliminar todo
docker-compose down -v

# Reconstruir e iniciar
docker-compose up --build

# En otra terminal, inicializar recursos
./setup-aws-resources.sh
```

## Próximos Pasos

Para extender el sistema puedes:
1. Añadir implementación real de envío de emails (SendGrid, AWS SES)
2. Añadir envío de SMS (Twilio, AWS SNS)
3. Implementar templates de mensajes
4. Añadir retry logic y dead letter queues
5. Implementar notificaciones push
