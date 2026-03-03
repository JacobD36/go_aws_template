# Messaging Service

Microservicio de mensajería basado en arquitectura hexagonal que se encarga de enviar notificaciones (correos electrónicos, SMS, etc.) en respuesta a eventos del sistema.

## Arquitectura

El servicio implementa **Arquitectura Hexagonal (Puertos y Adaptadores)** con las siguientes capas:

### Dominio (`internal/domain`)
- **user_event.go**: Eventos de usuario recibidos
- **message.go**: Entidad Message y lógica de negocio de mensajes
- **log_event.go**: Eventos de log para publicar
- **errors.go**: Errores del dominio

### Puertos (`internal/ports`)
Interfaces que definen contratos:
- **event_consumer.go**: Consumidor de eventos (SQS)
- **event_publisher.go**: Publicador de eventos (SQS)
- **message_sender.go**: Envío de mensajes (Email, SMS)
- **repository.go**: Persistencia de mensajes (DynamoDB)

### Infraestructura (`internal/infrastructure`)
Implementaciones concretas:
- **sqs_consumer.go**: Consume eventos de usuario desde SQS
- **sqs_publisher.go**: Publica eventos de log a SQS
- **simulated_sender.go**: Simulador de envío de mensajes (Email/SMS)
- **dynamodb_repository.go**: Persistencia en DynamoDB

### Aplicación (`internal/application`)
- **messaging_service.go**: Lógica de negocio y orquestación

## Flujo de Eventos

1. **Auth Service** crea un nuevo usuario → publica evento `user.created` a `user-queue`
2. **Messaging Service** consume el evento de `user-queue`
3. Crea y envía mensaje de bienvenida (simulado)
4. Guarda el mensaje en DynamoDB tabla `messages`
5. Publica evento `message.sent` a `employee-queue` (para logging)
6. **Logger Service** consume y registra el evento

## Principios SOLID

- **S (Single Responsibility)**: Cada componente tiene una única responsabilidad
- **O (Open/Closed)**: Fácil extender con nuevos tipos de mensajes sin modificar código existente
- **L (Liskov Substitution)**: Los adaptadores son intercambiables
- **I (Interface Segregation)**: Interfaces pequeñas y específicas
- **D (Dependency Inversion)**: Dependencias a través de interfaces (puertos)

## Patrones de Diseño

- **Ports & Adapters (Hexagonal Architecture)**
- **Dependency Injection**: Inyección de dependencias en el constructor
- **Strategy Pattern**: Diferentes estrategias de envío (Email, SMS)
- **Event-Driven Architecture**: Comunicación asíncrona basada en eventos

## Variables de Entorno

- `AWS_REGION`: Región de AWS (default: us-east-1)
- `AWS_ENDPOINT`: Endpoint de LocalStack para desarrollo local
- `AWS_ACCESS_KEY_ID`: Credenciales AWS
- `AWS_SECRET_ACCESS_KEY`: Credenciales AWS
- `USER_QUEUE_URL`: URL de la cola SQS para eventos de usuario
- `LOG_QUEUE_URL`: URL de la cola SQS para eventos de log
- `DYNAMODB_TABLE`: Tabla de DynamoDB para mensajes (default: messages)

## Testing

El servicio está diseñado para ser fácilmente testeable:
- Todas las dependencias son inyectadas a través de interfaces
- Se pueden crear mocks de los puertos para testing unitario
- La lógica de negocio está aislada en la capa de aplicación

## Ejemplo de Uso

Cuando se registra un nuevo usuario mediante el endpoint `/auth/register` del `auth-service`, automáticamente:

1. Se envía un correo de bienvenida (simulado)
2. Se registra el mensaje en DynamoDB
3. Se notifica al logger-service del envío

No se requiere ninguna acción manual adicional.
