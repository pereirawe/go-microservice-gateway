package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pereirawe/go-microservice-gateway/config"
)

// LoginResponse represents the successful login response.
type LoginResponse struct {
	Token string `json:"token"`
}

// TODO: Load config from env vars
// var cfg *config.Config

func handleLoginSuccess(w http.ResponseWriter, username string) {
	log.Printf("Login successful for user: %s", username)

	token := generateToken(username)
	response := &LoginResponse{Token: token}

	CreateJSONResponse(w, response, http.StatusOK)
}

// handleLoginFailure handles an invalid login.
func handleLoginFailure(w http.ResponseWriter, errorMessage string, statusCode int) {
	log.Printf("Login failed: %s", errorMessage)
	response := map[string]string{"error": errorMessage}
	CreateJSONResponse(w, response, statusCode)
}

// LoginHandler handles user authentication and token generation.
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}
	if r.Method != "POST" {
		handleLoginFailure(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok || username != cfg.APPUser || password != cfg.APPPass {
		handleLoginFailure(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	handleLoginSuccess(w, username)
}

// JWTAuthorizationHandler validates JWT tokens for incoming requests.
func JWTAuthorizationHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		handleLoginFailure(w, "Invalid or missing token", http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[7:]
	if _, err := verifyToken(tokenString); err != nil {
		handleLoginFailure(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	CreateJSONResponse(w, map[string]string{"message": "Authorization successful"}, http.StatusOK)
}

// verifyToken verifies the JWT token.
func verifyToken(tokenString string) (jwt.MapClaims, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading configuration: %s", err)
		return nil, fmt.Errorf("could not load configuration")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.APPSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}

// generateToken generates a new JWT token.
func generateToken(username string) string {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.APPSecret))
	if err != nil {
		log.Fatalf("failed to sign token: %v", err)
	}
	return tokenString
}

// JWTMiddleware is a middleware that validates JWT tokens.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			handleLoginFailure(w, "Invalid or missing token", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[7:]
		if _, err := verifyToken(tokenString); err != nil {
			log.Printf("Token verification failed: %v", err)
			handleLoginFailure(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
