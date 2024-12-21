package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/rand"
)

func TakeGoogleImage(imageUrl string) (imagePath string, err error) {
	sId, err := SnowflakeId(1)
	if err != nil {
		return "", fmt.Errorf("failed to generate Snowflake ID: %w", err)
	}
	savePath := fmt.Sprintf("./public/profilePic/%d.jpg", sId)

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

	return fmt.Sprintf("%d.jpg", sId), nil
}

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
