package models

import (
	"time"
)

type Driver struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	LicenseNumber string    `json:"license_number"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
 
type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LegalID   string    `json:"legal_id"`
	Type      string    `json:"type"` // "carrier_3pl", "shipper_retailer"
	Status    string    `json:"status"` // "active", "inactive", "suspended"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TenantUser struct {
	TenantID  string    `json:"tenant_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"` // "admin", "operator"
	Status    string    `json:"status"` // "active", "inactive"
	CreatedAt time.Time `json:"created_at"`
}

type TenantDriver struct {
	TenantID  string    `json:"tenant_id"`
	DriverID  string    `json:"driver_id"`
	Status    string    `json:"status"` // "active", "blocked", "pending"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TenantVehicle struct {
	TenantID  string    `json:"tenant_id"`
	VehicleID string    `json:"vehicle_id"`
	Status    string    `json:"status"` // "active", "blocked", "pending"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID           string    `json:"id"`
	DriverID     *string   `json:"driver_id,omitempty"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"` // "admin", "operator", "driver"
	VehicleID    *string   `json:"vehicle_id,omitempty"`
	VehiclePlate *string   `json:"vehicle_plate,omitempty"`
	VehicleModel *string   `json:"vehicle_model,omitempty"`
	Phone        string    `json:"phone"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserWithPassword struct {
	User
	PasswordHash string `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
	User      User   `json:"user"`
}

type SelectTenantRequest struct {
	TenantID string  `json:"tenant_id"`
	RouteID  *string `json:"route_id,omitempty"`
}

type SelectTenantResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
	TenantID  string `json:"tenant_id"`
	Role      string `json:"role"`
}

type UserTenantMembership struct {
	TenantID   string `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
	LegalID    string `json:"legal_id"`
	TenantType string `json:"tenant_type"`
	Role       string `json:"role"`
	Status     string `json:"status"`
}

type AssignedRoute struct {
	ID              string     `json:"id"`
	Code            string     `json:"code"`
	Date            string     `json:"date"`
	Status          string     `json:"status"`
	CarrierTenantID string     `json:"carrier_tenant_id"`
	CarrierName     string     `json:"carrier_name"`
	VehicleID       *string    `json:"vehicle_id,omitempty"`
	VehiclePlate    *string    `json:"vehicle_plate,omitempty"`
	TotalStops      int        `json:"total_stops"`
	CompletedStops  int        `json:"completed_stops"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

type TenantResolution struct {
	RouteID         string `json:"route_id"`
	DriverID        string `json:"driver_id"`
	CarrierTenantID string `json:"carrier_tenant_id"`
}

type Vehicle struct {
	ID          string    `json:"id"`
	Plate       string    `json:"plate"`
	Model       string    `json:"model"`
	Brand       string    `json:"brand"`
	CapacityKg  float64   `json:"capacity_kg"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type Order struct {
	ID               string    `json:"id"`
	OriginTenantID   *string   `json:"origin_tenant_id,omitempty"`
	CarrierTenantID  *string   `json:"carrier_tenant_id,omitempty"`
	TrackingNumber   string    `json:"tracking_number"`
	ExternalOrderID  string    `json:"external_order_id"`
	CustomerName     string    `json:"customer_name"`
	CustomerPhone    string    `json:"customer_phone"`
	CustomerEmail    string    `json:"customer_email"`
	DeliveryAddress  string    `json:"delivery_address"`
	Commune          string    `json:"commune"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	WindowStart      string    `json:"window_start"`
	WindowEnd        string    `json:"window_end"`
	PackagesCount    int       `json:"packages_count"`
	ItemsDescription string    `json:"items_description"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type Route struct {
	ID              string       `json:"id"`
	CarrierTenantID *string      `json:"carrier_tenant_id,omitempty"`
	Code            string       `json:"code"`
	Date            string       `json:"date"`
	DriverID        *string      `json:"driver_id"`
	DriverName      *string      `json:"driver_name,omitempty"`
	VehicleID       *string      `json:"vehicle_id"`
	VehiclePlate    *string      `json:"vehicle_plate,omitempty"`
	Status          string       `json:"status"`
	StartedAt       *time.Time   `json:"started_at,omitempty"`
	CompletedAt     *time.Time   `json:"completed_at,omitempty"`
	TotalStops      int          `json:"total_stops"`
	CompletedStops  int          `json:"completed_stops"`
	Stops           []RouteStop  `json:"stops,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
}

type RouteStop struct {
	ID          string         `json:"id"`
	RouteID     string         `json:"route_id"`
	OrderID     string         `json:"order_id"`
	StopNumber  int            `json:"stop_number"`
	Status      string         `json:"status"` // pending, en_route, arrived, delivered, failed
	ETA         *time.Time     `json:"eta,omitempty"`
	ArrivedAt   *time.Time     `json:"arrived_at,omitempty"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	Notes       string         `json:"notes"`
	Order       *Order         `json:"order,omitempty"`
	Proof       *DeliveryProof `json:"proof,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

type DeliveryProof struct {
	ID                    string     `json:"id"`
	RouteStopID           string     `json:"route_stop_id"`
	RecipientName         string     `json:"recipient_name"`
	RecipientDNI          string     `json:"recipient_dni"`
	RecipientRelationship string     `json:"recipient_relationship"`
	SignatureBase64       string     `json:"signature_base64,omitempty"`
	PhotoURL              string     `json:"photo_url,omitempty"`
	FailureReason         string     `json:"failure_reason,omitempty"`
	Notes                 string     `json:"notes"`
	Latitude              *float64   `json:"latitude,omitempty"`
	Longitude             *float64   `json:"longitude,omitempty"`
	SubmittedAt           time.Time  `json:"submitted_at"`
}

type TelemetryUpdate struct {
	VehicleID    string    `json:"vehicle_id"`
	DriverID     *string   `json:"driver_id,omitempty"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Speed        float64   `json:"speed"`
	Heading      float64   `json:"heading"`
	BatteryLevel *int      `json:"battery_level,omitempty"`
	RecordedAt   time.Time `json:"recorded_at"`
}

type VehiclePosition struct {
	VehicleID    string    `json:"vehicle_id"`
	Plate        string    `json:"plate"`
	Model        string    `json:"model"`
	DriverID     *string   `json:"driver_id"`
	DriverName   *string   `json:"driver_name"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Speed        float64   `json:"speed"`
	Heading      float64   `json:"heading"`
	RouteID      *string   `json:"route_id"`
	RouteCode    *string   `json:"route_code"`
	CurrentStop  *int      `json:"current_stop"`
	TotalStops   *int      `json:"total_stops"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PublicTrackingResponse struct {
	Order         Order            `json:"order"`
	StopStatus    string           `json:"stop_status"`
	StopNumber    int              `json:"stop_number"`
	TotalStops    int              `json:"total_stops"`
	DriverName    *string          `json:"driver_name,omitempty"`
	VehiclePlate  *string          `json:"vehicle_plate,omitempty"`
	VehicleLat    *float64         `json:"vehicle_lat,omitempty"`
	VehicleLng    *float64         `json:"vehicle_lng,omitempty"`
	Proof         *DeliveryProof   `json:"proof,omitempty"`
}

type WSMessage struct {
	Type    string      `json:"type"` // "telemetry_update", "order_status_change", "pod_submitted"
	Payload interface{} `json:"payload"`
}
