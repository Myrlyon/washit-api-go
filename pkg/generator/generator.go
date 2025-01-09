package generate

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	random "crypto/rand"

	"github.com/bwmarrin/snowflake"
	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/rand"
)

func RandomPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	password := make([]byte, 69)
	charsetLen := byte(len(charset))

	for i := range password {
		randomByte, err := random.Int(random.Reader, big.NewInt(int64(charsetLen)))
		if err != nil {
			return "", err
		}
		password[i] = charset[randomByte.Int64()]
	}

	return string(password), nil
}

func ImageFromUrl(imageUrl string) (imagePath string, err error) {
	timeID := time.Now().UnixNano()
	savePath := fmt.Sprintf("./public/profilePic/%d.jpg", timeID)

	err = os.MkdirAll("./public/profilePic", os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	client := resty.New()

	resp, err := client.R().SetOutput(savePath).Get(imageUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get image: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to download image, status: %s", resp.Status())
	}

	return fmt.Sprintf("%d.jpg", timeID), nil
}

func SnowflakeID(nodeID int64) (id int64, err error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		fmt.Println("Error creating Snowflake node:", err)
		return 0, err
	}
	return node.Generate().Int64(), nil
}

func AlphaNumericID(prefix string) (id string, err error) {
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

func SaveMediaToFile(imageData []byte, savePath string) error {
	err := os.MkdirAll("media/pelaporan", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = os.WriteFile(savePath, imageData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to save profile picture: %v", err)
	}

	return nil
}
