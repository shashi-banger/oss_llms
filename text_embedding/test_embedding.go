package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func getEmbeddings(texts []string) (*EmbeddingResponse, error) {
	url := "http://localhost:4000/v1/embeddings"

	reqBody := EmbeddingRequest{
		Model: "qwen-embedding",
		Input: texts,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-sb123")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var embeddingResp EmbeddingResponse
	err = json.Unmarshal(body, &embeddingResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &embeddingResp, nil
}

func main() {
	fmt.Println("Testing Qwen Embedding API...")

	texts := []string{
		"Hello, how are you?",
		"The weather is nice today.",
		"Machine learning is fascinating.",
		"Go is a great programming language.",
	}

	response, err := getEmbeddings(texts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Model used: %s\n", response.Model)
	fmt.Printf("Number of embeddings: %d\n", len(response.Data))
	if len(response.Data) > 0 {
		fmt.Printf("Embedding dimension: %d\n", len(response.Data[0].Embedding))
	}
	fmt.Printf("Total tokens used: %d\n", response.Usage.TotalTokens)

	// Print first few dimensions of each embedding
	for i, data := range response.Data {
		fmt.Printf("\nText %d: '%s'\n", i+1, texts[i])
		if len(data.Embedding) >= 5 {
			fmt.Printf("Embedding (first 5 dims): %.6f, %.6f, %.6f, %.6f, %.6f\n",
				data.Embedding[0], data.Embedding[1], data.Embedding[2],
				data.Embedding[3], data.Embedding[4])
		}
	}
}
