package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TokenValidation(w http.ResponseWriter, r *http.Request) {
	// Extract token from request headers
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		http.Error(w, "Authorization header not provided", http.StatusBadRequest)
		return
	}
	// Extract token from Authorization header (Bearer token)
	tokenParts := strings.Split(authorizationHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusBadRequest)
		return
	}
	accessToken := tokenParts[1]

	// Parse token without validation
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		http.Error(w, "Failed to parse token: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the token is expired
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusBadRequest)
		return
	}

	exp, exists := claims["exp"].(float64)
	if !exists {
		http.Error(w, "exp claim not found", http.StatusBadRequest)
		return
	}

	expirationTime := time.Unix(int64(exp), 0)
	if expirationTime.Before(time.Now()) {
		http.Error(w, "Token expired", http.StatusUnauthorized)
		return
	}

	// If we reach this point, token is valid
	w.WriteHeader(http.StatusOK)
}
