package userService

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	userModel "washit-api/internal/user/dto/model"
	userRequest "washit-api/internal/user/dto/request"
	userRepository "washit-api/internal/user/repository"
	auths "washit-api/pkg/auth"
	generate "washit-api/pkg/generator"
	jwt "washit-api/pkg/token"

	"firebase.google.com/go/auth"
	"github.com/fatih/camelcase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type IUserService interface {
	RefreshToken(c context.Context, userID int64) (string, error)
	Register(c context.Context, req *userRequest.Register) (*userModel.User, error)
	LoginWithGoogle(c context.Context, req *userRequest.Google, userInfo *auth.UserInfo) (*userModel.User, string, string, error)
	Login(c context.Context, req *userRequest.Login) (*userModel.User, string, string, error)
	Logout(c context.Context, userID int64) error
	BanUser(c context.Context, userID int64) (*userModel.User, error)
	UnbanUser(c context.Context, userID int64) (*userModel.User, error)
	GetMe(c context.Context, userID int64) (*userModel.User, error)
	GetUserByID(c context.Context, userID int64) (*userModel.User, error)
	GetUsers(c context.Context) ([]*userModel.User, error)
	GetBannedUsers(c context.Context) ([]*userModel.User, error)
	UpdateProfile(c context.Context, userID int64, req *userRequest.UpdateProfile) (*userModel.User, error)
	UpdatePassword(c context.Context, userID int64, req *userRequest.UpdatePassword) error
	UpdatePicture(c context.Context, userID int64, req *userRequest.UpdatePicture) (*userModel.User, error)
}

type UserService struct {
	repository userRepository.IUserRepository
	validator  *validator.Validate
}

func NewUserService(
	repository userRepository.IUserRepository, validator *validator.Validate) *UserService {
	return &UserService{
		repository: repository,
		validator:  validator,
	}
}

func (s *UserService) RefreshToken(c context.Context, userID int64) (string, error) {
	user, err := s.repository.GetUserByID(c, userID)
	if err != nil {
		log.Printf("Failed to get user by id: %v", err)
		return "", fmt.Errorf("user not found: %v", userID)
	}

	tokenData := gin.H{
		"id":   strconv.FormatInt(user.ID, 10),
		"role": user.Role,
	}

	accessToken, err := jwt.GenerateAccessToken(tokenData)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}

	return accessToken, nil
}

