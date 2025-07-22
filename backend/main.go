package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/ctrlaltvince/ask-my-doc-llm/internal"
)

var (
	clientID     = "39u7iped9gp9cfnfutjp1ras8b"
	clientSecret = "22hgbmveqbd36du39hbg43hgs18nm9vtjmqlop13o165b9ea3kj"
	redirectURL  = "http://localhost:5173/oauth/callback"
	issuerURL    = "https://cognito-idp.us-west-1.amazonaws.com/us-west-1_RdclhXSHD"
	oauth2Config oauth2.Config
	provider     *oidc.Provider
)

func initOIDC() {
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

// Middleware to handle CORS in Gin
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

// JWT Middleware for Gin
func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]

		ctx := context.Background()
		idToken, err := provider.Verifier(&oidc.Config{ClientID: clientID}).Verify(ctx, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token claims"})
			return
		}

		// Save claims in Gin context for handlers
		c.Set("claims", claims)

		c.Next()
	}
}

func callbackHandler(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or missing code"})
		return
	}

	// ✅ Now we can log the actual received code
	log.Printf("Received code: %s", req.Code)
	log.Printf("Trying to exchange using redirect URI: %s", oauth2Config.RedirectURL)

	ctx := c.Request.Context()
	token, err := oauth2Config.Exchange(ctx, req.Code)
	if err != nil {
		log.Printf("Token exchange error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Printf("No id_token found in token response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token field in oauth2 token"})
		return
	}

	idToken, err := provider.Verifier(&oidc.Config{ClientID: clientID}).Verify(ctx, rawIDToken)
	if err != nil {
		log.Printf("ID Token verification failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID Token: " + err.Error()})
		return
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("Failed to parse claims: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims: " + err.Error()})
		return
	}

	email, _ := claims["email"].(string)

	c.JSON(http.StatusOK, gin.H{
		"email":        email,
		"id_token":     rawIDToken,
		"access_token": token.AccessToken,
	})
}

func profileHandler(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No claims found"})
		return
	}

	claimsMap := claims.(map[string]interface{})
	email, _ := claimsMap["email"].(string)

	c.JSON(http.StatusOK, gin.H{"email": email})
}

func main() {
	initOIDC()

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Backend is running!\n")
	})

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Public route for callback (login)
	r.POST("/callback", callbackHandler)

	// Protected routes
	auth := r.Group("/")
	auth.Use(jwtMiddleware())
	{
		auth.GET("/profile", profileHandler)
		auth.POST("/ask", internal.AskQuestion)
		auth.POST("/upload", internal.UploadFile)
		auth.POST("/verify", internal.VerifyToken)
	}

	log.Println("Server running on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
