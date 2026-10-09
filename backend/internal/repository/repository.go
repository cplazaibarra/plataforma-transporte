package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"lastmile-backend/internal/database"
	"lastmile-backend/internal/models"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// ---------------- ORDERS ----------------

// GetAllOrders obtiene todas las órdenes (uso legacy / transicional sin contexto tenant)
func (r *Repository) GetAllOrders() ([]models.Order, error) {
	query := `
		SELECT id, origin_tenant_id::text, carrier_tenant_id::text, tracking_number, COALESCE(external_order_id, ''), customer_name, 
		       COALESCE(customer_phone, ''), COALESCE(customer_email, ''), delivery_address, 
		       COALESCE(commune, ''), latitude, longitude, window_start, window_end, 
		       packages_count, COALESCE(items_description, ''), status, created_at
		FROM orders
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var originID, carrierID *string
		if err := rows.Scan(
			&o.ID, &originID, &carrierID, &o.TrackingNumber, &o.ExternalOrderID, &o.CustomerName,
			&o.CustomerPhone, &o.CustomerEmail, &o.DeliveryAddress,
			&o.Commune, &o.Latitude, &o.Longitude, &o.WindowStart, &o.WindowEnd,
			&o.PackagesCount, &o.ItemsDescription, &o.Status, &o.CreatedAt,
		); err != nil {
			return nil, err
		}
		o.OriginTenantID = originID
		o.CarrierTenantID = carrierID
		orders = append(orders, o)
	}
	return orders, nil
}

// GetOrdersByTenant retorna las órdenes donde el tenant es propietario (origin) o transportista (carrier)
func (r *Repository) GetOrdersByTenant(tenantID string) ([]models.Order, error) {
	query := `
		SELECT id, origin_tenant_id::text, carrier_tenant_id::text, tracking_number, 
		       COALESCE(external_order_id, ''), customer_name, 
		       COALESCE(customer_phone, ''), COALESCE(customer_email, ''), delivery_address, 
		       COALESCE(commune, ''), latitude, longitude, window_start, window_end, 
		       packages_count, COALESCE(items_description, ''), status, created_at
		FROM orders
		WHERE origin_tenant_id = $1 OR carrier_tenant_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var originID, carrierID *string
		if err := rows.Scan(
			&o.ID, &originID, &carrierID, &o.TrackingNumber, &o.ExternalOrderID, &o.CustomerName,
			&o.CustomerPhone, &o.CustomerEmail, &o.DeliveryAddress,
			&o.Commune, &o.Latitude, &o.Longitude, &o.WindowStart, &o.WindowEnd,
			&o.PackagesCount, &o.ItemsDescription, &o.Status, &o.CreatedAt,
		); err != nil {
			return nil, err
		}
		o.OriginTenantID = originID
		o.CarrierTenantID = carrierID
		orders = append(orders, o)
	}
	return orders, nil
}

