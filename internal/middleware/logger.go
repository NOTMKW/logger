package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/notmkw/log/internal/models"
)

func NewEarlyRequestLogData(includeSensitive bool) *models.EarlyRequestLogData {
	return &models.EarlyRequestLogData{
		IncludeSensitiveHeaders: includeSensitive,
		SensitiveHeaders: []string{
			"Authorization",
			"Cookie",
			"X-api-key",
			"X-auth-token",
			"authentication",
			"X-access-token",
			"bearer",
		},
	}
}

func LoggerMiddleware(erl *models.EarlyRequestLogData) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()
		c.Locals("request_id", requestID)

		start := time.Now()

		err := c.Next()

		latency := time.Since(start)

		logData := CaptureRequestData(c, erl, requestID, latency)

		LogToConsole(logData)

		return err
	}
}

func CaptureRequestData(c *fiber.Ctx, erl *models.EarlyRequestLogData, requestID string, latency time.Duration) models.RequestLogData {
	converted := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		converted[string(key)] = string(value)
	})

	sanitizedHeaders := SanitizeHeaders(converted, erl)

	queryParams := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	scheme := "http"
	if c.Secure() {
		scheme = "https"
	}

	return models.RequestLogData{
		RequestID:     requestID,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		HTTPMethod:    c.Method(),
		URL:           c.OriginalURL(),
		Headers:       sanitizedHeaders,
		QueryParams:   queryParams,
		IPAddress:     c.IP(),
		UserAgent:     c.Get("User-Agent"),
		Latency:       latency.String(),
		ContentType:   c.Get("Content-Type"),
		ContentLength: c.Get("Content-Length"),
		Host:          c.Hostname(),
		Referer:       c.Get("Referer"),
		Protocol:      c.Protocol(),
		Scheme:        scheme,
	}
}

func SanitizeHeaders(headers map[string]string, erl *models.EarlyRequestLogData) map[string]string {
	sanitized := make(map[string]string)

	for key, value := range headers {
		lowerKey := strings.ToLower(key)
		isSensitive := false

		if erl.IncludeSensitiveHeaders {
			for _, sensitiveHeader := range erl.SensitiveHeaders {
				if strings.Contains(lowerKey, strings.ToLower(sensitiveHeader)) {
					isSensitive = true
					break
				}
			}
		}

		if isSensitive {
			sanitized[key] = "REDACTED"
		} else {
			sanitized[key] = value
		}
	}

	return sanitized
}

func LogToConsole(logData models.RequestLogData) {
	separator := strings.Repeat("=", 80)
	fmt.Println(separator)
	fmt.Printf("Request ID: %s\n", logData.RequestID)
	fmt.Printf("Timestamp: %s\n", logData.Timestamp)
	fmt.Printf("HTTP Method: %s\n", logData.HTTPMethod)
	fmt.Printf("URL: %s\n", logData.URL)
	fmt.Printf("IP Address: %s\n", logData.IPAddress)
	fmt.Printf("User Agent: %s\n", logData.UserAgent)
	fmt.Printf("Host: %s\n", logData.Host)
	fmt.Printf("Referer: %s\n", logData.Referer)
	fmt.Printf("Protocol: %s\n", logData.Protocol)
	fmt.Printf("Scheme: %s\n", logData.Scheme)
	fmt.Printf("Latency: %s\n", logData.Latency)

	if logData.ContentType != "" {
		fmt.Printf("Content Type: %s\n", logData.ContentType)
	}

	if logData.ContentLength != "" {
		fmt.Printf("Content Length: %s\n", logData.ContentLength)
	}

	if len(logData.QueryParams) > 0 {
		fmt.Printf("Query Parameters:\n")
		for key, value := range logData.QueryParams {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	fmt.Printf("Headers:\n")
	for key, value := range logData.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}

	fmt.Println(separator)
	fmt.Println()
}

func LogToJSON(logData models.RequestLogData) {
	jsonData, err := json.MarshalIndent(logData, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling log data to JSON: %v\n", err)
		return
	}

	fmt.Printf("Request Data (JSON):\n%s\n", string(jsonData))
}
