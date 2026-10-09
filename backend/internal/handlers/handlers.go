package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"lastmile-backend/internal/auth"
	"lastmile-backend/internal/models"
	"lastmile-backend/internal/repository"
	"lastmile-backend/internal/websocket"
)

type Handler struct {
	repo      *repository.Repository
	hub       *websocket.Hub
	jwtSecret string
}

func NewHandler(repo *repository.Repository, hub *websocket.Hub, jwtSecret string) *Handler {
	return &Handler{
		repo:      repo,
		hub:       hub,
		jwtSecret: jwtSecret,
	}
}

// ---------------- ORDERS ----------------

func (h *Handler) GetOrders(c *gin.Context) {
	if tenantID, exists := c.Get(auth.ContextTenantIDKey); exists && tenantID != "" {
		orders, err := h.repo.GetOrdersByTenant(tenantID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al consultar órdenes: " + err.Error()})
			return
		}
		if orders == nil {
			orders = []models.Order{}
		}
		c.JSON(http.StatusOK, orders)
		return
	}

	// Modo legacy / transicional (sin contexto de tenant autenticado)
	orders, err := h.repo.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if orders == nil {
		orders = []models.Order{}
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetOrderDetail(c *gin.Context) {
	id := c.Param("id")
	if tenantID, exists := c.Get(auth.ContextTenantIDKey); exists && tenantID != "" {
		o, err := h.repo.GetOrderByIDAndTenant(id, tenantID.(string))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "orden no encontrada o no autorizada para el tenant activo"})
			return
		}
		c.JSON(http.StatusOK, o)
		return
	}

	// Modo legacy / transicional
	o, err := h.repo.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var o models.Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tenantID, exists := c.Get(auth.ContextTenantIDKey); exists && tenantID != "" {
		// En contexto tenant: origin_tenant_id se toma obligatoria e inmutablemente del token autenticado.
		// Cualquier valor provisto por el cliente para origin_tenant_id es ignorado/sobreescrito por el servidor.
		if err := h.repo.CreateOrderWithTenant(&o, tenantID.(string)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error al crear orden tenantizada: " + err.Error()})
			return
		}

		h.hub.BroadcastEvent("order_created", o)
		c.JSON(http.StatusCreated, o)
		return
	}

	// Modo legacy / transicional
	if err := h.repo.CreateOrder(&o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.BroadcastEvent("order_created", o)
	c.JSON(http.StatusCreated, o)
}

type BatchOrdersInput struct {
	Orders []models.Order `json:"orders"`
}

func (h *Handler) ImportOrders(c *gin.Context) {
	var input BatchOrdersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var created []models.Order
	for _, o := range input.Orders {
		orderCopy := o
		if err := h.repo.CreateOrder(&orderCopy); err == nil {
			created = append(created, orderCopy)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Órdenes importadas con éxito",
		"total":   len(created),
		"orders":  created,
	})
}

// ---------------- ROUTES ----------------

func (h *Handler) GetRoutes(c *gin.Context) {
	if tenantID, exists := c.Get(auth.ContextTenantIDKey); exists && tenantID != "" {
		routes, err := h.repo.GetRoutesByTenant(tenantID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al consultar rutas: " + err.Error()})
			return
		}
		if routes == nil {
			routes = []models.Route{}
		}
		c.JSON(http.StatusOK, routes)
		return
	}

	// Modo legacy / transicional
	routes, err := h.repo.GetAllRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if routes == nil {
		routes = []models.Route{}
	}
	c.JSON(http.StatusOK, routes)
}

func (h *Handler) GetRouteDetail(c *gin.Context) {
	id := c.Param("id")
	if tenantID, exists := c.Get(auth.ContextTenantIDKey); exists && tenantID != "" {
		rt, err := h.repo.GetRouteWithStopsByTenant(id, tenantID.(string))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ruta no encontrada o no autorizada para el tenant activo"})
			return
		}
		c.JSON(http.StatusOK, rt)
		return
	}

	// Modo legacy / transicional
	rt, err := h.repo.GetRouteWithStops(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
		return
	}
	c.JSON(http.StatusOK, rt)
}

