package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Parse the request body into a LoginRequest struct
	var loginReq LoginRequest
	err = json.Unmarshal(body, &loginReq)
	if err != nil {
		http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
		return
	}
	region, err := getParameterValue("region")
	if err != nil {
		http.Error(w, "Unable to fetch region", http.StatusBadRequest)
		return
	}
	clientId, err := getParameterValue("clientIdUserPool")
	if err != nil {
		http.Error(w, "Unable to fetch Client ID", http.StatusBadRequest)
		return
	}
	// Extract username and password from the parsed struct
	username := loginReq.Username
	password := loginReq.Password

	// Initialize AWS session with your region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize Cognito client
	svc := cognitoidentityprovider.New(sess)

	// Create input for authentication
	authParams := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		ClientId: aws.String(clientId),
	}

	// Initiate user authentication
	authResp, err := svc.InitiateAuth(authParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if authentication was successful
	if authResp.AuthenticationResult != nil {
		// Authentication successful, return JWT token
		w.WriteHeader(http.StatusOK)
		response := struct {
			Message  string `json:"message"`
			JWTToken string `json:"jwt_token"`
			IdToken  string `json:"id_token"`
		}{
			Message:  "Authentication successful",
			JWTToken: *authResp.AuthenticationResult.AccessToken,
			IdToken:  *authResp.AuthenticationResult.IdToken,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		// Check if a new password is required
		if authResp.ChallengeName != nil && *authResp.ChallengeName == "NEW_PASSWORD_REQUIRED" {
			// Redirect to update password endpoint
			http.Error(w, "New Password required", http.StatusSeeOther)
			return
		} else {
			// Authentication failed, handle accordingly
			fmt.Fprintf(w, "Authentication failed\n")
		}
	}
}
