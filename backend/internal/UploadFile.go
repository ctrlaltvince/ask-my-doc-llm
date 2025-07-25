package internal

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

var S3Client *s3.Client

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("file missing: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file"})
		return
	}

	f, err := file.Open()
	if err != nil {
		log.Printf("file open error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		log.Printf("read error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// upload to S3
	_, err = S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("ask-my-doc-llm-files"),
		Key:    aws.String("uploads/" + file.Filename),
		Body:   bytes.NewReader(content),
		// Optional if bucket is already configured:
		// SSEKMSKeyId: aws.String("alias/ask-my-doc-s3"),
	})
	if err != nil {
		log.Printf("S3 upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to S3"})
		return
	}

	log.Printf("upload succeeded: %s", file.Filename)
	c.JSON(http.StatusOK, gin.H{"status": "Upload successful"})
}