// GetOrderByID busca una orden por ID (uso legacy / transicional)
func (r *Repository) GetOrderByID(id string) (*models.Order, error) {
	query := `
		SELECT id, origin_tenant_id::text, carrier_tenant_id::text, tracking_number, COALESCE(external_order_id, ''), customer_name, 
		       COALESCE(customer_phone, ''), COALESCE(customer_email, ''), delivery_address, 
		       COALESCE(commune, ''), latitude, longitude, window_start, window_end, 
		       packages_count, COALESCE(items_description, ''), status, created_at
		FROM orders WHERE id = $1`

	var o models.Order
	var originID, carrierID *string
	err := r.db.QueryRow(query, id).Scan(
		&o.ID, &originID, &carrierID, &o.TrackingNumber, &o.ExternalOrderID, &o.CustomerName,
		&o.CustomerPhone, &o.CustomerEmail, &o.DeliveryAddress,
		&o.Commune, &o.Latitude, &o.Longitude, &o.WindowStart, &o.WindowEnd,
		&o.PackagesCount, &o.ItemsDescription, &o.Status, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	o.OriginTenantID = originID
	o.CarrierTenantID = carrierID
	return &o, nil
}

// GetOrderByIDAndTenant busca una orden por ID asegurando visibilidad dual para el tenant (origin o carrier)
func (r *Repository) GetOrderByIDAndTenant(id, tenantID string) (*models.Order, error) {
	query := `
		SELECT id, origin_tenant_id::text, carrier_tenant_id::text, tracking_number, 
		       COALESCE(external_order_id, ''), customer_name, 
		       COALESCE(customer_phone, ''), COALESCE(customer_email, ''), delivery_address, 
		       COALESCE(commune, ''), latitude, longitude, window_start, window_end, 
		       packages_count, COALESCE(items_description, ''), status, created_at
		FROM orders 
		WHERE id = $1 AND (origin_tenant_id = $2 OR carrier_tenant_id = $2)`

	var o models.Order
	var originID, carrierID *string
	err := r.db.QueryRow(query, id, tenantID).Scan(
		&o.ID, &originID, &carrierID, &o.TrackingNumber, &o.ExternalOrderID, &o.CustomerName,
		&o.CustomerPhone, &o.CustomerEmail, &o.DeliveryAddress,
		&o.Commune, &o.Latitude, &o.Longitude, &o.WindowStart, &o.WindowEnd,
		&o.PackagesCount, &o.ItemsDescription, &o.Status, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	o.OriginTenantID = originID
	o.CarrierTenantID = carrierID
	return &o, nil
}

func (r *Repository) GetOrderByTrackingNumber(trk string) (*models.Order, error) {
	query := `
		SELECT id, origin_tenant_id::text, carrier_tenant_id::text, tracking_number, COALESCE(external_order_id, ''), customer_name, 
		       COALESCE(customer_phone, ''), COALESCE(customer_email, ''), delivery_address, 
		       COALESCE(commune, ''), latitude, longitude, window_start, window_end, 
		       packages_count, COALESCE(items_description, ''), status, created_at
		FROM orders WHERE tracking_number = $1`

	var o models.Order
	var originID, carrierID *string
	err := r.db.QueryRow(query, trk).Scan(
		&o.ID, &originID, &carrierID, &o.TrackingNumber, &o.ExternalOrderID, &o.CustomerName,
		&o.CustomerPhone, &o.CustomerEmail, &o.DeliveryAddress,
		&o.Commune, &o.Latitude, &o.Longitude, &o.WindowStart, &o.WindowEnd,
		&o.PackagesCount, &o.ItemsDescription, &o.Status, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	o.OriginTenantID = originID
	o.CarrierTenantID = carrierID
	return &o, nil
}

func (r *Repository) CreateOrder(o *models.Order) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	if o.TrackingNumber == "" {
		o.TrackingNumber = fmt.Sprintf("TRK-%d", time.Now().Unix()%1000000)
	}
	if o.Status == "" {
		o.Status = "ready_for_dispatch"
	}

	query := `
		INSERT INTO orders (id, origin_tenant_id, carrier_tenant_id, tracking_number, external_order_id, 
		                    customer_name, customer_phone, customer_email, delivery_address, commune, 
		                    latitude, longitude, window_start, window_end, packages_count, 
		                    items_description, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING created_at`

	return r.db.QueryRow(
		query, o.ID, o.OriginTenantID, o.CarrierTenantID, o.TrackingNumber, o.ExternalOrderID, 
		o.CustomerName, o.CustomerPhone, o.CustomerEmail, o.DeliveryAddress, o.Commune, 
		o.Latitude, o.Longitude, o.WindowStart, o.WindowEnd, o.PackagesCount, 
		o.ItemsDescription, o.Status,
	).Scan(&o.CreatedAt)
}

// CreateOrderWithTenant crea una orden forzando origin_tenant_id al activeTenantID del servidor
func (r *Repository) CreateOrderWithTenant(o *models.Order, activeTenantID string) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	if o.TrackingNumber == "" {
		o.TrackingNumber = fmt.Sprintf("TRK-%d", time.Now().Unix()%1000000)
	}
	if o.Status == "" {
		o.Status = "ready_for_dispatch"
	}
	// Forzar origin_tenant_id inmutablemente desde el contexto autenticado
	o.OriginTenantID = &activeTenantID

	// Si se especifica carrier_tenant_id, validar que exista como tenant activo
	if o.CarrierTenantID != nil && *o.CarrierTenantID != "" {
		var exists bool
		err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM tenants WHERE id = $1 AND status = 'active')`, *o.CarrierTenantID).Scan(&exists)
		if err != nil || !exists {
			return fmt.Errorf("carrier_tenant_id no válido o inactivo")
		}
	} else {
		o.CarrierTenantID = nil
	}

	query := `
		INSERT INTO orders (id, origin_tenant_id, carrier_tenant_id, tracking_number, external_order_id, 
		                    customer_name, customer_phone, customer_email, delivery_address, commune, 
		                    latitude, longitude, window_start, window_end, packages_count, 
		                    items_description, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING created_at`

	return r.db.QueryRow(
		query, o.ID, o.OriginTenantID, o.CarrierTenantID, o.TrackingNumber, o.ExternalOrderID, 
		o.CustomerName, o.CustomerPhone, o.CustomerEmail, o.DeliveryAddress, o.Commune, 
		o.Latitude, o.Longitude, o.WindowStart, o.WindowEnd, o.PackagesCount, 
		o.ItemsDescription, o.Status,
	).Scan(&o.CreatedAt)
}

func (r *Repository) UpdateOrderStatus(id, status string) error {
	_, err := r.db.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, id)
	return err
}

// ---------------- ROUTES & STOPS ----------------

// GetAllRoutes lista todas las rutas (uso legacy / transicional)
func (r *Repository) GetAllRoutes() ([]models.Route, error) {
	query := `
		SELECT r.id, r.carrier_tenant_id::text, r.code, r.date::text, r.driver_id, d.name as driver_name, 
		       r.vehicle_id, v.plate as vehicle_plate, r.status, r.started_at, 
		       r.completed_at, r.total_stops, r.completed_stops, r.created_at
		FROM routes r
		LEFT JOIN drivers d ON r.driver_id = d.id
		LEFT JOIN vehicles v ON r.vehicle_id = v.id
		ORDER BY r.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var rt models.Route
		var carrierID *string
		if err := rows.Scan(
			&rt.ID, &carrierID, &rt.Code, &rt.Date, &rt.DriverID, &rt.DriverName,
			&rt.VehicleID, &rt.VehiclePlate, &rt.Status, &rt.StartedAt,
			&rt.CompletedAt, &rt.TotalStops, &rt.CompletedStops, &rt.CreatedAt,
		); err != nil {
			return nil, err
		}
		rt.CarrierTenantID = carrierID
		routes = append(routes, rt)
	}
	return routes, nil
}

// GetRoutesByTenant lista las rutas asignadas exclusivamente a la empresa transportista (carrier)
func (r *Repository) GetRoutesByTenant(tenantID string) ([]models.Route, error) {
	query := `
		SELECT r.id, r.carrier_tenant_id::text, r.code, r.date::text, r.driver_id, d.name as driver_name, 
		       r.vehicle_id, v.plate as vehicle_plate, r.status, r.started_at, 
		       r.completed_at, r.total_stops, r.completed_stops, r.created_at
		FROM routes r
		LEFT JOIN drivers d ON r.driver_id = d.id
		LEFT JOIN vehicles v ON r.vehicle_id = v.id
		WHERE r.carrier_tenant_id = $1
		ORDER BY r.created_at DESC`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var rt models.Route
		var carrierID *string
		if err := rows.Scan(
			&rt.ID, &carrierID, &rt.Code, &rt.Date, &rt.DriverID, &rt.DriverName,
			&rt.VehicleID, &rt.VehiclePlate, &rt.Status, &rt.StartedAt,
			&rt.CompletedAt, &rt.TotalStops, &rt.CompletedStops, &rt.CreatedAt,
		); err != nil {
			return nil, err
		}
		rt.CarrierTenantID = carrierID
		routes = append(routes, rt)
	}
	return routes, nil
}

