package models 

import (
	"time"
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

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}