package internal

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"path/filepath"
	"time"
)

func UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file"})
		return
	}

	// Open uploaded file
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	// Read file into memory
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	fileBytes := buf.Bytes()

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AWS session"})
		return
	}

	// Generate a unique key for the S3 object
	key := fmt.Sprintf("uploads/%s_%d%s", uuid.New().String(), time.Now().Unix(), filepath.Ext(fileHeader.Filename))

	// Upload the file to S3
	s3Client := s3.New(sess)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("ask-my-doc-llm-files"),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileBytes),
		ACL:    aws.String("private"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to S3"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"status": "File uploaded successfully",
		"key":    key,
	})
}
