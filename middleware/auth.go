package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)
		if tokenString == "" {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("missing token"))
			return
		}

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Envs.JWT_SECRET), nil
		})
		if err != nil || !token.Valid {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}

		// Extract claims from the token
		claims := token.Claims.(jwt.MapClaims)

		// Extract userID and role from claims
		str, ok := claims["userID"].(string)
		if !ok {
			log.Printf("failed to extract userID from token claims")
			permissionDenied(w)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			log.Printf("failed to extract role from token claims")
			permissionDenied(w)
			return
		}

		// Convert userID to int
		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			permissionDenied(w)
			return
		}

		// Add userID and role to the request context
		ctx := context.WithValue(r.Context(), UserKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// AdminOnly middleware ensures that only users with the "admin" role can access the route.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve role from the request context
		role, ok := r.Context().Value(RoleKey).(string)
		if !ok {
			log.Printf("failed to retrieve role from context")
			permissionDenied(w)
			return
		}

		// Check if the user is an admin
		if role != "admin" {
			log.Printf("user does not have admin role")
			permissionDenied(w)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// permissionDenied writes a "permission denied" error response.
func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}
