-- Migración 006: Tenantización aditiva de orders y routes (Fase 1C)
-- Columnas NULLABLE y backfill hacia el Tenant Legacy Default

-- 1. Agregar columnas origin_tenant_id y carrier_tenant_id a orders si no existen
ALTER TABLE orders 
ADD COLUMN IF NOT EXISTS origin_tenant_id UUID REFERENCES tenants(id) ON DELETE RESTRICT;

ALTER TABLE orders 
ADD COLUMN IF NOT EXISTS carrier_tenant_id UUID REFERENCES tenants(id) ON DELETE RESTRICT;

-- 2. Agregar columna carrier_tenant_id a routes si no existe
ALTER TABLE routes 
ADD COLUMN IF NOT EXISTS carrier_tenant_id UUID REFERENCES tenants(id) ON DELETE RESTRICT;

-- 3. Índices para acelerar consultas multi-tenant sobre orders y routes
CREATE INDEX IF NOT EXISTS idx_orders_origin_tenant_id ON orders(origin_tenant_id);
CREATE INDEX IF NOT EXISTS idx_orders_carrier_tenant_id ON orders(carrier_tenant_id);
CREATE INDEX IF NOT EXISTS idx_routes_carrier_tenant_id ON routes(carrier_tenant_id);

-- 4. Backfill explícito e idempotente hacia el Tenant Legacy Default
-- Tenant Legacy Default: 00000000-0000-0000-0000-000000000001

UPDATE orders
SET origin_tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE origin_tenant_id IS NULL;

UPDATE orders
SET carrier_tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE carrier_tenant_id IS NULL;

UPDATE routes
SET carrier_tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE carrier_tenant_id IS NULL;
