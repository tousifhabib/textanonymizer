package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "regexp"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
)

const (
    openAIURL     = "https://api.openai.com/v1/chat/completions"
    contentType   = "application/json"
    authHeader    = "Authorization"
    openAITimeout = 10 * time.Second
)

type (
    Request struct {
        Text string `json:"text"`
    }

    Response struct {
        AnonymizedText string `json:"anonymizedText"`
    }

    OpenAIResponse struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }

    OpenAIClient interface {
        Anonymize(ctx context.Context, text string) (string, error)
    }

    openAIClient struct {
        apiKey string
        client *http.Client
    }
)

func NewOpenAIClient(apiKey string) OpenAIClient {
    return &openAIClient{
        apiKey: apiKey,
        client: &http.Client{Timeout: openAITimeout},
    }
}

func (c *openAIClient) Anonymize(ctx context.Context, text string) (string, error) {
    prompt := fmt.Sprintf(`Please anonymize the text by redacting any names of people with REDACTED. If there is a first and last name it should be replaced with one REDACTED. You should comprehensively search the text for names and redact them so that in the end result there should be not a single instance of any name. Here is the text: %s`, text)

    payload := map[string]interface{}{
        "model": "gpt-3.5-turbo-0125",
        "messages": []map[string]interface{}{
            {"role": "user", "content": prompt},
        },
        "max_tokens": 500,
    }

    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("failed to marshal JSON payload: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", openAIURL, bytes.NewBuffer(jsonPayload))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

    resp, err := c.client.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("OpenAI API returned non-200 status: %s", resp.Status)
    }

    var result OpenAIResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("failed to decode response: %w", err)
    }

    if len(result.Choices) == 0 {
        return "", fmt.Errorf("no choices returned from OpenAI API")
    }

    return postProcessAnonymizedText(result.Choices[0].Message.Content), nil
}

func postProcessAnonymizedText(text string) string {
    re := regexp.MustCompile(`<\|?endoftext\|?>`)
    return re.ReplaceAllString(text, "")
}

func main() {
    logger := logrus.New()
    logger.SetFormatter(&logrus.TextFormatter{
        FullTimestamp:   true,
        ForceColors:     true,
        DisableColors:   false,
        TimestampFormat: time.RFC3339,
    })
    logger.SetLevel(logrus.InfoLevel)

    if err := godotenv.Load(); err != nil {
        logger.Fatalf("Error loading .env file: %v", err)
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        logger.Fatal("Missing OpenAI API key")
    }

    openAIClient := NewOpenAIClient(apiKey)

    router := gin.Default()
    router.Use(cors.Default())
    router.POST("/anonymize", AnonymizeTextHandler(openAIClient, logger))

    logger.Info("Starting server on port 8080")
    if err := router.Run(":8080"); err != nil {
        logger.Fatalf("Failed to run server: %v", err)
    }
}

func AnonymizeTextHandler(openAIClient OpenAIClient, logger *logrus.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req Request
        if err := c.ShouldBindJSON(&req); err != nil {
            logger.Errorf("Failed to bind JSON: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        anonymizedText, err := openAIClient.Anonymize(c.Request.Context(), req.Text)
        if err != nil {
            logger.Errorf("Failed to call OpenAI API: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        logger.Info("Successfully anonymized text")
        c.JSON(http.StatusOK, Response{AnonymizedText: anonymizedText})
    }
}