type CreateRouteInput struct {
	Code      string   `json:"code"`
	DriverID  *string  `json:"driver_id"`
	VehicleID *string  `json:"vehicle_id"`
	OrderIDs  []string `json:"order_ids"`
}

func (h *Handler) CreateRoute(c *gin.Context) {
	var input CreateRouteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rt := models.Route{
		Code:      input.Code,
		DriverID:  input.DriverID,
		VehicleID: input.VehicleID,
	}

	if tenantID, exists := c.Get(auth.ContextTenantIDKey); exists && tenantID != "" {
		if err := h.repo.CreateRouteWithTenant(&rt, input.OrderIDs, tenantID.(string)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error al crear ruta con aislamiento de tenant: " + err.Error()})
			return
		}
		h.hub.BroadcastEvent("route_created", rt)
		c.JSON(http.StatusCreated, rt)
		return
	}

	// Modo legacy / transicional
	if err := h.repo.CreateRoute(&rt, input.OrderIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.BroadcastEvent("route_created", rt)
	c.JSON(http.StatusCreated, rt)
}

type UpdateStopStatusInput struct {
	Status string `json:"status" binding:"required"` // pending, en_route, arrived, delivered, failed
}

func (h *Handler) UpdateStopStatus(c *gin.Context) {
	stopID := c.Param("stopId")
	var input UpdateStopStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tenantID, exists := c.Get(auth.ContextTenantIDKey); exists && tenantID != "" {
		if err := h.repo.UpdateStopStatusWithTenant(stopID, input.Status, tenantID.(string)); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "error al actualizar parada en tenant activo: " + err.Error()})
			return
		}

		h.hub.BroadcastEvent("stop_status_changed", gin.H{
			"stop_id": stopID,
			"status":  input.Status,
		})

		c.JSON(http.StatusOK, gin.H{"message": "Estado de la parada actualizado", "status": input.Status})
		return
	}

	// Modo legacy / transicional
	if err := h.repo.UpdateStopStatus(stopID, input.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.BroadcastEvent("stop_status_changed", gin.H{
		"stop_id": stopID,
		"status":  input.Status,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Estado de la parada actualizado", "status": input.Status})
}

// ---------------- PROOF OF DELIVERY (POD) ----------------

type SubmitPODInput struct {
	RouteStopID           string   `json:"route_stop_id" binding:"required"`
	RecipientName         string   `json:"recipient_name"`
	RecipientDNI          string   `json:"recipient_dni"`
	RecipientRelationship string   `json:"recipient_relationship"`
	SignatureBase64       string   `json:"signature_base64"`
	PhotoURL              string   `json:"photo_url"`
	FailureReason         string   `json:"failure_reason"`
	Notes                 string   `json:"notes"`
	Latitude              *float64 `json:"latitude"`
	Longitude             *float64 `json:"longitude"`
	Status                string   `json:"status" binding:"required"` // delivered or failed
}

func (h *Handler) SubmitPOD(c *gin.Context) {
	var input SubmitPODInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pod := models.DeliveryProof{
		RouteStopID:           input.RouteStopID,
		RecipientName:         input.RecipientName,
		RecipientDNI:          input.RecipientDNI,
		RecipientRelationship: input.RecipientRelationship,
		SignatureBase64:       input.SignatureBase64,
		PhotoURL:              input.PhotoURL,
		FailureReason:         input.FailureReason,
		Notes:                 input.Notes,
		Latitude:              input.Latitude,
		Longitude:             input.Longitude,
	}

	// 1. Guardar la prueba de entrega en la BD
	if err := h.repo.SaveProofOfDelivery(&pod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el comprobante: " + err.Error()})
		return
	}

	// 2. Actualizar estado de la parada (delivered o failed)
	finalStatus := input.Status
	if finalStatus != "delivered" && finalStatus != "failed" {
		finalStatus = "delivered"
	}
	if err := h.repo.UpdateStopStatus(input.RouteStopID, finalStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar estado de la parada: " + err.Error()})
		return
	}

	// 3. Emitir evento por WebSocket para actualizar instantáneamente la Torre de Control y el cliente
	h.hub.BroadcastEvent("pod_submitted", gin.H{
		"route_stop_id": input.RouteStopID,
		"status":        finalStatus,
		"proof":         pod,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Comprobante de entrega registrado exitosamente",
		"proof":   pod,
		"status":  finalStatus,
	})
}

// ---------------- TELEMETRY & GPS ----------------

func (h *Handler) SubmitTelemetry(c *gin.Context) {
	var t models.TelemetryUpdate
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.SaveTelemetry(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Emitir por WebSocket a la Torre de Control
	h.hub.BroadcastEvent("telemetry_update", t)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) GetVehiclePositions(c *gin.Context) {
	positions, err := h.repo.GetAllVehicleLastPositions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, positions)
}

// ---------------- PUBLIC TRACKING ----------------

func (h *Handler) GetPublicTracking(c *gin.Context) {
	trackingNumber := c.Param("trackingNumber")
	data, err := h.repo.GetPublicTrackingData(trackingNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Número de seguimiento no encontrado"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ---------------- CATALOGS ----------------

func (h *Handler) GetDrivers(c *gin.Context) {
	drivers, err := h.repo.GetAllDrivers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, drivers)
}

func (h *Handler) GetVehicles(c *gin.Context) {
	vehicles, err := h.repo.GetAllVehicles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vehicles)
}

// ---------------- USERS & RBAC ----------------

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if u.Role != "admin" && u.Role != "operator" && u.Role != "driver" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rol inválido. Debe ser 'admin', 'operator' o 'driver'"})
		return
	}

	if err := h.repo.CreateUser(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario: " + err.Error()})
		return
	}

	h.hub.BroadcastEvent("user_created", u)
	c.JSON(http.StatusCreated, u)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.ID = id

	if err := h.repo.UpdateUser(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario: " + err.Error()})
		return
	}

	h.hub.BroadcastEvent("user_updated", u)
	c.JSON(http.StatusOK, u)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario: " + err.Error()})
		return
	}

	h.hub.BroadcastEvent("user_deleted", gin.H{"id": id})
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}

