package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds the configuration for JWT token generation and validation
type JWTConfig struct {
	SecretKey     string        // Secret key used for signing JWT tokens
	TokenDuration time.Duration // Duration for which the token will be valid
}

// NewJWTConfig creates a new JWT configuration instance
// Parameters:
//   - secretKey: The key used to sign the JWT tokens
//   - tokenDuration: How long the token will be valid for
//
// Returns:
//   - A pointer to the new JWTConfig instance
func NewJWTConfig(secretKey string, tokenDuration time.Duration) *JWTConfig {
	return &JWTConfig{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration,
	}
}

// GenerateToken creates a new JWT token with user information
// Parameters:
//   - userID: The unique identifier of the user
//   - username: The user's username
//   - role: The user's role in the system
//
// Returns:
//   - A signed JWT token string and any error that occurred during signing
func (c *JWTConfig) GenerateToken(userID, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(c.TokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.SecretKey))
}

// ValidateToken validates a JWT token string and returns the parsed token
// Parameters:
//   - tokenString: The JWT token string to validate
//
// Returns:
//   - The parsed JWT token and any error that occurred during validation
//   - Returns jwt.ErrSignatureInvalid if the signing method is invalid
func (c *JWTConfig) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(c.SecretKey), nil
	})
}
