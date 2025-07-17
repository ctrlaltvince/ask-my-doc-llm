package main

import (
	"github.com/ctrlaltvince/ask-my-doc-llm/internal"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/auth/verify", internal.VerifyToken)
	r.POST("/upload", internal.UploadFile)
	r.POST("/ask", internal.AskQuestion)

	r.Run(":8081")
}
