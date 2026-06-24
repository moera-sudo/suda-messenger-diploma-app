package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager struct {
	signingKey string
}

func NewManager(signingKey string) *TokenManager {
	return &TokenManager{signingKey: signingKey}
} 

func (m *TokenManager) NewJWT(userID uuid.UUID, role string, tokenType string, ttl time.Duration) (string, error) {
	return m.NewJWTWithSession(userID, role, tokenType, uuid.Nil, ttl)
}

// NewJWTWithSession mints a token that additionally carries the originating
// session id as the "sid" claim. The gateway forwards this id downstream so the
// service can mark the current session in the Active Sessions list. A Nil
// sessionID omits the claim (used by stateless tokens like webapp-session).
func (m *TokenManager) NewJWTWithSession(userID uuid.UUID, role string, tokenType string, sessionID uuid.UUID, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub" : userID.String(),
		"role" : role,
		"type" : tokenType,
		"exp" : time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}
	if sessionID != uuid.Nil {
		claims["sid"] = sessionID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.signingKey))
}

func (m *TokenManager) ParseRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	} )
	
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf("Invalid token")
	}

	if typ, ok := claims["type"].(string); !ok || typ != "refresh" {
		return uuid.Nil, fmt.Errorf("Invalid token type")	
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("Token does not contain user id")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid user id format")
	}

	_, ok = claims["role"].(string)
	if !ok {
		return userID, nil
	}

	return userID, nil

}

