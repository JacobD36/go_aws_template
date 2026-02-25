# Frontend Basic - Portal Administrativo

Este es el frontend del portal administrativo construido con Next.js, React, TypeScript y Tailwind CSS.

## Tecnologías

- **Next.js 16**: Framework de React con App Router
- **React 19**: Biblioteca de interfaz de usuario
- **TypeScript**: Tipado estático
- **Tailwind CSS 4**: Framework de estilos utility-first

## Características

- 🔐 **Autenticación**: Login con JWT
- 👥 **Gestión de Empleados**: Lista, creación con validación
- 🎨 **UI Moderna**: Componentes reutilizables con Tailwind
- 📱 **Responsive**: Diseño adaptable a diferentes dispositivos
- ✨ **Clean Code**: Arquitectura limpia y buenas prácticas

## Estructura del Proyecto

```
frontend_basic/
├── app/                      # App Router de Next.js
│   ├── login/               # Página de login
│   ├── employees/           # Módulo de empleados
│   │   ├── page.tsx        # Lista de empleados
│   │   └── layout.tsx      # Layout del dashboard
│   ├── layout.tsx          # Layout principal
│   └── page.tsx            # Página de inicio (redirect)
├── components/              # Componentes React
│   ├── ui/                 # Componentes UI base
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   ├── Modal.tsx
│   │   └── Table.tsx
│   ├── layout/             # Componentes de layout
│   │   ├── Sidebar.tsx
│   │   └── Header.tsx
│   └── employees/          # Componentes específicos de empleados
│       ├── EmployeeModal.tsx
│       └── EmployeeTable.tsx
├── lib/                    # Utilidades y servicios
│   ├── api.ts             # Cliente API
│   ├── auth.ts            # Utilidades de autenticación
│   └── constants.ts       # Constantes de la aplicación
└── types/                 # Tipos TypeScript
    └── index.ts

```

## Configuración

### Variables de Entorno

Crea un archivo `.env.local` con:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

### Instalación

```bash
# Instalar dependencias
npm install

# Modo desarrollo
npm run dev

# Build de producción
npm run build

# Iniciar producción
npm start
```

## Docker

### Build

```bash
docker build -t frontend-basic .
```

### Run

```bash
docker run -p 3000:3000 -e NEXT_PUBLIC_API_URL=http://localhost:8080/api frontend-basic
```

## Características Principales

### Autenticación

- Login con email y password
- Almacenamiento seguro de token JWT
- Validación de sesión
- Redirección automática según estado de autenticación

### Gestión de Empleados

- **Listar**: Tabla con paginación y estado vacío
- **Crear**: Modal con formulario validado
  - Validación de email
  - Validación de contraseña (8 caracteres, mayúscula, número, carácter especial)
- **Acciones**: Botones de editar y eliminar (UI only)

### Componentes UI

Todos los componentes siguen un diseño consistente con Tailwind CSS:

- **Button**: Variantes (primary, secondary, danger, ghost), tamaños, loading state
- **Input**: Con label, error messages, validación visual
- **Modal**: Backdrop, ESC para cerrar, animaciones
- **Table**: Responsive, estado vacío, columnas personalizables

## API Integration

El frontend se comunica con el API Gateway en:

- `POST /api/auth/login`: Autenticación
- `GET /api/employees`: Obtener empleados
- `POST /api/employees`: Crear empleado

## Mejores Prácticas

- ✅ Componentes reutilizables y modulares
- ✅ Tipado fuerte con TypeScript
- ✅ Manejo de errores centralizado
- ✅ Validación de formularios
- ✅ Estados de carga y error
- ✅ Clean Architecture
- ✅ Código autodocumentado

## Próximas Funcionalidades

- [ ] Editar empleado
- [ ] Eliminar empleado
- [ ] Búsqueda y filtros
- [ ] Paginación
- [ ] Cambio de contraseña
- [ ] Perfil de usuario
