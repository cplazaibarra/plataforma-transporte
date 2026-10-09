-- 003_users.sql: Tabla de usuarios y roles del sistema

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255) DEFAULT 'default_secret',
    role VARCHAR(30) NOT NULL CHECK (role IN ('admin', 'operator', 'driver')),
    vehicle_id UUID REFERENCES vehicles(id) ON DELETE SET NULL,
    phone VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_vehicle ON users(vehicle_id);

-- Usuarios iniciales de prueba
INSERT INTO users (id, name, email, role, vehicle_id, phone, status)
VALUES
    ('f1000000-0000-0000-0000-000000000001', 'Administrador Principal', 'admin@lastmile.com', 'admin', NULL, '+56900000001', 'active'),
    ('f2000000-0000-0000-0000-000000000002', 'Operador Monitoreo (Lectura)', 'operador@lastmile.com', 'operator', NULL, '+56900000002', 'active'),
    ('f3000000-0000-0000-0000-000000000003', 'Carlos Morales (Transportista)', 'cmorales@transporte.com', 'driver', 'a1000000-0000-0000-0000-000000000001', '+56911223344', 'active'),
    ('f4000000-0000-0000-0000-000000000004', 'Juan Pérez (Transportista)', 'jperez@transporte.com', 'driver', 'a2000000-0000-0000-0000-000000000002', '+56955667788', 'active')
ON CONFLICT (id) DO NOTHING;