// ---------------- AUTH & IDENTITY ----------------

// Login autentica las credenciales del usuario y emite un Global Identity Access Token
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credenciales requeridas: " + err.Error()})
		return
	}

	userWithPass, err := h.repo.GetUserByEmailWithPassword(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	if userWithPass.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"error": "usuario inactivo o suspendido"})
		return
	}

	if !auth.CheckPassword(req.Password, userWithPass.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	ttl := 24 * time.Hour
	tokenStr, err := auth.GenerateGlobalToken(
		userWithPass.ID,
		userWithPass.Email,
		userWithPass.Name,
		userWithPass.Role,
		ttl,
		h.jwtSecret,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar token de sesión"})
		return
	}

	resp := models.LoginResponse{
		Token:     tokenStr,
		TokenType: "Bearer",
		ExpiresIn: int64(ttl.Seconds()),
		User:      userWithPass.User,
	}

	c.JSON(http.StatusOK, resp)
}

// GetMe devuelve la identidad del usuario autenticado a partir del Global Identity Context
func (h *Handler) GetMe(c *gin.Context) {
	claimsVal, exists := c.Get(auth.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	claims, ok := claimsVal.(*auth.GlobalClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims no válidos"})
		return
	}

	user, err := h.repo.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetMyTenants devuelve la lista de tenants a los que pertenece el usuario autenticado
func (h *Handler) GetMyTenants(c *gin.Context) {
	claimsVal, exists := c.Get(auth.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	claims, ok := claimsVal.(*auth.GlobalClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims no válidos"})
		return
	}

	tenants, err := h.repo.GetUserTenants(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al consultar membresías: " + err.Error()})
		return
	}

	if tenants == nil {
		tenants = []models.UserTenantMembership{}
	}

	c.JSON(http.StatusOK, tenants)
}

// GetAssignedRoutes devuelve exclusivamente las rutas asignadas al conductor autenticado validando pertenencia y estado activo
func (h *Handler) GetAssignedRoutes(c *gin.Context) {
	claimsVal, exists := c.Get(auth.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	claims, ok := claimsVal.(*auth.GlobalClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims no válidos"})
		return
	}

	// 1. Validar que el rol sea 'driver'
	if claims.Role != "driver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo usuarios conductores pueden consultar rutas asignadas"})
		return
	}

	// 2. Obtener el usuario actual para verificar su driver_id
	user, err := h.repo.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	if user.DriverID == nil || *user.DriverID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "el usuario no tiene un perfil de conductor vinculado"})
		return
	}

	// 3. Consultar rutas asignadas con filtro de tenant_drivers y tenant activo
	routes, err := h.repo.GetAssignedRoutesByDriver(*user.DriverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al consultar rutas asignadas: " + err.Error()})
		return
	}

	if routes == nil {
		routes = []models.AssignedRoute{}
	}

	c.JSON(http.StatusOK, routes)
}

