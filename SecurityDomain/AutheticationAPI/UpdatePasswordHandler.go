package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"newPassword"`
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