// GetRouteWithStops obtiene una ruta y sus paradas (uso legacy / transicional)
func (r *Repository) GetRouteWithStops(routeID string) (*models.Route, error) {
	query := `
		SELECT r.id, r.carrier_tenant_id::text, r.code, r.date::text, r.driver_id, d.name as driver_name, 
		       r.vehicle_id, v.plate as vehicle_plate, r.status, r.started_at, 
		       r.completed_at, r.total_stops, r.completed_stops, r.created_at
		FROM routes r
		LEFT JOIN drivers d ON r.driver_id = d.id
		LEFT JOIN vehicles v ON r.vehicle_id = v.id
		WHERE r.id = $1`

	var rt models.Route
	var carrierID *string
	err := r.db.QueryRow(query, routeID).Scan(
		&rt.ID, &carrierID, &rt.Code, &rt.Date, &rt.DriverID, &rt.DriverName,
		&rt.VehicleID, &rt.VehiclePlate, &rt.Status, &rt.StartedAt,
		&rt.CompletedAt, &rt.TotalStops, &rt.CompletedStops, &rt.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	rt.CarrierTenantID = carrierID

	stopsQuery := `
		SELECT rs.id, rs.route_id, rs.order_id, rs.stop_number, rs.status, rs.eta, 
		       rs.arrived_at, rs.completed_at, COALESCE(rs.notes, ''), rs.created_at,
		       o.tracking_number, COALESCE(o.external_order_id, ''), o.customer_name, 
		       COALESCE(o.customer_phone, ''), o.delivery_address, COALESCE(o.commune, ''), 
		       o.latitude, o.longitude, o.window_start, o.window_end, o.packages_count, 
		       COALESCE(o.items_description, ''), o.status as order_status,
		       dp.id, dp.recipient_name, dp.recipient_dni, dp.recipient_relationship, 
		       COALESCE(dp.signature_base64, ''), COALESCE(dp.photo_url, ''), 
		       COALESCE(dp.failure_reason, ''), COALESCE(dp.notes, ''), dp.submitted_at
		FROM route_stops rs
		INNER JOIN orders o ON rs.order_id = o.id
		LEFT JOIN delivery_proofs dp ON dp.route_stop_id = rs.id
		WHERE rs.route_id = $1
		ORDER BY rs.stop_number ASC`

	rows, err := r.db.Query(stopsQuery, routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.RouteStop
		var o models.Order
		var dpID, dpName, dpDNI, dpRel, dpSig, dpPhoto, dpFail, dpNotes *string
		var dpSubAt *time.Time

		if err := rows.Scan(
			&s.ID, &s.RouteID, &s.OrderID, &s.StopNumber, &s.Status, &s.ETA,
			&s.ArrivedAt, &s.CompletedAt, &s.Notes, &s.CreatedAt,
			&o.TrackingNumber, &o.ExternalOrderID, &o.CustomerName,
			&o.CustomerPhone, &o.DeliveryAddress, &o.Commune,
			&o.Latitude, &o.Longitude, &o.WindowStart, &o.WindowEnd, &o.PackagesCount,
			&o.ItemsDescription, &o.Status,
			&dpID, &dpName, &dpDNI, &dpRel, &dpSig, &dpPhoto, &dpFail, &dpNotes, &dpSubAt,
		); err != nil {
			return nil, err
		}
		o.ID = s.OrderID
		s.Order = &o

		if dpID != nil {
			s.Proof = &models.DeliveryProof{
				ID:                    *dpID,
				RouteStopID:           s.ID,
				RecipientName:         stringValue(dpName),
				RecipientDNI:          stringValue(dpDNI),
				RecipientRelationship: stringValue(dpRel),
				SignatureBase64:       stringValue(dpSig),
				PhotoURL:              stringValue(dpPhoto),
				FailureReason:         stringValue(dpFail),
				Notes:                 stringValue(dpNotes),
				SubmittedAt:           timeValue(dpSubAt),
			}
		}

		rt.Stops = append(rt.Stops, s)
	}

	return &rt, nil
}

// GetRouteWithStopsByTenant obtiene el detalle de la ruta y sus paradas validando carrier_tenant_id
func (r *Repository) GetRouteWithStopsByTenant(routeID, tenantID string) (*models.Route, error) {
	query := `
		SELECT r.id, r.carrier_tenant_id::text, r.code, r.date::text, r.driver_id, d.name as driver_name, 
		       r.vehicle_id, v.plate as vehicle_plate, r.status, r.started_at, 
		       r.completed_at, r.total_stops, r.completed_stops, r.created_at
		FROM routes r
		LEFT JOIN drivers d ON r.driver_id = d.id
		LEFT JOIN vehicles v ON r.vehicle_id = v.id
		WHERE r.id = $1 AND r.carrier_tenant_id = $2`

	var rt models.Route
	var carrierID *string
	err := r.db.QueryRow(query, routeID, tenantID).Scan(
		&rt.ID, &carrierID, &rt.Code, &rt.Date, &rt.DriverID, &rt.DriverName,
		&rt.VehicleID, &rt.VehiclePlate, &rt.Status, &rt.StartedAt,
		&rt.CompletedAt, &rt.TotalStops, &rt.CompletedStops, &rt.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	rt.CarrierTenantID = carrierID

	stopsQuery := `
		SELECT rs.id, rs.route_id, rs.order_id, rs.stop_number, rs.status, rs.eta, 
		       rs.arrived_at, rs.completed_at, COALESCE(rs.notes, ''), rs.created_at,
		       o.origin_tenant_id::text, o.carrier_tenant_id::text,
		       o.tracking_number, COALESCE(o.external_order_id, ''), o.customer_name, 
		       COALESCE(o.customer_phone, ''), o.delivery_address, COALESCE(o.commune, ''), 
		       o.latitude, o.longitude, o.window_start, o.window_end, o.packages_count, 
		       COALESCE(o.items_description, ''), o.status as order_status,
		       dp.id, dp.recipient_name, dp.recipient_dni, dp.recipient_relationship, 
		       COALESCE(dp.signature_base64, ''), COALESCE(dp.photo_url, ''), 
		       COALESCE(dp.failure_reason, ''), COALESCE(dp.notes, ''), dp.submitted_at
		FROM route_stops rs
		INNER JOIN routes r ON rs.route_id = r.id
		INNER JOIN orders o ON rs.order_id = o.id
		LEFT JOIN delivery_proofs dp ON dp.route_stop_id = rs.id
		WHERE rs.route_id = $1 AND r.carrier_tenant_id = $2
		ORDER BY rs.stop_number ASC`

	rows, err := r.db.Query(stopsQuery, routeID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.RouteStop
		var o models.Order
		var originID, carrierOrderID *string
		var dpID, dpName, dpDNI, dpRel, dpSig, dpPhoto, dpFail, dpNotes *string
		var dpSubAt *time.Time

		if err := rows.Scan(
			&s.ID, &s.RouteID, &s.OrderID, &s.StopNumber, &s.Status, &s.ETA,
			&s.ArrivedAt, &s.CompletedAt, &s.Notes, &s.CreatedAt,
			&originID, &carrierOrderID,
			&o.TrackingNumber, &o.ExternalOrderID, &o.CustomerName,
			&o.CustomerPhone, &o.DeliveryAddress, &o.Commune,
			&o.Latitude, &o.Longitude, &o.WindowStart, &o.WindowEnd, &o.PackagesCount,
			&o.ItemsDescription, &o.Status,
			&dpID, &dpName, &dpDNI, &dpRel, &dpSig, &dpPhoto, &dpFail, &dpNotes, &dpSubAt,
		); err != nil {
			return nil, err
		}
		o.ID = s.OrderID
		o.OriginTenantID = originID
		o.CarrierTenantID = carrierOrderID
		s.Order = &o

		if dpID != nil {
			s.Proof = &models.DeliveryProof{
				ID:                    *dpID,
				RouteStopID:           s.ID,
				RecipientName:         stringValue(dpName),
				RecipientDNI:          stringValue(dpDNI),
				RecipientRelationship: stringValue(dpRel),
				SignatureBase64:       stringValue(dpSig),
				PhotoURL:              stringValue(dpPhoto),
				FailureReason:         stringValue(dpFail),
				Notes:                 stringValue(dpNotes),
				SubmittedAt:           timeValue(dpSubAt),
			}
		}

		rt.Stops = append(rt.Stops, s)
	}

	return &rt, nil
}

func (r *Repository) CreateRoute(rt *models.Route, orderIDs []string) error {
	if rt.ID == "" {
		rt.ID = uuid.New().String()
	}
	if rt.Code == "" {
		rt.Code = fmt.Sprintf("RUTA-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%1000)
	}
	rt.TotalStops = len(orderIDs)
	rt.Status = "assigned"

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertRouteQuery := `
		INSERT INTO routes (id, carrier_tenant_id, code, date, driver_id, vehicle_id, status, total_stops, completed_stops)
		VALUES ($1, $2, $3, CURRENT_DATE, $4, $5, $6, $7, 0)
		RETURNING created_at`

	if err := tx.QueryRow(insertRouteQuery, rt.ID, rt.CarrierTenantID, rt.Code, rt.DriverID, rt.VehicleID, rt.Status, rt.TotalStops).Scan(&rt.CreatedAt); err != nil {
		return err
	}

	for i, orderID := range orderIDs {
		stopID := uuid.New().String()
		stopNum := i + 1
		insertStopQuery := `
			INSERT INTO route_stops (id, route_id, order_id, stop_number, status)
			VALUES ($1, $2, $3, $4, 'pending')`
		if _, err := tx.Exec(insertStopQuery, stopID, rt.ID, orderID, stopNum); err != nil {
			return err
		}

		// Actualizar orden como asignada
		if _, err := tx.Exec(`UPDATE orders SET status = 'assigned' WHERE id = $1`, orderID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// CreateRouteWithTenant crea una ruta forzando carrier_tenant_id y validando que todas las órdenes pertenezcan al tenant activo
func (r *Repository) CreateRouteWithTenant(rt *models.Route, orderIDs []string, activeTenantID string) error {
	if rt.ID == "" {
		rt.ID = uuid.New().String()
	}
	if rt.Code == "" {
		rt.Code = fmt.Sprintf("RUTA-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%1000)
	}
	rt.CarrierTenantID = &activeTenantID
	rt.TotalStops = len(orderIDs)
	rt.Status = "assigned"

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Validar que si se asigna conductor, esté activo en tenant_drivers para este carrier_tenant_id
	if rt.DriverID != nil && *rt.DriverID != "" {
		var driverActive bool
		err := tx.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM tenant_drivers 
				WHERE tenant_id = $1 AND driver_id = $2 AND status = 'active'
			)`, activeTenantID, *rt.DriverID).Scan(&driverActive)
		if err != nil || !driverActive {
			return fmt.Errorf("el conductor asignado no está activo en el tenant transportista")
		}
	}

	// 2. Validar que si se asigna vehículo, esté activo en tenant_vehicles para este carrier_tenant_id
	if rt.VehicleID != nil && *rt.VehicleID != "" {
		var vehicleActive bool
		err := tx.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM tenant_vehicles 
				WHERE tenant_id = $1 AND vehicle_id = $2 AND status = 'active'
			)`, activeTenantID, *rt.VehicleID).Scan(&vehicleActive)
		if err != nil || !vehicleActive {
			return fmt.Errorf("el vehículo asignado no está activo en el tenant transportista")
		}
	}

	// 3. Validar que todas las órdenes pertenezcan al tenant activo (como origin o carrier)
	for _, orderID := range orderIDs {
		var orderAuthorized bool
		err := tx.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM orders 
				WHERE id = $1 AND (carrier_tenant_id = $2 OR (carrier_tenant_id IS NULL AND origin_tenant_id = $2))
			)`, orderID, activeTenantID).Scan(&orderAuthorized)
		if err != nil || !orderAuthorized {
			return fmt.Errorf("la orden %s no está autorizada para este tenant o pertenece a otra empresa", orderID)
		}
	}

	// 4. Insertar ruta
	insertRouteQuery := `
		INSERT INTO routes (id, carrier_tenant_id, code, date, driver_id, vehicle_id, status, total_stops, completed_stops)
		VALUES ($1, $2, $3, CURRENT_DATE, $4, $5, $6, $7, 0)
		RETURNING created_at`

	if err := tx.QueryRow(insertRouteQuery, rt.ID, rt.CarrierTenantID, rt.Code, rt.DriverID, rt.VehicleID, rt.Status, rt.TotalStops).Scan(&rt.CreatedAt); err != nil {
		return err
	}

	// 5. Insertar paradas y actualizar órdenes
	for i, orderID := range orderIDs {
		stopID := uuid.New().String()
		stopNum := i + 1
		insertStopQuery := `
			INSERT INTO route_stops (id, route_id, order_id, stop_number, status)
			VALUES ($1, $2, $3, $4, 'pending')`
		if _, err := tx.Exec(insertStopQuery, stopID, rt.ID, orderID, stopNum); err != nil {
			return err
		}

		if _, err := tx.Exec(`
			UPDATE orders 
			SET status = 'assigned', carrier_tenant_id = COALESCE(carrier_tenant_id, $1) 
			WHERE id = $2 AND (carrier_tenant_id = $1 OR origin_tenant_id = $1)`, activeTenantID, orderID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) UpdateStopStatus(stopID, status string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var orderID, routeID string
	err = tx.QueryRow(`SELECT order_id, route_id FROM route_stops WHERE id = $1`, stopID).Scan(&orderID, &routeID)
	if err != nil {
		return err
	}

	var updateStopQuery string
	now := time.Now()
	switch status {
	case "en_route":
		updateStopQuery = `UPDATE route_stops SET status = $1, arrived_at = NULL WHERE id = $2`
		if _, err := tx.Exec(updateStopQuery, status, stopID); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE orders SET status = 'in_transit' WHERE id = $1`, orderID); err != nil {
			return err
		}
	case "arrived":
		updateStopQuery = `UPDATE route_stops SET status = $1, arrived_at = $2 WHERE id = $3`
		if _, err := tx.Exec(updateStopQuery, status, now, stopID); err != nil {
			return err
		}
	case "delivered":
		updateStopQuery = `UPDATE route_stops SET status = $1, completed_at = $2 WHERE id = $3`
		if _, err := tx.Exec(updateStopQuery, status, now, stopID); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE orders SET status = 'delivered' WHERE id = $1`, orderID); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE routes SET completed_stops = completed_stops + 1 WHERE id = $1`, routeID); err != nil {
			return err
		}
	case "failed":
		updateStopQuery = `UPDATE route_stops SET status = $1, completed_at = $2 WHERE id = $3`
		if _, err := tx.Exec(updateStopQuery, status, now, stopID); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE orders SET status = 'failed' WHERE id = $1`, orderID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateStopStatusWithTenant actualiza el estado de una parada verificando que la ruta pertenezca al carrier_tenant_id activo
func (r *Repository) UpdateStopStatusWithTenant(stopID, status, activeTenantID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var orderID, routeID string
	err = tx.QueryRow(`
		SELECT rs.order_id, rs.route_id 
		FROM route_stops rs
		INNER JOIN routes r ON rs.route_id = r.id
		WHERE rs.id = $1 AND r.carrier_tenant_id = $2`, stopID, activeTenantID).Scan(&orderID, &routeID)
	if err != nil {
		return fmt.Errorf("parada no encontrada o no autorizada para el tenant")
	}

	var updateStopQuery string
	now := time.Now()
	switch status {
	case "en_route":
		updateStopQuery = `UPDATE route_stops SET status = $1, arrived_at = NULL WHERE id = $2`
		if _, err := tx.Exec(updateStopQuery, status, stopID); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE orders SET status = 'in_transit' WHERE id = $1`, orderID); err != nil {
			return err
		}
	case "arrived":
		updateStopQuery = `UPDATE route_stops SET status = $1, arrived_at = $2 WHERE id = $3`
		if _, err := tx.Exec(updateStopQuery, status, now, stopID); err != nil {
			return err
		}
	case "delivered":
		updateStopQuery = `UPDATE route_stops SET status = $1, completed_at = $2 WHERE id = $3`
		if _, err := tx.Exec(updateStopQuery, status, now, stopID); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE orders SET status = 'delivered' WHERE id = $1`, orderID); err != nil {
			return err
		}
		if _, err := tx.Exec(`
			UPDATE routes 
			SET completed_stops = (SELECT COUNT(*) FROM route_stops WHERE route_id = $1 AND status IN ('delivered', 'failed'))
			WHERE id = $1`, routeID); err != nil {
			return err
		}
	case "failed":
		updateStopQuery = `UPDATE route_stops SET status = $1, completed_at = $2 WHERE id = $3`
		if _, err := tx.Exec(updateStopQuery, status, now, stopID); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE orders SET status = 'failed' WHERE id = $1`, orderID); err != nil {
			return err
		}
		if _, err := tx.Exec(`
			UPDATE routes 
			SET completed_stops = (SELECT COUNT(*) FROM route_stops WHERE route_id = $1 AND status IN ('delivered', 'failed'))
			WHERE id = $1`, routeID); err != nil {
			return err
		}
	default:
		return fmt.Errorf("estado no válido: %s", status)
	}

	return tx.Commit()
}