func (s *UserService) LoginWithGoogle(c context.Context, req *userRequest.Google, userInfo *auth.UserInfo) (*userModel.User, string, string, error) {
	user, err := s.repository.GetUserByEmail(c, userInfo.Email)
	if err != nil {
		randomPassword, err := generate.RandomPassword()
		if err != nil {
			log.Printf("Error generating random password: %v", err)
			return nil, "", "", err
		}

		hashedPassword, err := auths.HashPassword(randomPassword)
		if err != nil {
			log.Printf("Error encrypting password: %v", err)
			return nil, "", "", err
		}

		snoflakeID, err := generate.SnowflakeID(1)
		if err != nil {
			log.Printf("Error generating snowflake ID: %v", err)
			return nil, "", "", err
		}

		imagePath, err := generate.ImageFromUrl(userInfo.PhotoURL)
		if err != nil {
			log.Printf("Error downloading Google profile image: %v", err)
			return nil, "", "", err
		}

		splittedName := camelcase.Split(userInfo.DisplayName)
		firstName, lastName := "", ""
		if len(splittedName) > 0 {
			firstName = splittedName[0]
		}
		if len(splittedName) > 1 {
			lastName = splittedName[1]
		}

		newUser := &userModel.User{
			ID:        snoflakeID,
			FirstName: firstName,
			LastName:  lastName,
			Email:     userInfo.Email,
			Password:  hashedPassword,
			Image:     imagePath,
		}

		if err := s.repository.CreateUser(c, newUser); err != nil {
			log.Printf("Error creating new user: %v", err)
			return nil, "", "", fmt.Errorf("failed to create user: %v", err)
		}

		user, err = s.repository.GetUserByEmail(c, newUser.Email)
		if err != nil {
			log.Printf("Error fetching newly created user: %v", err)
			return nil, "", "", fmt.Errorf("failed to fetch user after creation: %v", err)
		}
	}

	if req.FcmToken != "" {
		user.FcmToken = req.FcmToken
		if err := s.repository.UpdateUser(c, user); err != nil {
			log.Printf("Failed to update FCM token: %v", err)
			return nil, "", "", fmt.Errorf("failed to update FCM token: %v", err)
		}
	}

	if user.IsBanned {
		return nil, "", "", fmt.Errorf("user is banned")
	}

	tokenData := gin.H{
		"id":        strconv.FormatInt(user.ID, 10),
		"role":      user.Role,
		"fcm_token": req.FcmToken,
	}

	accessToken, err := jwt.GenerateAccessToken(tokenData)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		return nil, "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(tokenData)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Login(c context.Context, req *userRequest.Login) (*userModel.User, string, string, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Printf("Validation error for login request: %v", err)
		return nil, "", "", fmt.Errorf("validation error: %v", err)
	}

	user, err := s.repository.GetUserByEmail(c, req.Email)
	if err != nil {
		log.Printf("User not found with email: %s, error: %v", req.Email, err)
		return nil, "", "", fmt.Errorf("unable to find user with email: %s", req.Email)
	}

	if !auths.ComparePasswords(user.Password, []byte(req.Password)) {
		log.Printf("Invalid password for user: %s", req.Email)
		return nil, "", "", fmt.Errorf("invalid password")
	}

	if user.IsBanned {
		log.Printf("User is banned: %s", req.Email)
		return nil, "", "", fmt.Errorf("user is banned")
	}

	if req.FcmToken != "" {
		user.FcmToken = req.FcmToken
		if err := s.repository.UpdateUser(c, user); err != nil {
			log.Printf("Failed to update FCM token for user: %s, error: %v", req.Email, err)
			return nil, "", "", fmt.Errorf("failed to update FCM token: %v", err)
		}
	}

	tokenData := gin.H{
		"id":        strconv.FormatInt(user.ID, 10),
		"role":      user.Role,
		"fcm_token": req.FcmToken,
	}

	accessToken, err := jwt.GenerateAccessToken(tokenData)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		return nil, "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(tokenData)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Register(c context.Context, req *userRequest.Register) (*userModel.User, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Printf("Validation error for register request: %v", err)
		return nil, fmt.Errorf("validation error: %v", err)
	}

	if _, err := s.repository.GetUserByEmail(c, req.Email); err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	snoflakeID, err := generate.SnowflakeID(1)
	if err != nil {
		log.Printf("Error generating snowflake ID: %v", err)
		return nil, err
	}

	hashedPassword, err := auths.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, err
	}

	imagePath, err := generate.ImageFromUrl(
		fmt.Sprintf("https://avatar.iran.liara.run/username?username=%s+%s", req.FirstName, req.LastName))
	if err != nil {
		log.Printf("Error downloading profile image: %v", err)
		return nil, err
	}

	user := &userModel.User{
		ID:        snoflakeID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
		Image:     imagePath,
		Role:      "customer",
	}

	if err := s.repository.CreateUser(c, user); err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (s *UserService) Logout(c context.Context, userID int64) error {
	return nil
}

func (s *UserService) BanUser(c context.Context, userID int64) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(c, userID)
	if err != nil {
		log.Printf("Failed to get user by id: %v", err)
		return nil, fmt.Errorf("user not found: %d", userID)
	}

	if user.IsBanned {
		return nil, fmt.Errorf("user is already banned")
	}

	user.IsBanned = true

	if err := s.repository.UpdateUser(c, user); err != nil {
		log.Printf("Failed to ban user: %v", err)
		return nil, fmt.Errorf("failed to ban user: %v", err)
	}

	return user, nil
}

