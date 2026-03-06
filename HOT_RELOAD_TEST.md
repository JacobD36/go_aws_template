# 🔥 Hot Reload - Guía de Prueba Rápida

## ✅ Estado Actual

Tu entorno de desarrollo con hot reload está **completamente configurado y funcionando**:

- ✅ Todos los servicios Go corriendo con Air
- ✅ LocalStack configurado con DynamoDB y SQS
- ✅ Hot reload probado y funcionando

## 🚀 Prueba el Hot Reload

### 1. Ver logs en tiempo real

Abre una terminal y ejecuta:

\`\`\`bash
# Ver logs de un servicio específico
make dev-logs-employee
\`\`\`

### 2. Edita un archivo

Abre cualquier archivo `.go`, por ejemplo:
- `employee-service/cmd/main.go`
- `api-gateway/main.go`
- `auth-service/internal/application/auth_service.go`

Haz un pequeño cambio (añade un log, cambia un mensaje, etc.)

### 3. Guarda el archivo

Cuando guardes, verás en los logs:

\`\`\`
cmd/main.go has changed
building...
running...
\`\`\`

¡El servicio se recompila y reinicia automáticamente en 1-2 segundos!

### 4. Prueba el Frontend (Next.js)

El frontend también tiene hot reload con Turbopack:

\`\`\`bash
# En Terminal 1: Ver logs del frontend
make dev-logs-frontend

# En Terminal 2: Abre el navegador
open http://localhost:3000
\`\`\`

Edita cualquier archivo React/TypeScript en `frontend_basic/`:
- `app/page.tsx`
- `components/employees/EmployeeTable.tsx`
- `app/globals.css`

Al guardar, verás en el navegador que la página se actualiza automáticamente (Fast Refresh de Next.js).

## 📝 Comandos Útiles

\`\`\`bash
# Ver estado de contenedores
make dev-ps

# Ver logs de todos los servicios
make dev-logs

# Ver logs de un servicio específico
make dev-logs-employee
make dev-logs-api
make dev-logs-auth
make dev-logs-messaging
make dev-logs-logger

# Reiniciar un servicio específico
make dev-restart-employee

# Reiniciar todo
make dev-restart

# Detener todo
make dev-down

# Limpiar y empezar de cero
make clean
make clean-docker
make dev-build
\`\`\`

## 🧪 Probar los APIs

\`\`\`bash
# Health check
curl http://localhost:8080/health

# Crear empleado
curl -X POST http://localhost:8080/api/employees \\
  -H "Content-Type: application/json" \\
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "position": "Developer",
    "salary": 50000
  }'

# Listar empleados
curl http://localhost:8080/api/employees
\`\`\`

## 📊 Estado de los Servicios

\`\`\`bash
docker ps --format "table {{.Names}}\\t{{.Status}}\\t{{.Ports}}"
\`\`\`

Deberías ver:
- `api-gateway-dev` - Puerto 8080 - ✅ Hot reload activo
- `employee-service-dev` - Puerto 8081 - ✅ Hot reload activo
- `auth-service-dev` - Puerto 8082 - ✅ Hot reload activo
- `logger-service-dev` - ✅ Hot reload activo
- `messaging-service-dev` - ✅ Hot reload activo
- `frontend-basic-dev` - Puerto 3000 - ✅ Hot reload activo (Next.js + Turbopack)
- `localstack` - Puerto 4566

## 🎯 Ejemplo de Flujo de Desarrollo

### Backend (Go + Air)
1. **Terminal 1**: Deja corriendo `make dev-logs-employee`
2. **Terminal 2**: Edita código en tu IDE
3. Al guardar, verás en Terminal 1 cómo se recompila automáticamente
4. **Terminal 3**: Prueba el API con curl

### Frontend (Next.js + Turbopack)
1. **Terminal 1**: Deja corriendo `make dev-logs-frontend` (opcional)
2. **Navegador**: Abre http://localhost:3000
3. **IDE**: Edita archivos en `frontend_basic/`
4. **Navegador**: Los cambios aparecen automáticamente (Fast Refresh)

## 💡 Tips

### Backend (Go)
- Los cambios se detectan en archivos `.go`
- Los archivos `*_test.go` se ignoran
- El directorio `tmp/` se ignora (ahí Air guarda los binarios temporales)
- Si hay un error de compilación, Air lo muestra y espera a que lo corrijas

### Frontend (Next.js)
- Los cambios se detectan en archivos `.tsx`, `.ts`, `.jsx`, `.js`, `.css`
- Next.js usa Fast Refresh para actualizaciones casi instantáneas
- El directorio `.next/` se ignora (archivos de build)
- Los errores se muestran en el navegador y en la terminal
- Turbopack hace que las recargas sean más rápidas que Webpack

## 🐛 Si algo no funciona

\`\`\`bash
# Ver logs de un servicio específico
docker logs employee-service-dev

# Reiniciar un servicio
make dev-restart-employee

# Si todo falla, reinicia desde cero
make clean && make clean-docker
make dev-build
make setup-localstack
\`\`\`

## 📚 Más Información

- Documentación completa: [DEVELOPMENT.md](DEVELOPMENT.md)
- Ver todos los comandos: `make help`
- Air documentation: https://github.com/air-verse/air

---

**¡Disfruta del desarrollo con hot reload! 🎉**
