package jwt

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"washit-api/pkg/configs"
	"washit-api/pkg/utils"
)

const (
	AccessTokenExpiredTime  = 24 * 60 * 60   // 1 day
	RefreshTokenExpiredTime = 30 * 24 * 3600 // 30 days
	AccessTokenType         = "x-access"     // 5 minutes
	RefreshTokenType        = "x-refresh"    // 30 days
)

func GenerateAccessToken(payload map[string]interface{}) (string, error) {
	payload["type"] = AccessTokenType
	tokenContent := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Second * AccessTokenExpiredTime).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(configs.Envs.AuthSecret))
	if err != nil {
		log.Println("Failed to generate access token: ", err)
		return "", err
	}

	return token, nil
}

func GenerateRefreshToken(payload map[string]interface{}) (string, error) {
	payload["type"] = RefreshTokenType
	tokenContent := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Second * RefreshTokenExpiredTime).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(configs.Envs.AuthSecret))
	if err != nil {
		log.Println("Failed to generate refresh token: ", err)
		return "", err
	}

	return token, nil
}

func ValidateToken(jwtToken string) (map[string]interface{}, error) {
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Envs.AuthSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	var data map[string]interface{}
	utils.CopyTo(tokenData["payload"], &data)
	return data, nil
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