func (s *UserService) UnbanUser(c context.Context, userID int64) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(c, userID)
	if err != nil {
		log.Printf("Failed to get user by id: %v", err)
		return nil, fmt.Errorf("user not found: %v", userID)
	}

	if !user.IsBanned {
		return nil, fmt.Errorf("user is not banned")
	}

	user.IsBanned = false

	if err := s.repository.UpdateUser(c, user); err != nil {
		log.Printf("Failed to unban user: %v", err)
		return nil, fmt.Errorf("failed to unban user: %v", err)
	}

	return user, nil
}

func (s *UserService) UpdateProfile(c context.Context, userID int64, req *userRequest.UpdateProfile) (*userModel.User, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Printf("Validation error for update profile request: %v", err)
		return nil, fmt.Errorf("validation error: %v", err)
	}

	user, err := s.repository.GetUserByID(c, userID)
	if err != nil {
		log.Printf("User not found with id: %d, error: %v", userID, err)
		return nil, fmt.Errorf("user not found: %v", userID)
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.repository.UpdateUser(c, user); err != nil {
		log.Printf("Failed to update user: %d, error: %v", userID, err)
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}

func (s *UserService) UpdatePassword(c context.Context, userID int64, req *userRequest.UpdatePassword) error {
	if err := s.validator.Struct(req); err != nil {
		log.Printf("Validation error for update password request: %v", err)
		return fmt.Errorf("validation error: %v", err)
	}

	user, err := s.repository.GetUserByID(c, userID)
	if err != nil {
		log.Printf("Failed to get user by id: %v", err)
		return fmt.Errorf("user not found: %v", userID)
	}

	if !auths.ComparePasswords(user.Password, []byte(req.OldPassword)) {
		log.Printf("Invalid old password for user: %v", userID)
		return fmt.Errorf("invalid old password")
	}

	hashedPassword, err := auths.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("Failed to hash new password: %v", err)
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	user.Password = hashedPassword

	if err := s.repository.UpdateUser(c, user); err != nil {
		log.Printf("Failed to update user password: %v", err)
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}

func (s *UserService) UpdatePicture(c context.Context, userID int64, req *userRequest.UpdatePicture) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(c, userID)
	if err != nil {
		log.Printf("Failed to get user by ID: %v", err)
		return nil, fmt.Errorf("user not found: %v", userID)
	}

	mediaName := fmt.Sprintf("%d.jpg", time.Now().Unix())
	imagePath := fmt.Sprintf("./public/profilePic/%s", mediaName)
	if err := generate.SaveMediaToFile(req.Image, imagePath); err != nil {
		log.Printf("Failed to save image to directory: %v", err)
		return nil, fmt.Errorf("failed to save image to directory: %v", err)
	}

	user.Image = mediaName

	if err := s.repository.UpdateUser(c, user); err != nil {
		log.Printf("Failed to update user profile picture: %v", err)
		return nil, fmt.Errorf("failed to update profile picture: %v", err)
	}

	return user, nil
}

func (s *UserService) GetMe(c context.Context, userID int64) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(c, userID)
	if err != nil {
		log.Printf("Failed to get profile information for user ID %d: %v", userID, err)
		return nil, fmt.Errorf("user not found: %v", userID)
	}

	return user, nil
}

func (s *UserService) GetUsers(c context.Context) ([]*userModel.User, error) {
	users, err := s.repository.GetUsers(c)
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	return users, nil
}

func (s *UserService) GetBannedUsers(c context.Context) ([]*userModel.User, error) {
	users, err := s.repository.GetBannedUsers(c)
	if err != nil {
		log.Printf("Failed to get banned users: %v", err)
		return nil, fmt.Errorf("failed to get banned users: %v", err)
	}

	return users, nil
}

func (s *UserService) GetUserByID(c context.Context, userID int64) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(c, userID)
	if err != nil {
		log.Printf("Failed to get user by ID %d: %v", userID, err)
		return nil, fmt.Errorf("user not found: %v", userID)
	}

	return user, nil
}
