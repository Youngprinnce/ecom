package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/youngprinnce/go-ecom/config"
)

// CreateJWT generates a new JWT token for the given user ID and role.
func CreateJWT(secret []byte, userID int, role string) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWT_EXPIRE_IN_SECONDS)
	// Create the JWT claims
	claims := jwt.MapClaims{
		"userID": userID,                            // Include user ID in the claims
		"role":   role,                              // Include user role in the claims
		"expiresAt":    time.Now().Add(expiration).Unix(), // Token expires in 7 days
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}
