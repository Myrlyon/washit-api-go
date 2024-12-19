package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/rand"
)

var Validate = validator.New()

func MakeProfileImage(firstName string, lastName string) (imagePath string, err error) {
	sId, err := SnowflakeId(1)
	if err != nil {
		return "", fmt.Errorf("failed to generate Snowflake ID: %w", err)
	}

	imageURL := "https://avatar.iran.liara.run/username?username=" + firstName + "+" + lastName
	savePath := fmt.Sprintf("./public/profilePic/%d.jpg", sId)

	err = os.MkdirAll("./public/profilePic", os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	client := resty.New()

	resp, err := client.R().SetOutput(savePath).Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to get image: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to download image, status: %s", resp.Status())
	}

	return fmt.Sprintf("%d.jpg", sId), nil
}

func SnowflakeId(nodeID int64) (id int64, err error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		fmt.Println("Error creating Snowflake node:", err)
		return 0, err
	}
	return node.Generate().Int64(), nil
}

func AlphaNumericId(prefix string) (id string, err error) {
	length := 10
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var result []byte

	seed := uint64(time.Now().UnixNano())
	if seed == 0 {
		return "", errors.New("failed to generate seed for random number generator")
	}
	rand.Seed(seed)

	for i := 0; i < length; i++ {
		index := rand.Intn(len(charset))
		if index < 0 || index >= len(charset) {
			return "", errors.New("random index out of bounds")
		}
		result = append(result, charset[index])
	}

	return prefix + "-" + string(result), nil
}

func WriteJson(c *gin.Context, statusCode int, data interface{}) {
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

func ToData(title string, ConvertedData any) (responseData any) {
	responseData = map[string]interface{}{title: ConvertedData}
	return
}
