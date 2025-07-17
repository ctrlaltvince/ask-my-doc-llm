package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyToken(c *gin.Context) {
	// TODO: Parse and validate JWT from AWS Cognito
	c.JSON(http.StatusOK, gin.H{"status": "Token verified (stub)"})
}
