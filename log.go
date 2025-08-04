package log

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RequestLogData struct {
	RequestID     string            `json:"request_id"`
	Timestamp     string            `json:"timestamp"`
	HTTPMethod    string            `json:"http_method"`
	URL           string            `json:"url"`
	Headers       map[string]string `json:"headers"`
	QueryParams   map[string]string `json:"query_params"`
	IPAddress     string            `json:"ip_address"`
	UserAgent     string            `json:"user_agent"`
	Latency       string            `json:"latency"`
	ContentType   string            `json:"content_type"`
	ContentLength string            `json:"content_length"`
	Host          string            `json:"host"`
	Referer       string            `json:"referer"`
	Protocol      string            `json:"protocol"`
	Scheme        string            `json:"scheme"`
}

type EarlyRequestLogData struct {
	includeSensitiveHeaders bool
	sensitiveHeaders        []string
}

func NewEarlyRequestLogData(includeSensitive bool) *EarlyRequestLogData {
	return &EarlyRequestLogData{
		includeSensitiveHeaders: includeSensitive,
		sensitiveHeaders: []string{
			"Authorization",
			"Cookie",
			"x-api-key",
			"x-auth-token",
			"authentication",
			"x-access-token",
			"bearer",
		},
	}
}

func (erl *EarlyRequestLogData) LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()
		c.Locals("request_id", requestID)
		logData := erl.CaptureRequestData(c, requestID)

		erl.logToConsole(logData)
		return c.Next()
	}
}

func (erl *EarlyRequestLogData) CaptureRequestData(c *fiber.Ctx, requestID string) RequestLogData {
	rawHeaders := c.GetReqHeaders()
	converted := make(map[string]string)
	for k, v := range rawHeaders {
		converted[k] = v[0]
	}

	erl.sanitizeHeaders(converted)
	queryParams := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	scheme := "http"
	if c.Secure() {
		scheme = "https"
	}
	return RequestLogData{
		RequestID:     requestID,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		HTTPMethod:    c.Method(),
		URL:           c.OriginalURL(),
		Headers:       erl.sanitizeHeaders(converted),
		QueryParams:   queryParams,
		IPAddress:     c.IP(),
		UserAgent:     c.Get("User-Agent"),
		ContentType:   c.Get("Content-Type"),
		ContentLength: c.Get("Content-Length"),
		Host:          c.Hostname(),
		Referer:       c.Get("Referer"),
		Protocol:      c.Protocol(),
		Scheme:        scheme,
	}
}

func (erl *EarlyRequestLogData) sanitizeHeaders(headers map[string]string) map[string]string {
	sanitized := make(map[string]string)
	for key, values := range headers {
		lowerKey := strings.ToLower(key)

		headerValue := values

		isSensitive := false
		if !erl.includeSensitiveHeaders {
			for _, sensitiveHeader := range erl.sensitiveHeaders {
				if strings.Contains(lowerKey, sensitiveHeader) {
					isSensitive = true
					break
				}
			}
		}

		if isSensitive {
			sanitized[key] = "REDACTED"
		} else {
			sanitized[key] = headerValue
		}
	}
	return sanitized
}

func (erl *EarlyRequestLogData) logToConsole(logData RequestLogData) {
	seperator := strings.Repeat("=", 80)
	fmt.Println(seperator)
	fmt.Printf("Request ID: %s\n", logData.RequestID)
	fmt.Printf("Timestamp: %s\n", logData.Timestamp)
	fmt.Printf("HTTP Method: %s\n", logData.HTTPMethod)
	fmt.Printf("URL: %s\n", logData.URL)
	fmt.Printf("ip_address: %s\n", logData.IPAddress)
	fmt.Printf("User Agent: %s\n", logData.UserAgent)
	fmt.Printf("Host: %s\n", logData.Host)
	fmt.Printf("Referer: %s\n", logData.Referer)
	fmt.Printf("Protocol: %s\n", logData.Protocol)
	fmt.Printf("Scheme: %s\n", logData.Scheme)

	if logData.ContentType != "" {
		fmt.Printf("Content Type: %s\n", logData.ContentType)
	}

	if logData.ContentLength != "" {
		fmt.Printf("Content Length: %d\n", logData.ContentLength)
	}

	if logData.Referer != "" {
		fmt.Printf("Referer: %s\n", logData.Referer)
	}

	if len(logData.QueryParams) > 0 {
		fmt.Printf("Query Parameters: %v\n", logData.QueryParams)
		for key, value := range logData.QueryParams {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	fmt.Printf("Headers:\n")
	for key, value := range logData.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}

	fmt.Println("seperator")
	fmt.Println()
}

func (erl *EarlyRequestLogData) LogToJSON(logData RequestLogData) {
	jsonData, err := json.MarshalIndent(logData, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling log data to JSON: %v\n", err)
		return
	}

	fmt.Printf("Request Data (JSON):\n%s\n", string(jsonData))
}

func main() {
	app := fiber.New()
	logger := NewEarlyRequestLogData(false)
	app.Use(logger.LoggerMiddleware())

	app.Use(func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":    "Alpha API is running",
			"request_id": c.Locals("request_id"),
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
		})
	})

	app.Get("/user/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		return c.JSON(fiber.Map{
			"user_id":    userID,
			"request_id": c.Locals("request_id"),
			"message":    fmt.Sprintf("User ID: %s", userID),
		})
	})

	app.Post("/user", func(c *fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{
			"message":    "User created successfully",
			"request_id": c.Locals("request_id"),
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
		})
	})

	app.Put("/user/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		return c.JSON(fiber.Map{
			"user_id":    userID,
			"request_id": c.Locals("request_id"),
			"message":    fmt.Sprintf("User ID: %s updated successfully", userID),
		})
	})

	app.Delete("/user/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		return c.JSON(fiber.Map{
			"message":    fmt.Sprintf("User ID: %s deleted successfully", userID),
			"request_id": c.Locals("request_id"),
		})
	})

	fmt.Sprintf("Starting %s on port 3000...\n")

	log.Fatal(app.Listen(":3000"))
}