// ResolveRouteTenant valida de forma segura que la ruta pertenezca al conductor autenticado y esté habilitada
func (h *Handler) ResolveRouteTenant(user *models.User, routeID string) (*models.TenantResolution, error) {
	if user == nil || user.Role != "driver" || user.DriverID == nil || *user.DriverID == "" {
		return nil, auth.ErrInvalidToken
	}
	return h.repo.ResolveRouteTenant(*user.DriverID, routeID)
}

// SelectTenant valida la pertenencia y autorización del usuario para el tenant solicitado y emite un Tenant-Scoped Access Token
func (h *Handler) SelectTenant(c *gin.Context) {
	claimsVal, exists := c.Get(auth.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	claims, ok := claimsVal.(*auth.GlobalClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims no válidos"})
		return
	}

	var req models.SelectTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.TenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id es requerido"})
		return
	}

	user, err := h.repo.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	var membershipRole string
	var selectedTenantID = req.TenantID

	if user.Role == "admin" || user.Role == "operator" {
		// 1. Caso Backoffice: comprobar tenant_users con tenant activo
		role, err := h.repo.CheckUserTenantMembership(user.ID, req.TenantID)
		if err != nil || role == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "no posee membresía activa en el tenant solicitado"})
			return
		}
		membershipRole = role
	} else if user.Role == "driver" {
		// 2. Caso Conductor
		if user.DriverID == nil || *user.DriverID == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "el usuario no tiene un perfil de conductor vinculado"})
			return
		}

		// Si se especifica una ruta asignada, resolver a través de ResolveRouteTenant
		if req.RouteID != nil && *req.RouteID != "" {
			resolution, err := h.ResolveRouteTenant(user, *req.RouteID)
			if err != nil || resolution == nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "la ruta especificada no está autorizada o el tenant no está activo"})
				return
			}
			if resolution.CarrierTenantID != req.TenantID {
				c.JSON(http.StatusForbidden, gin.H{"error": "el tenant solicitado no coincide con la empresa transportista de la ruta"})
				return
			}
			selectedTenantID = resolution.CarrierTenantID
		} else {
			// Si no se especifica ruta, verificar que el conductor esté active en tenant_drivers para ese tenant
			active, err := h.repo.CheckDriverTenantActive(*user.DriverID, req.TenantID)
			if err != nil || !active {
				c.JSON(http.StatusForbidden, gin.H{"error": "el conductor no está habilitado en el tenant solicitado"})
				return
			}
		}
		membershipRole = "driver"
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "rol de usuario no autorizado"})
		return
	}

	// Emitir Tenant-Scoped Access Token
	ttl := 12 * time.Hour
	tokenStr, err := auth.GenerateTenantToken(
		user.ID,
		selectedTenantID,
		user.DriverID,
		user.Email,
		user.Name,
		user.Role,
		membershipRole,
		ttl,
		h.jwtSecret,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar token de tenant"})
		return
	}

	resp := models.SelectTenantResponse{
		Token:     tokenStr,
		TokenType: "Bearer",
		ExpiresIn: int64(ttl.Seconds()),
		TenantID:  selectedTenantID,
		Role:      membershipRole,
	}

	c.JSON(http.StatusOK, resp)
}

