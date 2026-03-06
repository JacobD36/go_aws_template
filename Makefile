.PHONY: help dev dev-build dev-down dev-logs dev-logs-api dev-logs-employee dev-logs-auth dev-logs-messaging dev-logs-logger dev-logs-frontend dev-restart dev-restart-api dev-restart-employee dev-restart-auth dev-restart-messaging dev-restart-logger prod prod-down setup-localstack clean install-air

help: ## Mostrar esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Comandos de desarrollo (con hot reload)
dev: ## Iniciar todos los servicios en modo desarrollo con hot reload
	docker-compose -f docker-compose.dev.yml up

dev-build: ## Construir y levantar servicios en modo desarrollo
	docker-compose -f docker-compose.dev.yml up --build

dev-down: ## Detener servicios de desarrollo
	docker-compose -f docker-compose.dev.yml down

dev-logs: ## Ver logs de todos los servicios en desarrollo
	docker-compose -f docker-compose.dev.yml logs -f

dev-logs-api: ## Ver logs del API Gateway
	docker-compose -f docker-compose.dev.yml logs -f api-gateway

dev-logs-employee: ## Ver logs del Employee Service
	docker-compose -f docker-compose.dev.yml logs -f employee-service

dev-logs-auth: ## Ver logs del Auth Service
	docker-compose -f docker-compose.dev.yml logs -f auth-service

dev-logs-messaging: ## Ver logs del Messaging Service
	docker-compose -f docker-compose.dev.yml logs -f messaging-service

dev-logs-logger: ## Ver logs del Logger Service
	docker-compose -f docker-compose.dev.yml logs -f logger-service

dev-logs-frontend: ## Ver logs del Frontend
	docker-compose -f docker-compose.dev.yml logs -f frontend-basic

dev-logs-localstack: ## Ver logs de LocalStack
	docker-compose -f docker-compose.dev.yml logs -f localstack

dev-restart: ## Reiniciar todos los servicios
	docker-compose -f docker-compose.dev.yml restart

dev-restart-api: ## Reiniciar API Gateway
	docker-compose -f docker-compose.dev.yml restart api-gateway

dev-restart-employee: ## Reiniciar Employee Service
	docker-compose -f docker-compose.dev.yml restart employee-service

dev-restart-auth: ## Reiniciar Auth Service
	docker-compose -f docker-compose.dev.yml restart auth-service

dev-restart-messaging: ## Reiniciar Messaging Service
	docker-compose -f docker-compose.dev.yml restart messaging-service

dev-restart-logger: ## Reiniciar Logger Service
	docker-compose -f docker-compose.dev.yml restart logger-service

dev-restart-frontend: ## Reiniciar Frontend
	docker-compose -f docker-compose.dev.yml restart frontend-basic

dev-ps: ## Ver estado de los contenedores
	docker-compose -f docker-compose.dev.yml ps

# Comandos de producción
prod: ## Iniciar servicios en modo producción
	docker-compose up

prod-build: ## Construir y levantar servicios en modo producción
	docker-compose up --build

prod-down: ## Detener servicios de producción
	docker-compose down

# Configuración de LocalStack
setup-localstack: ## Configurar recursos de AWS en LocalStack
	@echo "Esperando a que LocalStack esté listo..."
	@sleep 10
	@echo "Ejecutando script de inicialización..."
	@bash ./setup-aws-resources.sh

# Limpieza
clean: ## Limpiar archivos temporales y datos de LocalStack
	@echo "Limpiando archivos temporales..."
	@find . -type d -name "tmp" -exec rm -rf {} + 2>/dev/null || true
	@rm -rf localstack_data/
	@echo "Limpieza completada"

clean-docker: ## Limpiar contenedores, imágenes y volúmenes no usados
	docker-compose -f docker-compose.dev.yml down -v
	docker system prune -f

# Instalación
install-air: ## Instalar Air para desarrollo local (opcional)
	go install github.com/cosmtrek/air@latest

# Testing
test: ## Ejecutar tests (agregar comandos de test según necesites)
	@echo "Ejecutando tests..."
	@cd api-gateway && go test ./... || true
	@cd auth-service && go test ./... || true
	@cd employee-service && go test ./... || true
	@cd logger-service && go test ./... || true
	@cd messaging-service && go test ./... || true
