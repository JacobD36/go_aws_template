# Guía Completa del Frontend - frontend_basic

## 📋 Tabla de Contenidos
1. [¿Qué es Next.js y para qué sirve?](#qué-es-nextjs-y-para-qué-sirve)
2. [Estructura de Carpetas](#estructura-de-carpetas)
3. [Punto de Entrada - Cómo Inicia Todo](#punto-de-entrada---cómo-inicia-todo)
4. [Flujo de la Aplicación](#flujo-de-la-aplicación)
5. [Función de api.ts](#función-de-apits)
6. [Componentes Clave](#componentes-clave)
7. [Gestión de Autenticación](#gestión-de-autenticación)
8. [Cómo Iniciar desde Cero](#cómo-iniciar-desde-cero)

---

## 🚀 ¿Qué es Next.js y para qué sirve?

**Next.js** es un framework de React que facilita la creación de aplicaciones web modernas. Piensa en él como una **"capa mejorada sobre React"** que te da superpoderes adicionales:

### ¿Por qué usamos Next.js en lugar de React puro?

| Característica | React Puro | Next.js |
|---------------|------------|---------|
| **Enrutamiento** | Necesitas React Router (librería adicional) | ✅ Enrutamiento automático basado en carpetas |
| **SEO** | Difícil (todo renderiza en el cliente) | ✅ Server-Side Rendering y Static Generation |
| **Optimización** | Manual | ✅ Optimización automática de imágenes, fuentes, etc. |
| **Configuración** | Compleja (Webpack, Babel, etc.) | ✅ Configuración mínima - funciona de inmediato |
| **API Routes** | Necesitas servidor separado | ✅ Puedes crear APIs dentro del proyecto |

### Conceptos Clave de Next.js:

1. **App Router (Carpeta `app/`)**: El sistema de rutas basado en el sistema de archivos
2. **Server Components**: Componentes que se ejecutan en el servidor (por defecto)
3. **Client Components**: Componentes interactivos (con `'use client'`)
4. **Layouts**: Plantillas que envuelven múltiples páginas
5. **Standalone Mode**: Para producción en Docker (genera todo lo necesario en un solo directorio)

---

## 📁 Estructura de Carpetas

```
frontend_basic/
├── app/                          # ⭐ NÚCLEO - Sistema de rutas de Next.js
│   ├── layout.tsx               # Layout raíz (envuelve toda la aplicación)
│   ├── page.tsx                 # Página principal (/)
│   ├── globals.css              # Estilos globales
│   ├── login/                   # Ruta /login
│   │   └── page.tsx            # Página de login
│   └── employees/               # Ruta /employees
│       ├── layout.tsx          # Layout específico para empleados
│       └── page.tsx            # Página de empleados
│
├── components/                   # 🧩 Componentes reutilizables
│   ├── ui/                      # Componentes de UI básicos
│   │   ├── Button.tsx          # Botón reutilizable
│   │   ├── Input.tsx           # Input reutilizable
│   │   ├── Modal.tsx           # Modal reutilizable
│   │   └── Table.tsx           # Tabla reutilizable
│   ├── layout/                  # Componentes de layout
│   │   ├── Header.tsx          # Barra superior
│   │   └── Sidebar.tsx         # Menú lateral
│   └── employees/               # Componentes específicos de empleados
│       ├── EmployeeTable.tsx   # Tabla de empleados
│       └── EmployeeModal.tsx   # Modal para crear empleados
│
├── lib/                          # 📚 Lógica de negocio y utilidades
│   ├── api.ts                   # Cliente API (todas las peticiones HTTP)
│   ├── auth.ts                  # Utilidades de autenticación
│   └── constants.ts             # Constantes de la aplicación
│
├── types/                        # 🏷️ Definiciones de TypeScript
│   └── index.ts                 # Todos los tipos e interfaces
│
├── public/                       # 🖼️ Archivos estáticos (imágenes, etc.)
│
├── next.config.ts               # ⚙️ Configuración de Next.js
├── tsconfig.json                # ⚙️ Configuración de TypeScript
├── tailwind.config.ts           # ⚙️ Configuración de Tailwind CSS
├── package.json                 # 📦 Dependencias y scripts
└── Dockerfile                   # 🐳 Configuración de Docker
```

---

## 🎬 Punto de Entrada - Cómo Inicia Todo

### El Flujo de Inicio en Next.js:

```
1. Usuario visita http://localhost:3000
        ↓
2. Next.js busca app/layout.tsx (Layout raíz)
        ↓
3. Next.js busca app/page.tsx (Página raíz)
        ↓
4. Se renderiza la aplicación
```

### Archivos que Inician la Aplicación:

#### 1. `app/layout.tsx` - El Layout Raíz (Primera Pieza)

```tsx
// Este es el "esqueleto" de TODA la aplicación
export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        {children}  {/* Aquí se inyectan las páginas */}
      </body>
    </html>
  );
}
```

**¿Qué hace?**
- Define la estructura HTML básica (`<html>`, `<body>`)
- Carga fuentes globales
- Aplica estilos globales
- **Envuelve todas las páginas** - Este layout está presente en TODAS las rutas

#### 2. `app/page.tsx` - La Página Principal (/)

```tsx
// Ruta: http://localhost:3000/
export default function Home() {
  redirect('/login');  // Redirige inmediatamente al login
}
```

**¿Qué hace?**
- Es la primera página que carga cuando visitas `/`
- En nuestro caso, redirige inmediatamente a `/login`
- Es la "puerta de entrada" de la aplicación

#### 3. `app/login/page.tsx` - Página de Login (/login)

```tsx
'use client';  // ← Esto lo convierte en Client Component (interactivo)

export default function LoginPage() {
  // Lógica de login...
}
```

**¿Qué hace?**
- Maneja el formulario de autenticación
- Valida credenciales
- Redirige a `/employees` si el login es exitoso

---

## 🔄 Flujo de la Aplicación

### Flujo Completo de un Usuario:

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Usuario visita http://localhost:3000                     │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Next.js carga:                                            │
│    - app/layout.tsx (Layout raíz con HTML base)             │
│    - app/page.tsx (Redirige a /login)                       │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Usuario en /login (app/login/page.tsx)                   │
│    - Muestra formulario de login                             │
│    - Usuario ingresa email y contraseña                      │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Usuario hace clic en "Iniciar Sesión"                    │
│    ┌────────────────────────────────────────────────────┐   │
│    │ handleSubmit() en LoginPage                        │   │
│    │  ↓                                                 │   │
│    │ apiClient.login(credentials) en lib/api.ts        │   │
│    │  ↓                                                 │   │
│    │ fetch('http://localhost:8080/api/auth/login')     │   │
│    └────────────────────────────────────────────────────┘   │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. API Gateway responde con token                           │
│    {                                                          │
│      "token": "eyJhbGc...",                                  │
│      "user_id": "123",                                       │
│      "expires_at": 1234567890                                │
│    }                                                          │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. saveAuthData() guarda el token en localStorage           │
│    - localStorage.setItem('auth_token', token)               │
│    - localStorage.setItem('user_id', user_id)                │
│    - localStorage.setItem('expires_at', expires_at)          │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│ 7. Redirige a /employees (router.push('/employees'))        │
│    - Next.js carga app/employees/layout.tsx                 │
│    - Next.js carga app/employees/page.tsx                   │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│ 8. En /employees:                                            │
│    - Se muestra Header + Sidebar (layout)                   │
│    - Se cargan empleados (apiClient.getEmployees())         │
│    - Se muestra tabla de empleados                          │
└─────────────────────────────────────────────────────────────┘
```

---

## 🌐 Función de api.ts

`lib/api.ts` es el **cliente centralizado de API** - un patrón de diseño que centraliza todas las comunicaciones HTTP con el backend.

### ¿Por qué usar un cliente API centralizado?

**❌ Sin api.ts (Código duplicado):**
```tsx
// En LoginPage.tsx
fetch('http://localhost:8080/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(data)
})

// En EmployeesPage.tsx
fetch('http://localhost:8080/api/employees', {
  method: 'GET',
  headers: { 
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}` 
  }
})
```

**✅ Con api.ts (Centralizado):**
```tsx
// En LoginPage.tsx
apiClient.login(credentials)

// En EmployeesPage.tsx
apiClient.getEmployees()
```

### Estructura de api.ts:

```typescript
class ApiClient {
  private baseUrl: string;  // URL base del API

  // 🔒 Método privado: Maneja errores de respuesta
  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      throw error;  // Lanza error si la respuesta no es exitosa
    }
    return response.json();  // Parsea y retorna JSON
  }

  // 🔒 Método privado: Construye headers con autenticación
  private getHeaders(includeAuth: boolean = false): HeadersInit {
    const headers = { 'Content-Type': 'application/json' };
    
    if (includeAuth) {
      const token = getAuthToken();  // Obtiene token de localStorage
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }
    
    return headers;
  }

  // 🌍 Método público: Login
  async login(credentials: LoginCredentials): Promise<LoginResponse> {
    const response = await fetch(`${this.baseUrl}/auth/login`, {
      method: 'POST',
      headers: this.getHeaders(),  // Sin autenticación
      body: JSON.stringify(credentials),
    });
    return this.handleResponse<LoginResponse>(response);
  }

  // 🌍 Método público: Obtener empleados
  async getEmployees(): Promise<Employee[]> {
    const response = await fetch(`${this.baseUrl}/employees`, {
      method: 'GET',
      headers: this.getHeaders(true),  // CON autenticación
    });
    return this.handleResponse<Employee[]>(response);
  }

  // 🌍 Método público: Crear empleado
  async createEmployee(data: CreateEmployeeRequest): Promise<Employee> {
    const response = await fetch(`${this.baseUrl}/employees`, {
      method: 'POST',
      headers: this.getHeaders(true),  // CON autenticación
      body: JSON.stringify(data),
    });
    return this.handleResponse<Employee>(response);
  }
}

// Exporta una instancia única (Singleton)
export const apiClient = new ApiClient(API_BASE_URL);
```

### Ventajas del Cliente API:

1. **✅ DRY (Don't Repeat Yourself)**: No duplicamos código de fetch
2. **✅ Manejo centralizado de errores**: Los errores se manejan en un solo lugar
3. **✅ Autenticación automática**: El token se agrega automáticamente
4. **✅ TypeScript**: Tipos definidos para requests y responses
5. **✅ Fácil de testear**: Podemos mockear `apiClient` en tests
6. **✅ Fácil de cambiar**: Si cambia la URL base, se cambia en un solo lugar

---

## 🧩 Componentes Clave

### 1. Server Components vs Client Components

| Server Components | Client Components |
|-------------------|-------------------|
| Por defecto en Next.js | Necesitan `'use client'` |
| Se ejecutan en el servidor | Se ejecutan en el navegador |
| No tienen interactividad | Tienen interactividad (onClick, useState, etc.) |
| Pueden acceder a DB directamente | No pueden acceder a DB |
| Ejemplo: `layout.tsx` | Ejemplo: `page.tsx` en login |

### 2. Layouts

Los layouts son componentes que **envuelven** otras páginas:

```
app/
├── layout.tsx              ← Layout raíz (TODAS las páginas)
└── employees/
    ├── layout.tsx          ← Layout de empleados (solo /employees/*)
    └── page.tsx            ← Página de empleados
```

**Jerarquía de Layouts:**
```
┌─────────────────────────────────────────┐
│ app/layout.tsx (Raíz)                   │
│  ┌─────────────────────────────────┐    │
│  │ app/employees/layout.tsx        │    │
│  │  ┌───────────────────────────┐  │    │
│  │  │ app/employees/page.tsx    │  │    │
│  │  │ (Contenido de la página)  │  │    │
│  │  └───────────────────────────┘  │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

### 3. Componentes Reutilizables (components/ui/)

Componentes pequeños y reutilizables que siguen el patrón de **Componentes Tontos**:

- **Input**: Campo de entrada con validación
- **Button**: Botón con variantes (primary, secondary, danger)
- **Modal**: Ventana modal genérica
- **Table**: Tabla genérica

**Ejemplo de uso:**
```tsx
// En cualquier página
import { Button } from '@/components/ui/Button';

<Button variant="primary" onClick={handleClick}>
  Crear Empleado
</Button>
```

---

## 🔐 Gestión de Autenticación

### Flujo de Autenticación:

```
1. Usuario hace login
   ↓
2. Backend responde con token JWT
   ↓
3. saveAuthData() guarda en localStorage:
   - auth_token: "eyJhbGc..."
   - user_id: "123"
   - expires_at: 1234567890
   ↓
4. En cada petición, getAuthToken() obtiene el token
   ↓
5. apiClient agrega header: Authorization: Bearer eyJhbGc...
   ↓
6. Backend valida el token y responde
```

### Funciones de Autenticación (lib/auth.ts):

```typescript
// Guardar datos de autenticación
saveAuthData(loginResponse)

// Obtener token actual
const token = getAuthToken()

// Verificar si está autenticado
if (isAuthenticated()) {
  // Usuario autenticado
}

// Cerrar sesión
clearAuthData()
```

### localStorage vs Cookies:

| localStorage | Cookies |
|--------------|---------|
| ✅ Fácil de usar | ⚠️ Más complejo |
| ✅ Más espacio (10MB) | ⚠️ Limitado (4KB) |
| ⚠️ No se envía automáticamente | ✅ Se envía automáticamente |
| ⚠️ Vulnerable a XSS | ✅ Puede usar httpOnly |

En nuestro caso usamos **localStorage** por simplicidad en desarrollo.

---

## 🛠️ Cómo Iniciar desde Cero

### Opción 1: Desarrollo Local (Sin Docker)

#### Paso 1: Instalar Dependencias
```bash
# Navega a la carpeta del frontend
cd frontend_basic

# Instala las dependencias de Node.js
npm install
```

#### Paso 2: Configurar Variables de Entorno (Opcional)
```bash
# Crea un archivo .env.local
echo 'NEXT_PUBLIC_API_URL=http://localhost:8080/api' > .env.local
```

#### Paso 3: Iniciar en Modo Desarrollo
```bash
# Inicia el servidor de desarrollo
npm run dev

# La aplicación estará disponible en:
# http://localhost:3000
```

**¿Qué hace `npm run dev`?**
- Inicia el servidor de desarrollo de Next.js
- Habilita **Hot Reload** (cambios se reflejan automáticamente)
- Muestra errores detallados en el navegador
- Compila TypeScript automáticamente

#### Paso 4: Construir para Producción
```bash
# Construye la aplicación optimizada
npm run build

# Inicia en modo producción
npm start
```

---

### Opción 2: Con Docker (Producción)

#### Paso 1: Construir la Imagen
```bash
# Desde la raíz del proyecto
docker build -t frontend-basic ./frontend_basic
```

**¿Qué hace Docker?**
```dockerfile
# 1. Instala dependencias
COPY package.json package-lock.json ./
RUN npm ci

# 2. Construye la aplicación
COPY . .
RUN npm run build

# 3. Prepara para producción
# Crea una imagen ligera con solo lo necesario
```

#### Paso 2: Ejecutar el Contenedor
```bash
docker run -p 3000:3000 frontend-basic
```

#### Paso 3: Con Docker Compose (Recomendado)
```bash
# Levanta todos los servicios (frontend + backend + base de datos)
docker-compose up -d

# Ver logs
docker-compose logs -f frontend-basic

# Detener todos los servicios
docker-compose down
```

---

## 📊 Estructura de Datos (TypeScript)

### Tipos Principales (types/index.ts):

```typescript
// Credenciales de login
interface LoginCredentials {
  email: string;
  password: string;
}

// Respuesta del login
interface LoginResponse {
  token: string;       // JWT token
  user_id: string;     // ID del usuario
  expires_at: number;  // Timestamp de expiración
}

// Empleado
interface Employee {
  id: string;
  name: string;
  email: string;
  created_at: string;
}

// Crear empleado
interface CreateEmployeeRequest {
  name: string;
  email: string;
  password: string;
}

// Error de API
interface ApiError {
  message: string;
  status: number;
}
```

---

## 🎨 Estilos con Tailwind CSS

Tailwind CSS es un framework de **utilidades CSS** - en lugar de escribir CSS, usas clases predefinidas:

```tsx
// ❌ CSS tradicional
<button className="my-button">Click me</button>

// CSS separado
.my-button {
  background-color: blue;
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
}

// ✅ Tailwind CSS
<button className="bg-blue-500 text-white px-4 py-2 rounded">
  Click me
</button>
```

**Clases comunes:**
- `bg-blue-500`: Fondo azul
- `text-white`: Texto blanco
- `px-4 py-2`: Padding horizontal 16px, vertical 8px
- `rounded`: Bordes redondeados
- `hover:bg-blue-600`: Fondo azul oscuro al pasar el mouse

---

## 🔍 Debugging y Troubleshooting

### Ver logs en desarrollo:
```bash
# Los logs aparecen en la terminal donde ejecutaste npm run dev
npm run dev
```

### Ver logs en Docker:
```bash
# Ver logs en tiempo real
docker logs -f frontend-basic

# Ver últimas 50 líneas
docker logs --tail 50 frontend-basic
```

### Errores comunes:

| Error | Solución |
|-------|----------|
| `CORS policy` | Configurar CORS en el backend |
| `Module not found` | Ejecutar `npm install` |
| `Port 3000 already in use` | Matar proceso en puerto 3000: `lsof -ti:3000 \| xargs kill` |
| `fetch failed` | Verificar que el backend esté corriendo |

---

## 📚 Comandos Útiles

```bash
# Desarrollo
npm run dev              # Inicia servidor de desarrollo
npm run build            # Construye para producción
npm start                # Inicia en modo producción
npm run lint             # Ejecuta linter (ESLint)

# Docker
docker build -t frontend-basic ./frontend_basic
docker run -p 3000:3000 frontend-basic
docker-compose up -d frontend-basic
docker-compose logs -f frontend-basic
docker-compose down

# Debugging
lsof -i :3000           # Ver qué proceso usa el puerto 3000
kill -9 <PID>           # Matar proceso
```

---

## 🎓 Conceptos Importantes

### 1. Client-Side vs Server-Side

| Client-Side (Navegador) | Server-Side (Servidor) |
|-------------------------|------------------------|
| Código ejecutado en el navegador | Código ejecutado en el servidor |
| Tiene acceso a `window`, `localStorage` | No tiene acceso a navegador |
| Puede ser interactivo (onClick, etc.) | Solo genera HTML |
| Más lento (descarga código) | Más rápido (HTML ya renderizado) |

### 2. Renderizado

- **CSR (Client-Side Rendering)**: Todo se renderiza en el navegador
- **SSR (Server-Side Rendering)**: Se renderiza en el servidor y se envía HTML
- **SSG (Static Site Generation)**: Se genera HTML en tiempo de build

Next.js soporta los 3 y elige automáticamente el mejor.

### 3. Hydration

Proceso donde React "activa" el HTML estático convirtiéndolo en interactivo:

```
1. Servidor genera HTML estático
   ↓
2. Navegador muestra HTML (página visible pero no interactiva)
   ↓
3. JavaScript se descarga
   ↓
4. React "hidrata" el HTML (ahora es interactivo)
```

---

## 🚀 Próximos Pasos

1. **Agregar más páginas**: Crea nuevas carpetas en `app/`
2. **Mejorar autenticación**: Implementar refresh tokens
3. **Agregar tests**: Jest + React Testing Library
4. **Optimizar imágenes**: Usar `next/image` component
5. **Internacionalización**: Agregar soporte multi-idioma
6. **PWA**: Convertir en Progressive Web App

---

## 📖 Recursos Adicionales

- **Next.js Docs**: https://nextjs.org/docs
- **React Docs**: https://react.dev
- **Tailwind CSS**: https://tailwindcss.com/docs
- **TypeScript**: https://www.typescriptlang.org/docs

---

## 📝 Resumen Ejecutivo

1. **Next.js** es un framework sobre React que facilita el desarrollo
2. **app/layout.tsx** es el layout raíz que envuelve toda la aplicación
3. **app/page.tsx** es la página de inicio (ruta `/`)
4. **lib/api.ts** centraliza todas las peticiones HTTP al backend
5. **'use client'** marca componentes como interactivos (Client Components)
6. **localStorage** guarda el token de autenticación
7. **Tailwind CSS** provee utilidades CSS para estilos rápidos
8. **TypeScript** agrega tipos para mayor seguridad

**Flujo básico:**
```
Usuario visita → Next.js carga layout → Carga página → Login → 
Guarda token → Accede a empleados → apiClient hace fetch → 
Backend responde → Se muestran datos
```
