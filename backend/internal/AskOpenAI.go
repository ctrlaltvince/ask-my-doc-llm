package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func AskOpenAI(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("OPENAI_API_KEY not set")
	}

	payload := ChatCompletionRequest{
		Model: "gpt-4", // or "gpt-3.5-turbo"
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant that answers questions using only the provided context."},
			{Role: "user", Content: prompt},
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// Log the error but do not return it, as this is a deferred function
			// and we don't want to interrupt the main flow of the function.
			// log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", errors.New("OpenAI API error: " + string(b))
	}

	var result ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return result.Choices[0].Message.Content, nil
}
