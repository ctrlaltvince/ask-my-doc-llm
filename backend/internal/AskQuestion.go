package internal

import (
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type scoredChunk struct {
	Text      string
	Score     float32
	Embedding []float32
}

func AskQuestion(c *gin.Context) {
	var input struct {
		Question string `json:"question"`
		Filename string `json:"filename"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing question or filename"})
		return
	}

	ctx := c.Request.Context()

	// 🧼 Normalize filename (strip .pdf/.txt/etc)
	baseFilename := strings.TrimSuffix(input.Filename, filepath.Ext(input.Filename))

	// 1. Get the full extracted text from S3
	log.Printf("Looking for: uploads/%s.txt", baseFilename)
	fullText, err := GetExtractedTextFromS3(ctx, baseFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load file: " + err.Error()})
		return
	}

	// 2. Chunk the text
	chunks := ChunkText(fullText, 200)

	// 3. Embed all chunks
	var scoredChunks []scoredChunk
	for _, chunk := range chunks {
		embedding, err := GetEmbedding(chunk)
		if err != nil {
			log.Printf("embedding failed: %v", err)
			continue
		}
		scoredChunks = append(scoredChunks, scoredChunk{
			Text:      chunk,
			Embedding: embedding,
		})
	}

	if len(scoredChunks) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to embed document content"})
		return
	}

	// 4. Embed the question
	qEmbedding, err := GetEmbedding(input.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to embed question"})
		return
	}

	// 5. Score chunks
	for i := range scoredChunks {
		scoredChunks[i].Score = CosineSimilarity(scoredChunks[i].Embedding, qEmbedding)
	}

	// 6. Select top 3
	sort.Slice(scoredChunks, func(i, j int) bool {
		return scoredChunks[i].Score > scoredChunks[j].Score
	})
	topChunks := scoredChunks
	if len(topChunks) > 3 {
		topChunks = topChunks[:3]
	}

	// 7. Build prompt
	var contextText strings.Builder
	for _, chunk := range topChunks {
		contextText.WriteString(chunk.Text + "\n")
	}
	prompt := "Use the following context to answer the user's question:\n\n" + contextText.String() + "\nQuestion: " + input.Question

	// 8. Ask GPT-4
	answer, err := AskOpenAI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"question": input.Question,
		"answer":   answer,
	})
}
