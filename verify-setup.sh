#!/bin/bash

# Script de verificación del entorno de desarrollo
# Verifica que todos los archivos necesarios para hot reload estén presentes

echo "🔍 Verificando configuración de Hot Reload..."
echo ""

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Contador de verificaciones
CHECKS_PASSED=0
CHECKS_FAILED=0

# Función para verificar archivo
check_file() {
    local file=$1
    local description=$2
    
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓${NC} $description: $file"
        ((CHECKS_PASSED++))
    else
        echo -e "${RED}✗${NC} $description: $file ${RED}(NO ENCONTRADO)${NC}"
        ((CHECKS_FAILED++))
    fi
}

# Función para verificar directorio
check_dir() {
    local dir=$1
    local description=$2
    
    if [ -d "$dir" ]; then
        echo -e "${GREEN}✓${NC} $description: $dir"
        ((CHECKS_PASSED++))
    else
        echo -e "${YELLOW}⚠${NC} $description: $dir ${YELLOW}(NO EXISTE - se creará automáticamente)${NC}"
    fi
}

echo "📁 Verificando archivos de configuración..."
echo ""

# Verificar archivos principales
check_file "docker-compose.dev.yml" "Docker Compose Development"
check_file "Makefile" "Makefile"
check_file "DEVELOPMENT.md" "Documentación de Desarrollo"
check_file ".gitignore" "GitIgnore"

echo ""
echo "🔧 Verificando servicios Go..."
echo ""

# Lista de servicios
SERVICES=("api-gateway" "auth-service" "employee-service" "logger-service" "messaging-service")

for service in "${SERVICES[@]}"; do
    echo "  📦 $service:"
    check_file "$service/.air.toml" "    Air Config"
    check_file "$service/Dockerfile.dev" "    Dockerfile Dev"
    check_file "$service/go.mod" "    Go Module"
    check_dir "$service/tmp" "    Tmp Directory"
    echo ""
done

echo "  📦 frontend_basic:"
check_file "frontend_basic/Dockerfile.dev" "    Dockerfile Dev"
check_file "frontend_basic/package.json" "    Package.json"
echo ""

echo "🐳 Verificando Docker..."
echo ""

# Verificar Docker
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓${NC} Docker instalado: $(docker --version)"
    ((CHECKS_PASSED++))
else
    echo -e "${RED}✗${NC} Docker NO instalado"
    ((CHECKS_FAILED++))
fi

# Verificar Docker Compose
if command -v docker-compose &> /dev/null; then
    echo -e "${GREEN}✓${NC} Docker Compose instalado: $(docker-compose --version)"
    ((CHECKS_PASSED++))
else
    echo -e "${RED}✗${NC} Docker Compose NO instalado"
    ((CHECKS_FAILED++))
fi

echo ""
echo "📊 Resumen"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✓ Verificaciones exitosas: $CHECKS_PASSED${NC}"

if [ $CHECKS_FAILED -gt 0 ]; then
    echo -e "${RED}✗ Verificaciones fallidas: $CHECKS_FAILED${NC}"
    echo ""
    echo -e "${RED}⚠ Algunos archivos necesarios no se encontraron.${NC}"
    echo "Por favor, asegúrate de ejecutar este script desde el directorio raíz del proyecto."
    exit 1
else
    echo ""
    echo -e "${GREEN}✅ ¡Todo listo para desarrollo con hot reload!${NC}"
    echo ""
    echo "Para comenzar, ejecuta:"
    echo -e "  ${YELLOW}make dev-build${NC}"
    echo ""
    echo "Para ver todos los comandos disponibles:"
    echo -e "  ${YELLOW}make help${NC}"
    echo ""
fi
