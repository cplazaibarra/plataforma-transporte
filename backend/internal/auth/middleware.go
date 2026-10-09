package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserClaimsKey     = "user_claims"
	ContextUserIDKey         = "user_id"
	ContextTenantClaimsKey   = "tenant_claims"
	ContextTenantIDKey       = "tenant_id"
	ContextMembershipRoleKey = "membership_role"
	ContextDriverIDKey       = "driver_id"
)

// IdentityMiddleware valida el Bearer token en el header Authorization para el Global Identity Context
func IdentityMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cabecera Authorization requerida"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de autorización inválido, debe ser 'Bearer <token>'"})
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		claims, err := ValidateGlobalToken(tokenStr, secretKey)
		if err != nil {
			if err == ErrExpiredToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expirado"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido o no autorizado"})
			return
		}

		c.Set(ContextUserClaimsKey, claims)
		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}

// TenantContextMiddleware valida obligatoriamente un Tenant-Scoped Access Token en el header Authorization.
// Rechaza tokens globales de identidad, tokens expirados, manipulados o con firma inválida.
// Inyecta en el contexto de Gin: tenant_claims, tenant_id, user_id, membership_role y driver_id (si aplica).
func TenantContextMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cabecera Authorization requerida para contexto tenant"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de autorización inválido, debe ser 'Bearer <token>'"})
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		claims, err := ValidateTenantToken(tokenStr, secretKey)
		if err != nil {
			if err == ErrExpiredToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de tenant expirado"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de tenant inválido, manipulado o no autorizado"})
			return
		}

		if claims.TenantID == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "token no contiene un tenant_id válido"})
			return
		}

		c.Set(ContextTenantClaimsKey, claims)
		c.Set(ContextTenantIDKey, claims.TenantID)
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextMembershipRoleKey, claims.MembershipRole)
		if claims.DriverID != nil && *claims.DriverID != "" {
			c.Set(ContextDriverIDKey, *claims.DriverID)
		}
		c.Next()
	}
}

// TryExtractTenantContext inspecciona opcionalmente la cabecera Authorization para detectar
// si el cliente envió un Tenant-Scoped Access Token. Si es válido, inyecta las variables de contexto.
// Si no hay token o no es de tenant, permite continuar sin abortar la petición (útil para endpoints de compatibilidad transicional).
func TryExtractTenantContext(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenStr := strings.TrimSpace(parts[1])
				claims, err := ValidateTenantToken(tokenStr, secretKey)
				if err == nil && claims != nil && claims.TenantID != "" {
					c.Set(ContextTenantClaimsKey, claims)
					c.Set(ContextTenantIDKey, claims.TenantID)
					c.Set(ContextUserIDKey, claims.UserID)
					c.Set(ContextMembershipRoleKey, claims.MembershipRole)
					if claims.DriverID != nil && *claims.DriverID != "" {
						c.Set(ContextDriverIDKey, *claims.DriverID)
					}
				}
			}
		}
		c.Next()
	}
}
