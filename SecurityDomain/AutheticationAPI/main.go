package main

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/updatepassword", updatePasswordHandler)
	http.HandleFunc("/TokenValidation", TokenValidation)
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
