# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar go.mod y descargar dependencias
COPY go.mod ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o messaging-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario desde el build stage
COPY --from=builder /app/messaging-service .

# Exponer puerto (si en el futuro se añade API HTTP)
# EXPOSE 8083

CMD ["./messaging-service"]
