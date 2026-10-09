-- Migración 007: Vinculación segura de identidad User <-> Driver (Fase 2B)
-- Agrega columna driver_id nullable a la tabla users con FK a drivers(id)

-- 1. Agregar columna driver_id a users si no existe
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS driver_id UUID REFERENCES drivers(id) ON DELETE SET NULL;

-- 2. Crear índice para optimizar búsquedas por driver_id
CREATE INDEX IF NOT EXISTS idx_users_driver_id ON users(driver_id);

-- 3. Backfill explícito e inequívoco para usuarios legacy con role='driver'
-- Solo asocia coincidencias exactas 1:1 basadas en email normalizado
UPDATE users u
SET driver_id = d.id
FROM drivers d
WHERE u.role = 'driver'
  AND u.driver_id IS NULL
  AND LOWER(u.email) = LOWER(d.email)
  AND (
      SELECT COUNT(*) 
      FROM drivers d2 
      WHERE LOWER(d2.email) = LOWER(u.email)
  ) = 1;
