# Guía de Desarrollo con Hot Reload

Esta guía explica cómo trabajar con el entorno de desarrollo que incluye hot reload automático para todos los servicios Go.

## 🚀 Inicio Rápido

### Iniciar el entorno de desarrollo

```bash
make dev-build
```

Esto levantará todos los servicios con hot reload activado. Cualquier cambio en el código Go se detectará automáticamente y el servicio se reconstruirá y reiniciará.

### Detener el entorno

```bash
make dev-down
```

## 📋 Comandos Disponibles

### Desarrollo

| Comando | Descripción |
|---------|-------------|
| `make dev` | Iniciar todos los servicios en modo desarrollo |
| `make dev-build` | Construir y levantar servicios en modo desarrollo |
| `make dev-down` | Detener servicios de desarrollo |
| `make dev-ps` | Ver estado de los contenedores |

### Logs

| Comando | Descripción |
|---------|-------------|
| `make dev-logs` | Ver logs de todos los servicios |
| `make dev-logs-api` | Ver logs del API Gateway |
| `make dev-logs-employee` | Ver logs del Employee Service |
| `make dev-logs-auth` | Ver logs del Auth Service |
| `make dev-logs-messaging` | Ver logs del Messaging Service |
| `make dev-logs-logger` | Ver logs del Logger Service |
| `make dev-logs-frontend` | Ver logs del Frontend |
| `make dev-logs-localstack` | Ver logs de LocalStack |

### Reiniciar Servicios

| Comando | Descripción |
|---------|-------------|
| `make dev-restart` | Reiniciar todos los servicios |
| `make dev-restart-api` | Reiniciar API Gateway |
| `make dev-restart-employee` | Reiniciar Employee Service |
| `make dev-restart-auth` | Reiniciar Auth Service |
| `make dev-restart-messaging` | Reiniciar Messaging Service |
| `make dev-restart-logger` | Reiniciar Logger Service |

### LocalStack

```bash
make setup-localstack
```

Configura todos los recursos de AWS necesarios en LocalStack (tablas DynamoDB, colas SQS, etc.)

### Limpieza

| Comando | Descripción |
|---------|-------------|
| `make clean` | Limpiar archivos temporales y datos de LocalStack |
| `make clean-docker` | Limpiar contenedores, imágenes y volúmenes no usados |

### Producción

| Comando | Descripción |
|---------|-------------|
| `make prod` | Iniciar servicios en modo producción |
| `make prod-build` | Construir y levantar servicios en modo producción |
| `make prod-down` | Detener servicios de producción |

## 🔥 Hot Reload

### ¿Cómo funciona?

