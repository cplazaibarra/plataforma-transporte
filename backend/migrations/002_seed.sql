-- 002_seed.sql: Datos de prueba con UUIDs hexadecimales válidos

-- Choferes de prueba
INSERT INTO drivers (id, name, email, phone, license_number, status)
VALUES 
    ('d1000000-0000-0000-0000-000000000001', 'Carlos Morales', 'cmorales@transporte.com', '+56911223344', 'CL-B-998811', 'active'),
    ('d2000000-0000-0000-0000-000000000002', 'Juan Pérez', 'jperez@transporte.com', '+56955667788', 'CL-B-774411', 'active')
ON CONFLICT (id) DO NOTHING;

-- Vehículos de prueba (a100... en lugar de v100...)
INSERT INTO vehicles (id, plate, model, brand, capacity_kg, status)
VALUES 
    ('a1000000-0000-0000-0000-000000000001', 'KJ-DF-44', 'Partner', 'Peugeot', 800.00, 'active'),
    ('a2000000-0000-0000-0000-000000000002', 'LP-RR-92', 'Kangoo', 'Renault', 750.00, 'active')
ON CONFLICT (id) DO NOTHING;

-- Pedidos / Guías de despacho (b100... en lugar de o100...)
INSERT INTO orders (id, tracking_number, external_order_id, customer_name, customer_phone, customer_email, delivery_address, commune, latitude, longitude, window_start, window_end, packages_count, items_description, status)
VALUES
    ('b1000000-0000-0000-0000-000000000001', 'TRK-1001', 'ORD-98210', 'María Silva', '+56987654321', 'maria.silva@email.com', 'Av. Providencia 1240, Depto 402', 'Providencia', -33.4285, -70.6184, '09:00', '13:00', 1, '1x Smart TV 55 Pulgadas', 'in_transit'),
    ('b2000000-0000-0000-0000-000000000002', 'TRK-1002', 'ORD-98211', 'Gonzalo Tapia', '+56999887766', 'gonzalo.t@email.com', 'Av. Las Condes 8540, Of 301', 'Las Condes', -33.3912, -70.5489, '10:00', '14:00', 2, '2x Cajas de Herramientas Eléctricas', 'assigned'),
    ('b3000000-0000-0000-0000-000000000003', 'TRK-1003', 'ORD-98212', 'Carolina Soto', '+56944332211', 'carolina.soto@email.com', 'Moneda 1050', 'Santiago Centro', -33.4429, -70.6539, '13:00', '17:00', 1, '1x Microondas Digital', 'delivered'),
    ('b4000000-0000-0000-0000-000000000004', 'TRK-1004', 'ORD-98213', 'Felipe Arancibia', '+56933221100', 'felipe.a@email.com', 'Av. Vitacura 4300', 'Vitacura', -33.3980, -70.5890, '14:00', '18:00', 3, '3x Muebles de Oficina', 'ready_for_dispatch')
ON CONFLICT (id) DO NOTHING;

-- Ruta de despacho activa para Carlos Morales (c100... en lugar de r100...)
INSERT INTO routes (id, code, date, driver_id, vehicle_id, status, started_at, total_stops, completed_stops)
VALUES 
    ('c1000000-0000-0000-0000-000000000001', 'RUTA-CENTRO-01', CURRENT_DATE, 'd1000000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001', 'in_progress', NOW() - INTERVAL '45 minutes', 3, 1)
ON CONFLICT (id) DO NOTHING;

-- Paradas de la ruta RUTA-CENTRO-01 (e100... en lugar de s100...)
-- Parada 1: Ya entregada (con POD)
INSERT INTO route_stops (id, route_id, order_id, stop_number, status, arrived_at, completed_at)
VALUES 
    ('e1000000-0000-0000-0000-000000000001', 'c1000000-0000-0000-0000-000000000001', 'b3000000-0000-0000-0000-000000000003', 1, 'delivered', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '25 minutes')
ON CONFLICT (id) DO NOTHING;

-- POD para la parada 1
INSERT INTO delivery_proofs (route_stop_id, recipient_name, recipient_dni, recipient_relationship, failure_reason, latitude, longitude, submitted_at)
VALUES 
    ('e1000000-0000-0000-0000-000000000001', 'Carolina Soto', '17.432.190-5', 'titular', NULL, -33.4429, -70.6539, NOW() - INTERVAL '25 minutes')
ON CONFLICT (route_stop_id) DO NOTHING;

-- Parada 2: En ruta hacia Providencia
INSERT INTO route_stops (id, route_id, order_id, stop_number, status, arrived_at, completed_at)
VALUES 
    ('e2000000-0000-0000-0000-000000000002', 'c1000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001', 2, 'en_route', NULL, NULL)
ON CONFLICT (id) DO NOTHING;

-- Parada 3: Pendiente en Las Condes
INSERT INTO route_stops (id, route_id, order_id, stop_number, status)
VALUES 
    ('e3000000-0000-0000-0000-000000000003', 'c1000000-0000-0000-0000-000000000001', 'b2000000-0000-0000-0000-000000000002', 3, 'pending')
ON CONFLICT (id) DO NOTHING;

-- Posición de telemetría del móvil KJ-DF-44
INSERT INTO vehicle_last_positions (vehicle_id, driver_id, latitude, longitude, speed, heading, updated_at)
VALUES 
    ('a1000000-0000-0000-0000-000000000001', 'd1000000-0000-0000-0000-000000000001', -33.4310, -70.6250, 38.5, 75.0, NOW())
ON CONFLICT (vehicle_id) DO UPDATE SET 
    latitude = EXCLUDED.latitude,
    longitude = EXCLUDED.longitude,
    speed = EXCLUDED.speed,
    heading = EXCLUDED.heading,
    updated_at = EXCLUDED.updated_at;
