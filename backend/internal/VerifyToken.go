package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyToken(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No claims found"})
		return
	}

	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
		return
	}

	// Optionally, filter or return the full claims
	c.JSON(http.StatusOK, gin.H{
		"status": "Token verified",
		"claims": claimsMap,
	})
}