// ---------------- PROOF OF DELIVERY (POD) ----------------

func (r *Repository) SaveProofOfDelivery(pod *models.DeliveryProof) error {
	if pod.ID == "" {
		pod.ID = uuid.New().String()
	}

	query := `
		INSERT INTO delivery_proofs (id, route_stop_id, recipient_name, recipient_dni, 
		                             recipient_relationship, signature_base64, photo_url, 
		                             failure_reason, notes, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (route_stop_id) DO UPDATE SET
			recipient_name = EXCLUDED.recipient_name,
			recipient_dni = EXCLUDED.recipient_dni,
			recipient_relationship = EXCLUDED.recipient_relationship,
			signature_base64 = EXCLUDED.signature_base64,
			photo_url = EXCLUDED.photo_url,
			failure_reason = EXCLUDED.failure_reason,
			notes = EXCLUDED.notes,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			submitted_at = NOW()
		RETURNING submitted_at`

	return r.db.QueryRow(
		query, pod.ID, pod.RouteStopID, pod.RecipientName, pod.RecipientDNI,
		pod.RecipientRelationship, pod.SignatureBase64, pod.PhotoURL,
		pod.FailureReason, pod.Notes, pod.Latitude, pod.Longitude,
	).Scan(&pod.SubmittedAt)
}

