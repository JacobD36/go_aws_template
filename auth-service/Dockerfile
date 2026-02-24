# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum* ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario compilado
COPY --from=builder /app/auth-service .

# Exponer el puerto
EXPOSE 8082

CMD ["./auth-service"]
