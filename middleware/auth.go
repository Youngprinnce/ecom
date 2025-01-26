package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/youngprinnce/go-ecom/config"
	"github.com/youngprinnce/go-ecom/utils"
)

// Define context keys for userID and role
type contextKey string

const (
	UserKey contextKey = "userID"
	RoleKey contextKey = "role"
)

// JWTAuth middleware validates the JWT token and sets the userID and role in the request context.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := utils.GetTokenFromRequest(c.Request)
		if tokenString == "" {
			utils.WriteError(c.Writer, http.StatusUnauthorized, fmt.Errorf("missing token"))
			c.Abort()
			return
		}

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Envs.JWT_SECRET), nil
		})
		if err != nil || !token.Valid {
			utils.WriteError(c.Writer, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			c.Abort()
			return
		}

		// Extract claims from the token
		claims := token.Claims.(jwt.MapClaims)

		// Extract userID and role from claims
		str, ok := claims["userID"].(float64)
		if !ok {
			log.Printf("failed to extract userID from token claims")
			permissionDenied(c.Writer)
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			log.Printf("failed to extract role from token claims")
			permissionDenied(c.Writer)
			c.Abort()
			return
		}

		// Convert userID to int
		userID := int(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			permissionDenied(c.Writer)
			c.Abort()
			return
		}

		// Add userID and role to the request context
		c.Set(string(UserKey), userID)
		c.Set(string(RoleKey), role)

		// Call the next handler
		c.Next()
	}
}

// AdminOnly middleware ensures that only users with the "admin" role can access the route.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve role from the request context
		role, exists := c.Get(string(RoleKey))
		if !exists {
			log.Printf("failed to retrieve role from context")
			permissionDenied(c.Writer)
			c.Abort()
			return
		}

		// Check if the user is an admin
		if role != "admin" {
			log.Printf("user does not have admin role")
			permissionDenied(c.Writer)
			c.Abort()
			return
		}

		// Call the next handler
		c.Next()
	}
}

// permissionDenied writes a "permission denied" error response.
func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}
