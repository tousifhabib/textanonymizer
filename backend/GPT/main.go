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
    spacyURL      = "http://localhost:5000/anonymize"
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

    Client interface {
        Anonymize(ctx context.Context, text string) (string, error)
    }

    HTTPClient struct {
        apiKey string
        url    string
        client *http.Client
    }
)

func NewHTTPClient(apiKey, url string, timeout time.Duration) *HTTPClient {
    return &HTTPClient{
        apiKey: apiKey,
        url:    url,
        client: &http.Client{Timeout: timeout},
    }
}

func (c *HTTPClient) Anonymize(ctx context.Context, text string) (string, error) {
    var payload interface{}

    if c.url == spacyURL {
        payload = map[string]interface{}{
            "text": text,
        }
    } else {
        payload = map[string]interface{}{
            "model": "gpt-3.5-turbo-0125",
            "messages": []map[string]interface{}{
                {"role": "user", "content": fmt.Sprintf(`Please anonymize the text by redacting any names of people with REDACTED. Here is the text: %s`, text)},
            },
            "max_tokens": 500,
        }
    }

    return c.makeRequest(ctx, payload)
}

func (c *HTTPClient) makeRequest(ctx context.Context, payload interface{}) (string, error) {
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("failed to marshal JSON payload: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(jsonPayload))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", contentType)
    if c.apiKey != "" {
        req.Header.Set(authHeader, "Bearer "+c.apiKey)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API returned non-200 status: %s", resp.Status)
    }

    var result OpenAIResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("failed to decode response: %w", err)
    }

    if len(result.Choices) == 0 {
        return "", fmt.Errorf("no choices returned from API")
    }

    return postProcessAnonymizedText(result.Choices[0].Message.Content), nil
}

func postProcessAnonymizedText(text string) string {
    re := regexp.MustCompile(`<\|?endoftext\|?>`)
    return re.ReplaceAllString(text, "")
}

func main() {
    logger := initLogger()
    loadEnv(logger)

    apiKey := getAPIKey(logger)
    openAIClient := NewHTTPClient(apiKey, openAIURL, openAITimeout)
    spacyClient := NewHTTPClient("", spacyURL, 0)

    router := gin.Default()
    router.Use(cors.Default())
    router.POST("/anonymize-gpt", anonymizeHandler(openAIClient, logger, false))
    router.POST("/anonymize-spacy", anonymizeHandler(spacyClient, logger, true))

    logger.Info("Starting server on port 8080")
    if err := router.Run(":8080"); err != nil {
        logger.Fatalf("Failed to run server: %v", err)
    }
}

func initLogger() *logrus.Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.TextFormatter{
        FullTimestamp:   true,
        ForceColors:     true,
        DisableColors:   false,
        TimestampFormat: time.RFC3339,
    })
    logger.SetLevel(logrus.InfoLevel)
    return logger
}

func loadEnv(logger *logrus.Logger) {
    if err := godotenv.Load(); err != nil {
        logger.Fatalf("Error loading .env file: %v", err)
    }
}

func getAPIKey(logger *logrus.Logger) string {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        logger.Fatal("Missing OpenAI API key")
    }
    return apiKey
}

func anonymizeHandler(client Client, logger *logrus.Logger, isSpacy bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req Request
        if err := c.ShouldBindJSON(&req); err != nil {
            logger.Errorf("Failed to bind JSON: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        anonymizedText, err := client.Anonymize(c.Request.Context(), req.Text)
        if err != nil {
            logger.Errorf("Failed to anonymize text: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        logger.Infof("Successfully anonymized text using %s", func() string {
            if isSpacy {
                return "spaCy"
            }
            return "OpenAI"
        }())
        c.JSON(http.StatusOK, Response{AnonymizedText: anonymizedText})
    }
}
