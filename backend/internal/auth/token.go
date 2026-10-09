package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken = errors.New("token inválido o manipulado")
	ErrExpiredToken = errors.New("token expirado")
	ErrMalformed    = errors.New("formato de token malformado")
)

type TokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// GlobalClaims representa el Global Identity Context (usuario autenticado a nivel de plataforma)
type GlobalClaims struct {
	UserID    string `json:"sub"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "global_identity"
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// CheckPassword valida una contraseña en texto plano contra el hash almacenado.
// Soporta contraseñas hasheadas con bcrypt y, por compatibilidad con el seed inicial, texto plano "default_secret".
func CheckPassword(plainPassword, storedHash string) bool {
	if strings.HasPrefix(storedHash, "$2a$") || strings.HasPrefix(storedHash, "$2b$") {
		err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(plainPassword))
		return err == nil
	}
	// Comparación en tiempo constante para contraseñas de transición legacy
	return hmac.Equal([]byte(plainPassword), []byte(storedHash))
}

// HashPassword genera un hash bcrypt seguro para contraseñas nuevas o actualizadas
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// GenerateGlobalToken emite un token JWT HS256 para el contexto global de identidad
func GenerateGlobalToken(userID, email, name, role string, ttl time.Duration, secretKey string) (string, error) {
	now := time.Now()
	header := TokenHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	claims := GlobalClaims{
		UserID:    userID,
		Email:     email,
		Name:      name,
		Role:      role,
		TokenType: "global_identity",
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(ttl).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("error encoding header: %w", err)
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("error encoding claims: %w", err)
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	payload := headerB64 + "." + claimsB64
	sig := createHMACSHA256(payload, secretKey)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return payload + "." + sigB64, nil
}

// ValidateGlobalToken verifica firma criptográfica y expiración, devolviendo los GlobalClaims
func ValidateGlobalToken(tokenStr, secretKey string) (*GlobalClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, ErrMalformed
	}

	headerB64, claimsB64, sigB64 := parts[0], parts[1], parts[2]
	payload := headerB64 + "." + claimsB64

	sigBytes, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return nil, ErrInvalidToken
	}

	expectedSig := createHMACSHA256(payload, secretKey)
	if !hmac.Equal(sigBytes, expectedSig) {
		return nil, ErrInvalidToken
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(claimsB64)
	if err != nil {
		return nil, ErrMalformed
	}

	var claims GlobalClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, ErrMalformed
	}

	if claims.TokenType != "global_identity" {
		return nil, ErrInvalidToken
	}

	now := time.Now().Unix()
	if claims.ExpiresAt <= now {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

// TenantClaims representa el Tenant-Scoped Access Context
type TenantClaims struct {
	UserID         string  `json:"sub"`
	TenantID       string  `json:"tenant_id"`
	DriverID       *string `json:"driver_id,omitempty"`
	Email          string  `json:"email"`
	Name           string  `json:"name"`
	Role           string  `json:"role"`
	MembershipRole string  `json:"membership_role"`
	TokenType      string  `json:"token_type"` // "tenant_access"
	IssuedAt       int64   `json:"iat"`
	ExpiresAt      int64   `json:"exp"`
}

// GenerateTenantToken emite un token JWT HS256 para el contexto específico de un tenant autorizado
func GenerateTenantToken(userID, tenantID string, driverID *string, email, name, role, membershipRole string, ttl time.Duration, secretKey string) (string, error) {
	now := time.Now()
	header := TokenHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	claims := TenantClaims{
		UserID:         userID,
		TenantID:       tenantID,
		DriverID:       driverID,
		Email:          email,
		Name:           name,
		Role:           role,
		MembershipRole: membershipRole,
		TokenType:      "tenant_access",
		IssuedAt:       now.Unix(),
		ExpiresAt:      now.Add(ttl).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("error encoding header: %w", err)
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("error encoding claims: %w", err)
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	payload := headerB64 + "." + claimsB64
	sig := createHMACSHA256(payload, secretKey)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return payload + "." + sigB64, nil
}

// ValidateTenantToken verifica firma criptográfica y expiración de un Tenant-Scoped Access Token
func ValidateTenantToken(tokenStr, secretKey string) (*TenantClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, ErrMalformed
	}

	headerB64, claimsB64, sigB64 := parts[0], parts[1], parts[2]
	payload := headerB64 + "." + claimsB64

	sigBytes, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return nil, ErrInvalidToken
	}

	expectedSig := createHMACSHA256(payload, secretKey)
	if !hmac.Equal(sigBytes, expectedSig) {
		return nil, ErrInvalidToken
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(claimsB64)
	if err != nil {
		return nil, ErrMalformed
	}

	var claims TenantClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, ErrMalformed
	}

	if claims.TokenType != "tenant_access" {
		return nil, ErrInvalidToken
	}

	now := time.Now().Unix()
	if claims.ExpiresAt <= now {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func createHMACSHA256(data, secretKey string) []byte {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return h.Sum(nil)
}
