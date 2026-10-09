# 🚚 Plataforma de Gestión de Última Milla (Last-Mile Delivery) en Go

Solución integral de última milla para logística y distribución urbana inspirada en **Beetrack LastMile / DispatchTrack**, construida con arquitectura de alto rendimiento con **Go (Golang)** en el backend, **PostgreSQL + PostGIS**, y frontend web/móvil interactivo con **OpenStreetMap + Leaflet**.

---

## 🌟 Características Principales

### 1. 🏢 Torre de Control (Backoffice de Despacho y Monitoreo)
- **Mapa interactivo en tiempo real**: Visualización de la flota de vehículos en movimiento satelital con OpenStreetMap y Leaflet.
- **Auditoría de Entregas (POD)**: Inspección inmediata de la firma digital, foto de evidencia de la entrega, datos del receptor (Nombre, RUT/DNI) y geolocalización de la parada.
- **Gestión de Guías de Despacho**: Creación, importación masiva y asignación a rutas.
- **Métricas operacionales en vivo**: Pedidos en ruta, entregas exitosas, rechazos tipificados y móviles conectados vía WebSockets.

### 2. 📱 App Móvil del Conductor (PWA Mobile-First)
- **Hoja de ruta digital**: Secuencia ordenada de paradas con estados: *Pendiente*, *En camino*, *Entregado*, *Fallido*.
- **Transmisión de GPS en segundo plano**: Conexión nativa con el GPS del smartphone del conductor (`navigator.geolocation`) que transmite la ubicación al backend de Go en tiempo real.
- **Integración con Waze y Google Maps**: Botón directo para iniciar la navegación a la dirección de entrega con un toque.
- **Prueba de Entrega Digital (POD)**:
  - Firma táctil en pantalla (HTML5 Canvas).
  - Captura fotográfica con la cámara del dispositivo móvil.
  - Registro de causales de fallo (Cliente ausente, dirección incorrecta, etc.).

### 3. 📦 Portal de Seguimiento al Cliente Final (Live Tracking)
- Acceso público vía enlace único (`/tracking/[tracking_number]`) sin login ni instalación de apps.
- **Mapa satelital en vivo**: El cliente ve aproximarse el móvil de reparto a su domicilio con tiempo estimado (ETA).
- **Comprobante digital descargable**: Al concretarse la entrega, el cliente visualiza su recibo con firma y hora.

### 4. ⚡ Simulador de Telemetría GPS en Go
- Módulo en Go (`cmd/simulator/main.go`) que simula el movimiento fluido de un vehículo por la ciudad para demostraciones y pruebas sin necesidad de salir a la calle.

---

## 🏗️ Arquitectura del Sistema

```text
PLATAFORMA-TRANSPORTE/
├── docker-compose.yml              # PostgreSQL 16 + PostGIS
├── backend/                        # Backend de alto rendimiento en Go
│   ├── go.mod
│   ├── cmd/
│   │   ├── server/main.go          # Servidor HTTP Gin + WebSockets
│   │   └── simulator/main.go       # Simulador GPS en tiempo real
│   ├── internal/
│   │   ├── config/config.go        # Variables de entorno
│   │   ├── database/db.go          # Conexión DB y auto-migraciones
│   │   ├── handlers/handlers.go    # Endpoints REST (Órdenes, Rutas, POD, Telemetría)
│   │   ├── models/models.go        # Estructuras de dominio
│   │   ├── repository/repository.go# Queries SQL nativas para PostgreSQL
│   │   └── websocket/hub.go        # Hub WebSocket concurrente
│   └── migrations/
│       ├── 001_init.sql            # Esquema DDL con índices geoespaciales
│       └── 002_seed.sql            # Datos de prueba (choferes, móviles, pedidos)
└── frontend/                       # Next.js 14 + TailwindCSS + Leaflet
    ├── package.json
    ├── src/
    │   ├── app/
    │   │   ├── page.tsx            # Torre de Control (Dashboard Despachador)
    │   │   ├── driver/page.tsx     # App Móvil Conductor (PWA / GPS / POD)
    │   │   └── tracking/[id]/page.tsx # Tracking en vivo del Cliente Final
    │   ├── components/
    │   │   ├── Map.tsx             # Componente Leaflet con OpenStreetMap
    │   │   └── PODModal.tsx        # Modal de firma digital y foto
    │   ├── lib/api.ts              # Cliente HTTP para backend Go
    │   └── types/index.ts          # Interfaces TypeScript
```

---

## 🚀 Puesta en Marcha Rápida

### 1. Iniciar la Base de Datos (PostgreSQL + PostGIS)
Ejecuta en la raíz del proyecto:
```bash
docker compose up -d
```
> Esto iniciará PostgreSQL en el puerto `5432` y ejecutará automáticamente los scripts `001_init.sql` y `002_seed.sql`.

---

### 2. Iniciar el Backend en Go
En una terminal:
```bash
cd backend
go run cmd/server/main.go
```
El servidor arrancará en `http://localhost:8080`:
- REST API: `http://localhost:8080/api/v1`
- WebSocket Server: `ws://localhost:8080/ws`
- Healthcheck: `http://localhost:8080/health`

---

### 3. Iniciar el Frontend (Torre de Control + App Conductor + Tracking)
En otra terminal:
```bash
cd frontend
npm install
npm run dev
```
La aplicación web estará disponible en `http://localhost:3000`:
- **Torre de Control**: `http://localhost:3000`
- **App Conductor**: `http://localhost:3000/driver`
- **Tracking Cliente en Vivo**: `http://localhost:3000/tracking/TRK-1001`

---

### 4. (Opcional) Probar con el Simulador GPS en Vivo
Para ver el furgón moverse en tiempo real por el mapa satelital sin conducir:
```bash
cd backend
go run cmd/simulator/main.go
```
Verás las coordenadas emitidas segundo a segundo y cómo el camión avanza fluidamente tanto en la Torre de Control como en el celular del cliente.

---

## 📡 Documentación de la API REST (Backend Go)

| Método | Endpoint | Descripción |
|---|---|---|
| `GET` | `/api/v1/orders` | Lista todas las órdenes y guías de despacho |
| `POST` | `/api/v1/orders` | Crea una nueva orden de despacho |
| `POST` | `/api/v1/orders/import` | Carga masiva de pedidos (Batch) |
| `GET` | `/api/v1/routes` | Lista de rutas activas e históricas |
| `GET` | `/api/v1/routes/:id` | Detalle de ruta con paradas ordenadas y POD |
| `POST` | `/api/v1/routes` | Creación y armado de ruta con móvil y chofer |
| `PATCH` | `/api/v1/routes/stops/:stopId/status` | Cambia estado de la parada (`en_route`, `delivered`, etc.) |
| `POST` | `/api/v1/pod` | Registro de Prueba de Entrega (firma táctil, foto, RUT receptor) |
| `POST` | `/api/v1/telemetry` | Recepción de telemetría GPS del móvil |
| `GET` | `/api/v1/telemetry/positions` | Última ubicación conocida de todos los móviles |
| `GET` | `/api/v1/tracking/:trackingNumber` | Consulta pública del estado del pedido para el cliente final |
| `GET` | `/api/v1/drivers` | Catálogo de choferes |
| `GET` | `/api/v1/vehicles` | Catálogo de móviles / vehículos |
| `GET` | `/ws` | Canal WebSocket para eventos y telemetría en tiempo real |
