package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found or unable to load")
	}
}

type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func OpenaiGenerate(app *fiber.App) {
	app.Post("/generate", generateAIResponse)
	app.Get("/", testfunc)
}

func callOpenAi(prompt string) (string, error) {
	apiKey := os.Getenv("OPEN_AI_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY is not set")
	}
	url := "https://api.openai.com/v1/chat/completions"

	payload := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: prompt},
		},
		MaxTokens: 150,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Make HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(body))

	// Parse OpenAI response
	var openAIResp OpenAIResponse
	err = json.Unmarshal(body, &openAIResp)
	if err != nil {
		return "", err
	}

	// Extract response content
	if openAIResp.Choices.Message.Content != "" {
		return openAIResp.Choices.Message.Content, nil
	}
	return "", fmt.Errorf("no response from OpenAI")

}

func generateAIResponse(c *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Prompt string `json:"prompt"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Call OpenAI
	response, err := callOpenAi(req.Prompt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return response
	return c.JSON(fiber.Map{"response": response})
}
