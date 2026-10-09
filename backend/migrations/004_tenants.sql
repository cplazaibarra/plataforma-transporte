-- 004_tenants.sql: Tabla de tenants y vinculación con usuarios para soporte multi-tenant

CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    legal_id VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('carrier_3pl', 'shipper_retailer')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Aseguramos que el identificador legal/fiscal sea único
CREATE UNIQUE INDEX IF NOT EXISTS idx_tenants_legal_id ON tenants(legal_id);
-- Índice para filtrado de empresas activas
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);

CREATE TABLE IF NOT EXISTS tenant_users (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'operator' CHECK (role IN ('admin', 'operator')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (tenant_id, user_id)
);

-- Índice para consultar rápidamente los tenants de un usuario
CREATE INDEX IF NOT EXISTS idx_tenant_users_user_id ON tenant_users(user_id);

-- Seed del Tenant Legacy Default
INSERT INTO tenants (id, name, legal_id, type, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Empresa Principal Default',
    '76.000.000-0',
    'carrier_3pl',
    'active'
) ON CONFLICT (id) DO NOTHING;

-- Vinculación automática en tenant_users de los usuarios 'admin' y 'operator' existentes (creados en 003_users.sql)
INSERT INTO tenant_users (tenant_id, user_id, role, status)
SELECT 
    '00000000-0000-0000-0000-000000000001'::UUID,
    id,
    role,
    'active'
FROM users 
WHERE role IN ('admin', 'operator')
ON CONFLICT (tenant_id, user_id) DO NOTHING;