// ---------------- TELEMETRY & GPS ----------------

func (r *Repository) SaveTelemetry(t *models.TelemetryUpdate) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertQuery := `
		INSERT INTO vehicle_telemetry (vehicle_id, driver_id, latitude, longitude, speed, heading, battery_level)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := tx.Exec(insertQuery, t.VehicleID, t.DriverID, t.Latitude, t.Longitude, t.Speed, t.Heading, t.BatteryLevel); err != nil {
		return err
	}

	upsertLastPosQuery := `
		INSERT INTO vehicle_last_positions (vehicle_id, driver_id, latitude, longitude, speed, heading, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (vehicle_id) DO UPDATE SET
			driver_id = EXCLUDED.driver_id,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			speed = EXCLUDED.speed,
			heading = EXCLUDED.heading,
			updated_at = EXCLUDED.updated_at`
	if _, err := tx.Exec(upsertLastPosQuery, t.VehicleID, t.DriverID, t.Latitude, t.Longitude, t.Speed, t.Heading); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetAllVehicleLastPositions() ([]models.VehiclePosition, error) {
	query := `
		SELECT vlp.vehicle_id, v.plate, v.model, vlp.driver_id, d.name as driver_name,
		       vlp.latitude, vlp.longitude, vlp.speed, vlp.heading,
		       r.id as route_id, r.code as route_code, r.completed_stops + 1 as current_stop, 
		       r.total_stops, vlp.updated_at
		FROM vehicle_last_positions vlp
		INNER JOIN vehicles v ON vlp.vehicle_id = v.id
		LEFT JOIN drivers d ON vlp.driver_id = d.id
		LEFT JOIN routes r ON r.vehicle_id = v.id AND r.status = 'in_progress'
		ORDER BY vlp.updated_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.VehiclePosition
	for rows.Next() {
		var vp models.VehiclePosition
		if err := rows.Scan(
			&vp.VehicleID, &vp.Plate, &vp.Model, &vp.DriverID, &vp.DriverName,
			&vp.Latitude, &vp.Longitude, &vp.Speed, &vp.Heading,
			&vp.RouteID, &vp.RouteCode, &vp.CurrentStop, &vp.TotalStops, &vp.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, vp)
	}
	return list, nil
}

