package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientID     = "39u7iped9gp9cfnfutjp1ras8b"
	clientSecret = "22hgbmveqbd36du39hbg43hgs18nm9vtjmqlop13o165b9ea3kj"
	redirectURL  = "http://localhost:5173"
	issuerURL    = "https://cognito-idp.us-west-1.amazonaws.com/us-west-1_RdclhXSHD"
	oauth2Config oauth2.Config
	provider     *oidc.Provider
)

type CodeRequest struct {
	Code string `json:"code"`
}

func init() {
	var err error
	provider, err = oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
	}

	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CodeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Code == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	token, err := oauth2Config.Exchange(ctx, req.Code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token", http.StatusInternalServerError)
		return
	}

	idToken, err := provider.Verifier(&oidc.Config{ClientID: clientID}).Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	email, _ := claims["email"].(string)

	resp := map[string]string{
		"email":    email,
		"id_token": rawIDToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Also add Authorization here
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Verify the token
		ctx := r.Context()
		idToken, err := provider.Verifier(&oidc.Config{ClientID: clientID}).Verify(ctx, tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Store claims in context for handlers to use
		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(map[string]interface{})
	if !ok {
		http.Error(w, "No claims found", http.StatusInternalServerError)
		return
	}

	email, _ := claims["email"].(string)

	resp := map[string]string{
		"email": email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", callbackHandler)
	mux.Handle("/profile", jwtMiddleware(http.HandlerFunc(profileHandler)))

	handler := enableCORS(mux)
	log.Println("Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
