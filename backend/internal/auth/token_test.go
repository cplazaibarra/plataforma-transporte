package auth

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"
)

func TestTokenGenerationAndValidation(t *testing.T) {
	secret := "test_secret_key_super_secure_32bytes"
	userID := "f1000000-0000-0000-0000-000000000001"
	email := "admin@lastmile.com"
	name := "Admin"
	role := "admin"

	// 1. Token válido
	token, err := GenerateGlobalToken(userID, email, name, role, 1*time.Hour, secret)
	if err != nil {
		t.Fatalf("Error al generar token: %v", err)
	}

	claims, err := ValidateGlobalToken(token, secret)
	if err != nil {
		t.Fatalf("Token debería ser válido, pero falló: %v", err)
	}

	if claims.UserID != userID || claims.Email != email || claims.Role != role {
		t.Errorf("Claims inconsistentes: %+v", claims)
	}

	// 2. Token con secreto incorrecto
	_, err = ValidateGlobalToken(token, "wrong_secret")
	if err != ErrInvalidToken {
		t.Errorf("Se esperaba ErrInvalidToken, se obtuvo: %v", err)
	}

	// 3. Token expirado
	expiredToken, err := GenerateGlobalToken(userID, email, name, role, -1*time.Second, secret)
	if err != nil {
		t.Fatalf("Error al generar token expirado: %v", err)
	}
	_, err = ValidateGlobalToken(expiredToken, secret)
	if err != ErrExpiredToken {
		t.Errorf("Se esperaba ErrExpiredToken, se obtuvo: %v", err)
	}

	// 4. Token malformado
	_, err = ValidateGlobalToken("invalido.token", secret)
	if err != ErrMalformed {
		t.Errorf("Se esperaba ErrMalformed, se obtuvo: %v", err)
	}
}

func TestPasswordCheck(t *testing.T) {
	plain := "my_secret_pass"
	hash, err := HashPassword(plain)
	if err != nil {
		t.Fatalf("Error al hashear password: %v", err)
	}

	if !CheckPassword(plain, hash) {
		t.Errorf("CheckPassword con hash bcrypt debería ser true")
	}

	if CheckPassword("wrong_pass", hash) {
		t.Errorf("CheckPassword con password incorrecta debería ser false")
	}

	// Compatibilidad legacy
	if !CheckPassword("default_secret", "default_secret") {
		t.Errorf("CheckPassword con legacy plain secret debería ser true")
	}
	if CheckPassword("wrong_secret", "default_secret") {
		t.Errorf("CheckPassword con legacy plain secret incorrecto debería ser false")
	}
}

func TestTenantTokenGenerationAndValidation(t *testing.T) {
	secret := "test_secret_key_super_secure_32bytes"
	userID := "f3000000-0000-0000-0000-000000000003"
	tenantID := "00000000-0000-0000-0000-000000000001"
	driverID := "d1000000-0000-0000-0000-000000000001"
	email := "cmorales@transporte.com"
	name := "Carlos Morales"
	role := "driver"
	membershipRole := "driver"

	// 1. Tenant-Scoped Token válido
	token, err := GenerateTenantToken(userID, tenantID, &driverID, email, name, role, membershipRole, 12*time.Hour, secret)
	if err != nil {
		t.Fatalf("Error al generar TenantToken: %v", err)
	}

	claims, err := ValidateTenantToken(token, secret)
	if err != nil {
		t.Fatalf("TenantToken debería ser válido, pero falló: %v", err)
	}

	if claims.UserID != userID || claims.TenantID != tenantID || *claims.DriverID != driverID || claims.TokenType != "tenant_access" {
		t.Errorf("Claims inconsistentes en TenantToken: %+v", claims)
	}

	// 2. Intentar validar TenantToken como GlobalToken debe fallar
	_, err = ValidateGlobalToken(token, secret)
	if err != ErrInvalidToken {
		t.Errorf("Un TenantToken NO debe ser aceptado como GlobalToken. Se obtuvo: %v", err)
	}

	// 3. Modificación del payload (manipulación de tenant_id) debe invalidar la firma
	parts := strings.Split(token, ".")
	if len(parts) == 3 {
		decodedClaims, _ := base64.RawURLEncoding.DecodeString(parts[1])
		tamperedClaimsJSON := strings.Replace(string(decodedClaims), tenantID, "ffffffff-ffff-ffff-ffff-ffffffffffff", 1)
		tamperedClaimsB64 := base64.RawURLEncoding.EncodeToString([]byte(tamperedClaimsJSON))
		tamperedToken := parts[0] + "." + tamperedClaimsB64 + "." + parts[2]
		_, err = ValidateTenantToken(tamperedToken, secret)
		if err != ErrInvalidToken {
			t.Errorf("Token con payload manipulado debería dar ErrInvalidToken, se obtuvo: %v", err)
		}
	}

	// 4. Token con firma incorrecta
	_, err = ValidateTenantToken(token, "otra_firma_secreta")
	if err != ErrInvalidToken {
		t.Errorf("Se esperaba ErrInvalidToken para firma incorrecta, se obtuvo: %v", err)
	}

	// 5. Token expirado
	expiredToken, err := GenerateTenantToken(userID, tenantID, &driverID, email, name, role, membershipRole, -1*time.Minute, secret)
	if err != nil {
		t.Fatalf("Error generando token expirado: %v", err)
	}
	_, err = ValidateTenantToken(expiredToken, secret)
	if err != ErrExpiredToken {
		t.Errorf("Se esperaba ErrExpiredToken, se obtuvo: %v", err)
	}
}