// ---------------- PUBLIC TRACKING ----------------

func (r *Repository) GetPublicTrackingData(trackingNumber string) (*models.PublicTrackingResponse, error) {
	query := `
		SELECT o.id, o.tracking_number, COALESCE(o.external_order_id, ''), o.customer_name,
		       COALESCE(o.customer_phone, ''), COALESCE(o.customer_email, ''), o.delivery_address,
		       COALESCE(o.commune, ''), o.latitude, o.longitude, o.window_start, o.window_end,
		       o.packages_count, COALESCE(o.items_description, ''), o.status, o.created_at,
		       COALESCE(rs.status, 'pending') as stop_status, COALESCE(rs.stop_number, 0),
		       COALESCE(r.total_stops, 0), d.name as driver_name, v.plate as vehicle_plate,
		       vlp.latitude as vehicle_lat, vlp.longitude as vehicle_lng,
		       dp.id as dp_id, dp.recipient_name, dp.recipient_dni, dp.recipient_relationship,
		       COALESCE(dp.signature_base64, ''), COALESCE(dp.photo_url, ''),
		       COALESCE(dp.failure_reason, ''), COALESCE(dp.notes, ''), dp.submitted_at
		FROM orders o
		LEFT JOIN route_stops rs ON rs.order_id = o.id
		LEFT JOIN routes r ON rs.route_id = r.id
		LEFT JOIN drivers d ON r.driver_id = d.id
		LEFT JOIN vehicles v ON r.vehicle_id = v.id
		LEFT JOIN vehicle_last_positions vlp ON vlp.vehicle_id = v.id
		LEFT JOIN delivery_proofs dp ON dp.route_stop_id = rs.id
		WHERE o.tracking_number = $1`

	var res models.PublicTrackingResponse
	var o models.Order
	var dpID, dpName, dpDNI, dpRel, dpSig, dpPhoto, dpFail, dpNotes *string
	var dpSubAt *time.Time

	err := r.db.QueryRow(query, trackingNumber).Scan(
		&o.ID, &o.TrackingNumber, &o.ExternalOrderID, &o.CustomerName,
		&o.CustomerPhone, &o.CustomerEmail, &o.DeliveryAddress,
		&o.Commune, &o.Latitude, &o.Longitude, &o.WindowStart, &o.WindowEnd,
		&o.PackagesCount, &o.ItemsDescription, &o.Status, &o.CreatedAt,
		&res.StopStatus, &res.StopNumber, &res.TotalStops,
		&res.DriverName, &res.VehiclePlate, &res.VehicleLat, &res.VehicleLng,
		&dpID, &dpName, &dpDNI, &dpRel, &dpSig, &dpPhoto, &dpFail, &dpNotes, &dpSubAt,
	)
	if err != nil {
		return nil, err
	}

	res.Order = o
	if dpID != nil {
		res.Proof = &models.DeliveryProof{
			ID:                    *dpID,
			RecipientName:         stringValue(dpName),
			RecipientDNI:          stringValue(dpDNI),
			RecipientRelationship: stringValue(dpRel),
			SignatureBase64:       stringValue(dpSig),
			PhotoURL:              stringValue(dpPhoto),
			FailureReason:         stringValue(dpFail),
			Notes:                 stringValue(dpNotes),
			SubmittedAt:           timeValue(dpSubAt),
		}
	}

	return &res, nil
}