Cada servicio Go utiliza [Air](https://github.com/cosmtrek/air) para detectar cambios en el código fuente y recompilar automáticamente.

**Configuración:**
- Cada servicio tiene un archivo `.air.toml` con su configuración
- Los cambios se detectan en archivos `.go`
- El servicio se recompila y reinicia automáticamente
- Los logs de compilación se muestran en tiempo real

### ¿Qué archivos vigila?

**Backend (Air):**

Air vigila todos los archivos con estas extensiones:
- `.go`
- `.tpl`
- `.tmpl`
- `.html`

**Excluye:**
- Archivos de test (`*_test.go`)
- Directorio `tmp/`
- Directorio `vendor/`
- Directorio `testdata/`

**Frontend (Next.js):**

Next.js vigila automáticamente:
- `.tsx`, `.ts` (TypeScript/React)
- `.jsx`, `.js` (JavaScript/React)
- `.css`, `.scss` (Estilos)
- Archivos en `app/`, `components/`, `lib/`, etc.

**Excluye:**
- Directorio `.next/`
- Directorio `node_modules/`

### Ejemplo de flujo de trabajo

1. Levanta el entorno:
   ```bash
   make dev-build
   ```

2. En otra terminal, observa los logs del servicio que modificarás:
   ```bash
   make dev-logs-employee
   ```

3. Edita el código en `employee-service/internal/...`

4. Guarda el archivo - automáticamente verás:
   ```
   [employee-service-dev] building...
   [employee-service-dev] running...
   ```

5. ¡El servicio está actualizado! Prueba los cambios inmediatamente.

## 🛠️ Configuración por Servicio

### API Gateway
- **Puerto:** 8080
- **Comando hot reload:** `air -c .air.toml`
- **Path main:** `main.go`

### Auth Service
- **Puerto:** 8082
- **Comando hot reload:** `air -c .air.toml`
- **Path main:** `cmd/main.go`

### Employee Service
- **Puerto:** 8081
- **Comando hot reload:** `air -c .air.toml`
- **Path main:** `cmd/main.go`

### Logger Service
- **Comando hot reload:** `air -c .air.toml`
- **Path main:** `cmd/main.go`

### Messaging Service
- **Comando hot reload:** `air -c .air.toml`
- **Path main:** `cmd/main.go`

### Frontend (Next.js)
- **Puerto:** 3000
- **Comando hot reload:** `npm run dev`
- **Hot reload:** Turbopack + Fast Refresh
- **Framework:** Next.js 16 con React 19

## 📝 Volúmenes

Los Dockerfiles de desarrollo montan volúmenes para permitir el hot reload:

```yaml
volumes:
  - ./service-name:/app      # Código fuente
  - /app/tmp                 # Excluir directorio tmp de Air
```

Esto significa que:
- ✅ Los cambios en tu código local se reflejan instantáneamente en el contenedor
- ✅ No necesitas reconstruir las imágenes para cada cambio
- ✅ Los archivos temporales de Air no contaminan tu directorio local

## 🐛 Troubleshooting

### El hot reload no está funcionando

1. Verifica que el servicio esté usando el Dockerfile correcto:
   ```bash
   make dev-ps
   ```

2. Revisa los logs del servicio:
   ```bash
   make dev-logs-employee  # o el servicio que sea
   ```

3. Si ves errores de permisos, limpia y reconstruye:
   ```bash
   make clean
   make dev-build
   ```

### Los cambios no se detectan

1. Asegúrate de que el archivo tiene extensión `.go`
2. Verifica que el archivo no esté en un directorio excluido (`tmp/`, `vendor/`, etc.)
3. Reinicia el servicio específico:
   ```bash
   make dev-restart-employee
   ```

### LocalStack no está listo

Si los servicios fallan al conectarse a LocalStack:

```bash
# Verifica el health de LocalStack
docker-compose -f docker-compose.dev.yml logs localstack

# Espera a que esté listo y configura recursos
make setup-localstack
```

### Errores de compilación persistentes

Si Air muestra errores de compilación:

1. Revisa `build-errors.log` en el directorio del servicio
2. Ejecuta `go build` localmente para verificar el error
3. Corrige el código y guarda - Air detectará el cambio

## 💡 Tips

### Ver logs en tiempo real mientras desarrollas

Abre varias terminales:

```bash
# Terminal 1: Servicios corriendo
make dev

# Terminal 2: Logs del servicio que estás editando
make dev-logs-employee

# Terminal 3: Logs de LocalStack (si trabajas con AWS)
make dev-logs-localstack
```

### Reinicio rápido de un servicio específico

En lugar de reiniciar todo:

```bash
make dev-restart-employee
```

### Limpiar y empezar desde cero

```bash
make clean
make clean-docker
make dev-build
make setup-localstack
```

## 🔄 Diferencias entre Desarrollo y Producción

| Aspecto | Desarrollo | Producción |
|---------|-----------|------------|
| **Dockerfile** | `Dockerfile.dev` | `Dockerfile` |
| **Hot Reload** | ✅ Activado | ❌ Desactivado |
| **Volúmenes** | ✅ Código montado | ❌ Código copiado |
| **Build** | Incremental | Multi-stage optimizado |
| **Tamaño imagen** | ~500MB | ~20MB |
| **Air instalado** | ✅ Sí | ❌ No |

## 📚 Recursos

- [Air - Live reload for Go apps](https://github.com/cosmtrek/air)
- [LocalStack Documentation](https://docs.localstack.cloud/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

## 🎯 Siguiente Pasos

1. ✅ Configuración completada
2. Personaliza `.air.toml` según tus necesidades
3. Agrega tests automáticos
4. Configura CI/CD

---

¿Necesitas ayuda? Revisa el README principal o los logs del servicio específico.
