package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

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
