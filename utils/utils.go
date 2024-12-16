package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var Validate = validator.New()

func WriteJson(c *gin.Context, statusCode int, data map[string]interface{}) {
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(statusCode)
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Writer.Write(jsonData)
}

func WriteError(c *gin.Context, status int, err error) {
	WriteJson(c, status, map[string]interface{}{"error": err.Error()})
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
