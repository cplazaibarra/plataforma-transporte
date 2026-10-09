-- 005_tenant_fleet.sql: Tablas de vinculación multi-tenant para conductores y vehículos

-- 1. Vinculación de conductores por tenant con soporte de habilitación y bloqueo local
CREATE TABLE IF NOT EXISTS tenant_drivers (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    driver_id UUID NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'blocked', 'pending')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (tenant_id, driver_id)
);

-- Índice para consultar rápidamente los tenants donde un conductor está habilitado
CREATE INDEX IF NOT EXISTS idx_tenant_drivers_driver_id ON tenant_drivers(driver_id);
-- Índice para filtrado de conductores por estado dentro de un tenant
CREATE INDEX IF NOT EXISTS idx_tenant_drivers_tenant_status ON tenant_drivers(tenant_id, status);

-- 2. Vinculación de vehículos por tenant con soporte de habilitación y bloqueo local
CREATE TABLE IF NOT EXISTS tenant_vehicles (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'blocked', 'pending')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (tenant_id, vehicle_id)
);

-- Índice para consultar rápidamente los tenants donde un vehículo está habilitado
CREATE INDEX IF NOT EXISTS idx_tenant_vehicles_vehicle_id ON tenant_vehicles(vehicle_id);
-- Índice para filtrado de vehículos por estado dentro de un tenant
CREATE INDEX IF NOT EXISTS idx_tenant_vehicles_tenant_status ON tenant_vehicles(tenant_id, status);

-- 3. Vinculación automática inicial de la flota existente hacia el Tenant Legacy Default
-- Asocia todos los conductores existentes en drivers al Tenant Legacy Default
INSERT INTO tenant_drivers (tenant_id, driver_id, status)
SELECT 
    '00000000-0000-0000-0000-000000000001'::UUID,
    id,
    'active'
FROM drivers
ON CONFLICT (tenant_id, driver_id) DO NOTHING;

-- Asocia todos los vehículos existentes en vehicles al Tenant Legacy Default
INSERT INTO tenant_vehicles (tenant_id, vehicle_id, status)
SELECT 
    '00000000-0000-0000-0000-000000000001'::UUID,
    id,
    'active'
FROM vehicles
ON CONFLICT (tenant_id, vehicle_id) DO NOTHING;