// ---------------- DRIVERS & VEHICLES ----------------

func (r *Repository) GetAllDrivers() ([]models.Driver, error) {
	rows, err := r.db.Query(`SELECT id, name, email, COALESCE(phone, ''), COALESCE(license_number, ''), status, created_at FROM drivers ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Driver
	for rows.Next() {
		var d models.Driver
		if err := rows.Scan(&d.ID, &d.Name, &d.Email, &d.Phone, &d.LicenseNumber, &d.Status, &d.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *Repository) GetAllVehicles() ([]models.Vehicle, error) {
	rows, err := r.db.Query(`SELECT id, plate, COALESCE(model, ''), COALESCE(brand, ''), capacity_kg, status, created_at FROM vehicles ORDER BY plate ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		if err := rows.Scan(&v.ID, &v.Plate, &v.Model, &v.Brand, &v.CapacityKg, &v.Status, &v.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}

// ---------------- USERS & RBAC ----------------

func (r *Repository) GetAllUsers() ([]models.User, error) {
	query := `
		SELECT u.id, u.driver_id, u.name, u.email, u.role, u.vehicle_id, 
		       v.plate as vehicle_plate, v.model as vehicle_model, 
		       COALESCE(u.phone, ''), COALESCE(u.status, 'active'), u.created_at
		FROM users u
		LEFT JOIN vehicles v ON u.vehicle_id = v.id
		ORDER BY u.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.DriverID, &u.Name, &u.Email, &u.Role, &u.VehicleID,
			&u.VehiclePlate, &u.VehicleModel, &u.Phone, &u.Status, &u.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *Repository) CreateUser(u *models.User) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.Status == "" {
		u.Status = "active"
	}
	if u.Role != "driver" {
		u.VehicleID = nil
		u.DriverID = nil
	}

	query := `
		INSERT INTO users (id, driver_id, name, email, role, vehicle_id, phone, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at`

	return r.db.QueryRow(
		query, u.ID, u.DriverID, u.Name, u.Email, u.Role, u.VehicleID, u.Phone, u.Status,
	).Scan(&u.CreatedAt)
}

func (r *Repository) UpdateUser(u *models.User) error {
	if u.Role != "driver" {
		u.VehicleID = nil
		u.DriverID = nil
	}

	query := `
		UPDATE users
		SET driver_id = $1, name = $2, email = $3, role = $4, vehicle_id = $5, phone = $6, status = $7
		WHERE id = $8`

	_, err := r.db.Exec(query, u.DriverID, u.Name, u.Email, u.Role, u.VehicleID, u.Phone, u.Status, u.ID)
	return err
}

func (r *Repository) DeleteUser(id string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *Repository) GetUserByEmailWithPassword(email string) (*models.UserWithPassword, error) {
	query := `
		SELECT u.id, u.driver_id, u.name, u.email, COALESCE(u.password_hash, ''), u.role, u.vehicle_id,
		       v.plate as vehicle_plate, v.model as vehicle_model,
		       COALESCE(u.phone, ''), COALESCE(u.status, 'active'), u.created_at
		FROM users u
		LEFT JOIN vehicles v ON u.vehicle_id = v.id
		WHERE LOWER(u.email) = LOWER($1)`

	var u models.UserWithPassword
	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.DriverID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.VehicleID,
		&u.VehiclePlate, &u.VehicleModel, &u.Phone, &u.Status, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT u.id, u.driver_id, u.name, u.email, u.role, u.vehicle_id,
		       v.plate as vehicle_plate, v.model as vehicle_model,
		       COALESCE(u.phone, ''), COALESCE(u.status, 'active'), u.created_at
		FROM users u
		LEFT JOIN vehicles v ON u.vehicle_id = v.id
		WHERE u.id = $1`

	var u models.User
	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.DriverID, &u.Name, &u.Email, &u.Role, &u.VehicleID,
		&u.VehiclePlate, &u.VehicleModel, &u.Phone, &u.Status, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetUserTenants(userID string) ([]models.UserTenantMembership, error) {
	query := `
		SELECT tu.tenant_id, t.name, t.legal_id, t.type, tu.role, tu.status
		FROM tenant_users tu
		JOIN tenants t ON tu.tenant_id = t.id
		WHERE tu.user_id = $1 AND tu.status = 'active' AND t.status = 'active'
		ORDER BY t.name ASC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.UserTenantMembership
	for rows.Next() {
		var m models.UserTenantMembership
		if err := rows.Scan(
			&m.TenantID, &m.TenantName, &m.LegalID, &m.TenantType, &m.Role, &m.Status,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// GetAssignedRoutesByDriver obtiene las rutas asignadas a un conductor validando que el tenant_driver y el tenant estén activos
func (r *Repository) GetAssignedRoutesByDriver(driverID string) ([]models.AssignedRoute, error) {
	query := `
		SELECT r.id, r.code, r.date::text, r.status, r.carrier_tenant_id, t.name as carrier_name,
		       r.vehicle_id, v.plate as vehicle_plate, r.total_stops, r.completed_stops,
		       r.started_at, r.completed_at
		FROM routes r
		JOIN tenants t ON r.carrier_tenant_id = t.id
		JOIN tenant_drivers td ON td.tenant_id = r.carrier_tenant_id AND td.driver_id = r.driver_id
		LEFT JOIN vehicles v ON r.vehicle_id = v.id
		WHERE r.driver_id = $1
		  AND td.status = 'active'
		  AND t.status = 'active'
		ORDER BY r.date DESC, r.created_at DESC`

	rows, err := r.db.Query(query, driverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []models.AssignedRoute
	for rows.Next() {
		var ar models.AssignedRoute
		if err := rows.Scan(
			&ar.ID, &ar.Code, &ar.Date, &ar.Status, &ar.CarrierTenantID, &ar.CarrierName,
			&ar.VehicleID, &ar.VehiclePlate, &ar.TotalStops, &ar.CompletedStops,
			&ar.StartedAt, &ar.CompletedAt,
		); err != nil {
			return nil, err
		}
		routes = append(routes, ar)
	}
	return routes, nil
}

// ResolveRouteTenant valida de forma segura que la ruta pertenezca al conductor autenticado y que tanto el tenant como la vinculación estén activos
func (r *Repository) ResolveRouteTenant(driverID, routeID string) (*models.TenantResolution, error) {
	query := `
		SELECT r.id, r.driver_id, r.carrier_tenant_id
		FROM routes r
		JOIN tenants t ON r.carrier_tenant_id = t.id
		JOIN tenant_drivers td ON td.tenant_id = r.carrier_tenant_id AND td.driver_id = r.driver_id
		WHERE r.id = $1
		  AND r.driver_id = $2
		  AND r.carrier_tenant_id IS NOT NULL
		  AND td.status = 'active'
		  AND t.status = 'active'`

	var tr models.TenantResolution
	err := r.db.QueryRow(query, routeID, driverID).Scan(&tr.RouteID, &tr.DriverID, &tr.CarrierTenantID)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

// CheckUserTenantMembership valida si un usuario de backoffice tiene membresía activa en un tenant activo y retorna el rol
func (r *Repository) CheckUserTenantMembership(userID, tenantID string) (string, error) {
	query := `
		SELECT tu.role
		FROM tenant_users tu
		JOIN tenants t ON tu.tenant_id = t.id
		WHERE tu.user_id = $1 
		  AND tu.tenant_id = $2 
		  AND tu.status = 'active' 
		  AND t.status = 'active'`

	var role string
	err := r.db.QueryRow(query, userID, tenantID).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

// CheckDriverTenantActive valida si un conductor está habilitado (active) en un tenant activo
func (r *Repository) CheckDriverTenantActive(driverID, tenantID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM tenant_drivers td
			JOIN tenants t ON td.tenant_id = t.id
			WHERE td.driver_id = $1 
			  AND td.tenant_id = $2 
			  AND td.status = 'active' 
			  AND t.status = 'active'
		)`

	var exists bool
	err := r.db.QueryRow(query, driverID, tenantID).Scan(&exists)
	return exists, err
}

func stringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func timeValue(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}
