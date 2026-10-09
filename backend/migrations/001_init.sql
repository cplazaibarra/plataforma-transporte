-- 001_init.sql: Inicialización de tablas para Plataforma Last-Mile
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

-- 1. Choferes
CREATE TABLE IF NOT EXISTS drivers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    phone VARCHAR(50),
    license_number VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Vehículos / Móviles
CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plate VARCHAR(20) UNIQUE NOT NULL,
    model VARCHAR(100),
    brand VARCHAR(100),
    capacity_kg NUMERIC(10, 2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Pedidos / Guías de Despacho
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tracking_number VARCHAR(50) UNIQUE NOT NULL,
    external_order_id VARCHAR(100),
    customer_name VARCHAR(150) NOT NULL,
    customer_phone VARCHAR(50),
    customer_email VARCHAR(150),
    delivery_address TEXT NOT NULL,
    commune VARCHAR(100),
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    window_start VARCHAR(10) DEFAULT '09:00',
    window_end VARCHAR(10) DEFAULT '18:00',
    packages_count INT DEFAULT 1,
    items_description TEXT,
    status VARCHAR(30) DEFAULT 'ready_for_dispatch', -- ready_for_dispatch, assigned, in_transit, delivered, failed
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. Rutas Diarias
CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    driver_id UUID REFERENCES drivers(id) ON DELETE SET NULL,
    vehicle_id UUID REFERENCES vehicles(id) ON DELETE SET NULL,
    status VARCHAR(30) DEFAULT 'draft', -- draft, assigned, in_progress, completed, cancelled
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    total_stops INT DEFAULT 0,
    completed_stops INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Paradas de la Ruta (Secuencia de entrega)
CREATE TABLE IF NOT EXISTS route_stops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    stop_number INT NOT NULL,
    status VARCHAR(30) DEFAULT 'pending', -- pending, en_route, arrived, delivered, failed
    eta TIMESTAMPTZ,
    arrived_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_route_stop UNIQUE (route_id, stop_number)
);

-- 6. Prueba de Entrega (POD - Proof of Delivery)
CREATE TABLE IF NOT EXISTS delivery_proofs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    route_stop_id UUID NOT NULL UNIQUE REFERENCES route_stops(id) ON DELETE CASCADE,
    recipient_name VARCHAR(150),
    recipient_dni VARCHAR(50),
    recipient_relationship VARCHAR(50), -- titular, conserje, familiar, vecino, etc.
    signature_base64 TEXT,
    photo_url TEXT,
    failure_reason VARCHAR(100), -- cliente_ausente, direccion_incorrecta, rechazo_producto, zona_peligrosa, etc.
    notes TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    submitted_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. Telemetría GPS en tiempo real de los móviles
CREATE TABLE IF NOT EXISTS vehicle_telemetry (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id UUID REFERENCES vehicles(id) ON DELETE CASCADE,
    driver_id UUID REFERENCES drivers(id) ON DELETE SET NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    speed DOUBLE PRECISION DEFAULT 0,
    heading DOUBLE PRECISION DEFAULT 0,
    battery_level INT,
    recorded_at TIMESTAMPTZ DEFAULT NOW()
);

-- 8. Última posición conocida (para consultas rápidas del mapa de torre de control)
CREATE TABLE IF NOT EXISTS vehicle_last_positions (
    vehicle_id UUID PRIMARY KEY REFERENCES vehicles(id) ON DELETE CASCADE,
    driver_id UUID REFERENCES drivers(id) ON DELETE SET NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    speed DOUBLE PRECISION DEFAULT 0,
    heading DOUBLE PRECISION DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Índices para alta concurrencia y consultas geoespaciales
CREATE INDEX IF NOT EXISTS idx_orders_tracking ON orders(tracking_number);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_route_stops_route ON route_stops(route_id);
CREATE INDEX IF NOT EXISTS idx_route_stops_order ON route_stops(order_id);
CREATE INDEX IF NOT EXISTS idx_telemetry_vehicle ON vehicle_telemetry(vehicle_id, recorded_at DESC);
