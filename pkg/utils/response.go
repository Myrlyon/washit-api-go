package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type HypermediaLink struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

type SuccessResponseFormat struct {
	Status     string                    `json:"status"`
	StatusCode int                       `json:"statusCode"`
	Message    string                    `json:"message"`
	Data       interface{}               `json:"data,omitempty"`
	Meta       MetaInfo                  `json:"meta"`
	Links      map[string]HypermediaLink `json:"_links,omitempty"`
}

type ErrorResponseFormat struct {
	Status     string                    `json:"status"`
	StatusCode int                       `json:"statusCode"`
	Message    string                    `json:"message"`
	Error      ErrorInfo                 `json:"error,omitempty"`
	Meta       MetaInfo                  `json:"meta"`
	Links      map[string]HypermediaLink `json:"_links,omitempty"`
}

type ErrorInfo struct {
	Type    string `json:"type"`
	Details string `json:"details"`
}

type MetaInfo struct {
	RequestID string    `json:"requestId,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, links map[string]HypermediaLink) {
	response := SuccessResponseFormat{
		Status:     "success",
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta: MetaInfo{
			RequestID: c.GetString("requestId"),
			Timestamp: time.Now().UTC(),
		},
		Links: links,
	}
	c.JSON(statusCode, response)
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := ErrorResponseFormat{
		Status:     "error",
		StatusCode: statusCode,
		Message:    message,
		Error: ErrorInfo{
			Type:    fmt.Sprintf("%T", err),
			Details: err.Error(),
		},
		Meta: MetaInfo{
			RequestID: c.GetString("requestId"),
			Timestamp: time.Now().UTC(),
		},
	}
	c.JSON(statusCode, response)
}

func ParseJson(c *gin.Context, v any) error {
	if c.Request.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(c.Request.Body).Decode(v)
}

func CopyTo(src interface{}, dest interface{}) {
	data, _ := json.Marshal(src)
	_ = json.Unmarshal(data, dest)
}

func ToData(title string, ConvertedData any) (responseData any) {
	responseData = map[string]interface{}{
		"message": title + " is collected successfully",
		title:     ConvertedData,
	}
	return
}
