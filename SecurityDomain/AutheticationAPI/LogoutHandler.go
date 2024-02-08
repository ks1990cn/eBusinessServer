package main

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract token from request headers
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		http.Error(w, "Authorization header not provided", http.StatusBadRequest)
		return
	}
	region, err := getParameterValue("region")
	if err != nil {
		http.Error(w, "Unable to fetch region", http.StatusBadRequest)
		return
	}
	// Extract token from Authorization header (Bearer token)
	tokenParts := strings.Split(authorizationHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusBadRequest)
		return
	}
	accessToken := tokenParts[1]

	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize Cognito client
	svc := cognitoidentityprovider.New(sess)

	// Create input for logout
	logoutInput := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	}

	// Perform logout
	_, err = svc.GlobalSignOut(logoutInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}
