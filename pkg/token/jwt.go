package jwt

import (
	"log"
	"strings"
	"time"

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

func GenerateAccessToken(payload map[string]interface{}) string {
	payload["type"] = AccessTokenType
	tokenContent := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Second * AccessTokenExpiredTime).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(configs.Envs.AuthSecret))
	if err != nil {
		log.Println("Failed to generate access token: ", err)
		return ""
	}

	return token
}

func GenerateRefreshToken(payload map[string]interface{}) string {
	payload["type"] = RefreshTokenType
	tokenContent := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Second * RefreshTokenExpiredTime).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(configs.Envs.AuthSecret))
	if err != nil {
		log.Println("Failed to generate refresh token: ", err)
		return ""
	}

	return token
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

// func ValidateGoogleToken(jwtToken string) (map[string]interface{}, error) {
// 	tokenData := jwt.MapClaims{}
// 	token, err := jwt.ParseWithClaims(jwtToken, tokenData, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(configs.Envs.GoogleAuthSecret), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if !token.Valid {
// 		return nil, jwt.ErrInvalidKey
// 	}

// 	var data map[string]interface{}
// 	utils.CopyTo(tokenData["payload"], &data)
// 	return data, nil
// }
