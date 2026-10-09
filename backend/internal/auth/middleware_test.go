package auth

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestTenantContextMiddleware(t *testing.T) {
	secret := "test_secret_for_middleware_32bytes"
	userID := "f1000000-0000-0000-0000-000000000001"
	tenantID := "00000000-0000-0000-0000-000000000001"
	driverID := "d1000000-0000-0000-0000-000000000001"

	// Helper to setup test router
	setupRouter := func() *gin.Engine {
		r := gin.New()
		r.Use(TenantContextMiddleware(secret))
		r.GET("/test-tenant", func(c *gin.Context) {
			tID, _ := c.Get(ContextTenantIDKey)
			uID, _ := c.Get(ContextUserIDKey)
			mRole, _ := c.Get(ContextMembershipRoleKey)
			dID, _ := c.Get(ContextDriverIDKey)

			c.JSON(http.StatusOK, gin.H{
				"tenant_id":       tID,
				"user_id":         uID,
				"membership_role": mRole,
				"driver_id":       dID,
			})
		})
		return r
	}

	router := setupRouter()

	// 1. Caso Exitoso: Token Tenant válido con driver_id
	validTenantToken, err := GenerateTenantToken(userID, tenantID, &driverID, "test@lastmile.com", "Test User", "driver", "driver", 1*time.Hour, secret)
	if err != nil {
		t.Fatalf("Error generating tenant token: %v", err)
	}

	req, _ := http.NewRequest("GET", "/test-tenant", nil)
	req.Header.Set("Authorization", "Bearer "+validTenantToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Se esperaba 200 OK con token válido, se obtuvo %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), tenantID) || !strings.Contains(w.Body.String(), driverID) {
		t.Errorf("Respuesta no contiene datos esperados: %s", w.Body.String())
	}

	// 2. Caso Rechazo: Token de Identidad Global (token_type = 'global_identity') en endpoint de Tenant
	globalToken, err := GenerateGlobalToken(userID, "test@lastmile.com", "Test User", "admin", 1*time.Hour, secret)
	if err != nil {
		t.Fatalf("Error generating global token: %v", err)
	}

	reqGlobal, _ := http.NewRequest("GET", "/test-tenant", nil)
	reqGlobal.Header.Set("Authorization", "Bearer "+globalToken)
	wGlobal := httptest.NewRecorder()
	router.ServeHTTP(wGlobal, reqGlobal)

	if wGlobal.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba 401 Unauthorized para Global Token en Tenant Middleware, se obtuvo %d", wGlobal.Code)
	}

	// 3. Caso Rechazo: Petición sin Header Authorization
	reqNoAuth, _ := http.NewRequest("GET", "/test-tenant", nil)
	wNoAuth := httptest.NewRecorder()
	router.ServeHTTP(wNoAuth, reqNoAuth)

	if wNoAuth.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba 401 Unauthorized sin cabecera Authorization, se obtuvo %d", wNoAuth.Code)
	}

	// 4. Caso Rechazo: Token Expirado
	expiredToken, err := GenerateTenantToken(userID, tenantID, &driverID, "test@lastmile.com", "Test User", "driver", "driver", -1*time.Hour, secret)
	if err != nil {
		t.Fatalf("Error generating expired token: %v", err)
	}

	reqExpired, _ := http.NewRequest("GET", "/test-tenant", nil)
	reqExpired.Header.Set("Authorization", "Bearer "+expiredToken)
	wExpired := httptest.NewRecorder()
	router.ServeHTTP(wExpired, reqExpired)

	if wExpired.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba 401 Unauthorized para token expirado, se obtuvo %d", wExpired.Code)
	}

	// 5. Caso Rechazo: Token con firma adulterada o manipulado
	parts := strings.Split(validTenantToken, ".")
	decodedClaims, _ := base64.RawURLEncoding.DecodeString(parts[1])
	tamperedClaims := strings.Replace(string(decodedClaims), tenantID, "00000000-0000-0000-0000-999999999999", 1)
	tamperedToken := parts[0] + "." + base64.RawURLEncoding.EncodeToString([]byte(tamperedClaims)) + "." + parts[2]

	reqTampered, _ := http.NewRequest("GET", "/test-tenant", nil)
	reqTampered.Header.Set("Authorization", "Bearer "+tamperedToken)
	wTampered := httptest.NewRecorder()
	router.ServeHTTP(wTampered, reqTampered)

	if wTampered.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba 401 Unauthorized para token manipulado, se obtuvo %d", wTampered.Code)
	}
}
