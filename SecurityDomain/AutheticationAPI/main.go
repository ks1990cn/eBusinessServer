package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/updatepassword", updatePasswordHandler)
	http.ListenAndServe(":8080", addCorsHeaders(http.DefaultServeMux)) // Add CORS middleware
}

// Middleware function to add CORS headers to every response
func addCorsHeaders(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
func getParameterValue(parameterName string) (string, error) {
	// Initialize a new AWS session using environment credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a SSM client
	ssmClient := ssm.New(sess)

	// Get the parameter value from Parameter Store
	output, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(parameterName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	return *output.Parameter.Value, nil
}
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
		// Authentication successful, return JWT token or any other response
		fmt.Fprintf(w, "Authentication successful!\n")
		fmt.Fprintf(w, "JWT Token: %s\n", *authResp.AuthenticationResult.IdToken)
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
func updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
		return
	}
	// Initialize a new AWS session using environment credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	userPoolId, err := getParameterValue("userPoolId")
	if err != nil {
		http.Error(w, "Unable to fetch Client ID", http.StatusBadRequest)
		return
	}
	// Create a Cognito Identity Provider client
	svc := cognitoidentityprovider.New(sess)
	requestBody.Password = "12345678"
	// Parameters for updating the user's password
	input := &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(userPoolId), // Replace with your actual Cognito user pool ID
		Username:   aws.String(requestBody.Username),
		Password:   aws.String(requestBody.Password),
		Permanent:  aws.Bool(true), // Set to true if the password change is permanent
	}

	// Update the user's password
	_, err = svc.AdminSetUserPassword(input)
	if err != nil {
		http.Error(w, "Failed to update password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password updated successfully"))
}
