package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lastmile-backend/internal/auth"
	"lastmile-backend/internal/config"
	"lastmile-backend/internal/database"
	"lastmile-backend/internal/handlers"
	"lastmile-backend/internal/repository"
	"lastmile-backend/internal/websocket"
)

func main() {
	cfg := config.LoadConfig()

	// 1. Conexión a Base de Datos PostgreSQL
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// 2. Repositorio y WebSocket Hub
	repo := repository.NewRepository(db)
	hub := websocket.NewHub()
	go hub.Run()

	// 3. Controladores
	h := handlers.NewHandler(repo, hub, cfg.JWTSecret)

	// 4. Servidor Gin HTTP
	r := gin.Default()

	// Configuración CORS permisiva para frontend en desarrollo
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "LastMile Go API"})
	})

	// WebSocket para telemetría y eventos en vivo
	r.GET("/ws", func(c *gin.Context) {
		hub.ServeWS(c.Writer, c.Request)
	})

	// API REST v1
	v1 := r.Group("/api/v1")
	{
		// Autenticación de Identidad Global
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", h.Login)
			authGroup.POST("/select-tenant", auth.IdentityMiddleware(cfg.JWTSecret), h.SelectTenant)
		}

		// Identidad y contexto global del usuario autenticado (Protegido por Global Identity Token)
		me := v1.Group("/me")
		me.Use(auth.IdentityMiddleware(cfg.JWTSecret))
		{
			me.GET("", h.GetMe)
			me.GET("/tenants", h.GetMyTenants)
			me.GET("/assigned-routes", h.GetAssignedRoutes)
		}

		// =========================================================================
		// RECURSOS OPERACIONALES MULTI-TENANT (FASE 3B)
		// Protegidos estrictamente por TenantContextMiddleware.
		// Requiere exclusivamente Tenant-Scoped Access Token (token_type = "tenant_access").
		// Aislamiento estricto de datos en runtime:
		// - Orders: origin_tenant_id = active_tenant_id OR carrier_tenant_id = active_tenant_id
		// - Routes: carrier_tenant_id = active_tenant_id
		// =========================================================================
		tenantOps := v1.Group("/tenant")
		tenantOps.Use(auth.TenantContextMiddleware(cfg.JWTSecret))
		{
			// Órdenes aisladas por tenant
			tenantOps.GET("/orders", h.GetOrders)
			tenantOps.POST("/orders", h.CreateOrder)
			tenantOps.GET("/orders/:id", h.GetOrderDetail)

			// Rutas aisladas por tenant transportista
			tenantOps.GET("/routes", h.GetRoutes)
			tenantOps.POST("/routes", h.CreateRoute)
			tenantOps.GET("/routes/:id", h.GetRouteDetail)
			tenantOps.PATCH("/routes/stops/:stopId/status", h.UpdateStopStatus)
		}

		// =========================================================================
		// ENDPOINTS OPERACIONALES TRANSICIONALES / LEGACY
		// NOTA DE ARQUITECTURA: Estos endpoints se mantienen de forma transicional
		// para preservar compatibilidad con frontend Next.js y PWA móvil antes de su
		// migración completa. Utilizan TryExtractTenantContext: si reciben un Bearer token
		// de tenant activo aíslan los recursos en runtime; si no lo reciben, operan
		// en modo legacy para no causar una ruptura abrupta del servicio visual.
		// =========================================================================
		opsLegacy := v1.Group("")
		opsLegacy.Use(auth.TryExtractTenantContext(cfg.JWTSecret))
		{
			// Órdenes / Guías
			opsLegacy.GET("/orders", h.GetOrders)
			opsLegacy.POST("/orders", h.CreateOrder)
			opsLegacy.GET("/orders/:id", h.GetOrderDetail)
			opsLegacy.POST("/orders/import", h.ImportOrders)

			// Rutas y Paradas
			opsLegacy.GET("/routes", h.GetRoutes)
			opsLegacy.POST("/routes", h.CreateRoute)
			opsLegacy.GET("/routes/:id", h.GetRouteDetail)
			opsLegacy.PATCH("/routes/stops/:stopId/status", h.UpdateStopStatus)
		}

		// Prueba de Entrega (POD - Proof of Delivery)
		v1.POST("/pod", h.SubmitPOD)

		// Telemetría GPS de los móviles
		v1.POST("/telemetry", h.SubmitTelemetry)
		v1.GET("/telemetry/positions", h.GetVehiclePositions)

		// Tracking público para el cliente final
		v1.GET("/tracking/:trackingNumber", h.GetPublicTracking)

		// Catálogos
		v1.GET("/drivers", h.GetDrivers)
		v1.GET("/vehicles", h.GetVehicles)

		// Gestión de Usuarios y Roles (RBAC)
		v1.GET("/users", h.GetUsers)
		v1.POST("/users", h.CreateUser)
		v1.PUT("/users/:id", h.UpdateUser)
		v1.DELETE("/users/:id", h.DeleteUser)
	}

	log.Printf("🚀 Servidor LastMile iniciado en puerto %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
