package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetExtractedTextFromS3(ctx context.Context, filename string) (string, error) {
	key := fmt.Sprintf("uploads/%s.txt", filename)

	resp, err := S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("ask-my-doc-llm-files"),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("failed to get object from S3: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		log.Printf("failed to read object body: %v", err)
		return "", err
	}

	return buf.String(), nil
}
