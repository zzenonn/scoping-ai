package db

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	tnamessage "github.com/zzenonn/trainocate-tna/internal/message"
)

func init() {

	// Set log level based on environment variables
	switch logLevel := strings.ToLower(os.Getenv("LOG_LEVEL")); logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}

}

type OpenAiRepository struct {
	ApiKey      string
	OpenAiUrl   string
	Model       string
	Temperature float32
}

type OpenAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewOpenAiRepository(apiKey string, openAiUrl string, model string, temperature float32) OpenAiRepository {
	return OpenAiRepository{
		ApiKey:      apiKey,
		OpenAiUrl:   openAiUrl,
		Model:       model,
		Temperature: temperature,
	}
}

func (repo *OpenAiRepository) createRequestPayload(openAiContext string, input string) []byte {
	data := map[string]interface{}{
		"model":       repo.Model,
		"temperature": repo.Temperature,
		"messages": []OpenAiMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: input,
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to encode data to JSON: %v", err)
	}

	return jsonData
}

func (repo *OpenAiRepository) PostPrompt(ctx context.Context, aiContext string, prompt string) (tnamessage.ChatCompletion, error) {

	log.Debug("Posting prompt . . .")

	payload := repo.createRequestPayload(aiContext, prompt)

	req, err := http.NewRequest("POST", repo.OpenAiUrl, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+repo.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to make the request: %v", err)
		return tnamessage.ChatCompletion{}, err
	}
	defer resp.Body.Close()

	var chatCompletion tnamessage.ChatCompletion

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to read response body: %v", err)
		return tnamessage.ChatCompletion{}, err
	}

	// Logging the response body
	bodyString := string(bodyBytes)
	log.Debug(bodyString)

	// Decoding the body bytes into chatCompletion
	err = json.Unmarshal(bodyBytes, &chatCompletion)
	if err != nil {
		log.Errorf("Failed to unmarshal response body: %v", err)
		return tnamessage.ChatCompletion{}, err
	}

	return chatCompletion, nil
}

// The sendRequest function remains the same as you provided
