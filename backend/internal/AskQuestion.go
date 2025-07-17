package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AskQuestion(c *gin.Context) {
	var input struct {
		Question string `json:"question"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing question"})
		return
	}

	// TODO: Search embeddings & query OpenAI
	c.JSON(http.StatusOK, gin.H{
		"question": input.Question,
		"answer":   "This is a stubbed answer.",
	})
}
