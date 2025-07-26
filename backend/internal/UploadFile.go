package internal

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var acceptedExtensions = map[string]bool{
	".pdf": true,
	".txt": true,
	".md":  true,
	".csv": true,
}

var S3Client *s3.Client

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("file missing: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !acceptedExtensions[ext] {
		log.Printf("Rejected file with unsupported extension: %s", ext)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
		return
	}

	f, err := file.Open()
	if err != nil {
		log.Printf("file open error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	// Extract text based on file extension
	text, err := ExtractTextFromFile(f, ext)
	if err != nil {
		log.Printf("text extraction failed: %v", err)
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Text extraction failed: " + err.Error()})
		return
	}

	// Upload extracted text as a .txt file to S3
	key := "uploads/" + strings.TrimSuffix(file.Filename, ext) + ".txt"
	_, err = S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("ask-my-doc-llm-files"),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(text)),
	})
	if err != nil {
		log.Printf("S3 upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to S3"})
		return
	}

	log.Printf("upload succeeded: %s", key)
	c.JSON(http.StatusOK, gin.H{"status": "Upload successful"})
}
