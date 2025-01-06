package userService

import (
	"context"
	"fmt"
	"log"
	"strconv"

	userModel "washit-api/internal/user/dto/model"
	userRequest "washit-api/internal/user/dto/request"
	userRepository "washit-api/internal/user/repository"
	auths "washit-api/pkg/auth"
	generate "washit-api/pkg/generator"
	jwt "washit-api/pkg/token"
	"washit-api/pkg/utils"

	"firebase.google.com/go/auth"
	"github.com/fatih/camelcase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type IUserService interface {
	RefreshToken(c context.Context, userId string) (string, error)
	Register(c context.Context, req *userRequest.Register) (*userModel.User, error)
	LoginWithGoogle(c context.Context, req *userRequest.Google, userInfo *auth.UserInfo) (*userModel.User, any, any, error)
	Login(c context.Context, req *userRequest.Login) (*userModel.User, any, any, error)
	Logout(c context.Context, userId string) error
	BanUser(c context.Context, userId string) (*userModel.User, error)
	UnbanUser(c context.Context, userId string) (*userModel.User, error)
	GetMe(c context.Context, userId string) (*userModel.User, error)
	GetUserByID(c context.Context, userId string) (*userModel.User, error)
	GetUsers(c context.Context) ([]*userModel.User, error)
	GetBannedUsers(c context.Context) ([]*userModel.User, error)
	UpdateProfile(c context.Context, userId string, req *userRequest.UpdateProfile) (*userModel.User, error)
	UpdatePassword(c context.Context, userId string, req *userRequest.UpdatePassword) error
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

func (s *UserService) RefreshToken(c context.Context, userId string) (string, error) {
	user, err := s.repository.GetUserByID(c, userId)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return "", fmt.Errorf("user not found: %v", userId)
	}

	tokenData := gin.H{"id": strconv.FormatInt(user.ID, 10), "role": user.Role}

	accessToken := jwt.GenerateAccessToken(tokenData)
	return accessToken, nil
}

func (s *UserService) LoginWithGoogle(c context.Context, req *userRequest.Google, userInfo *auth.UserInfo) (*userModel.User, any, any, error) {
	user, err := s.repository.GetUserByEmail(c, userInfo.Email)
	if err != nil {
		randomPassword, err := generate.RandomPassword()
		if err != nil {
			log.Println("Error generating random password: ", err)
			return nil, nil, nil, err
		}

		hashedPassword, err := auths.HashPassword(randomPassword)
		if err != nil {
			log.Println("Error encrypting password: ", err)
			return nil, nil, nil, err
		}

		sId, err := generate.SnowflakeId(1)
		if err != nil {
			log.Println("Error generating snowflake ID: ", err)
			return nil, nil, nil, err
		}

		imagePath, err := generate.ImageFromUrl(userInfo.PhotoURL)
		if err != nil {
			log.Println("Error downloading Google profile image: ", err)
			return nil, nil, nil, err
		}

		SplittedName := camelcase.Split(userInfo.DisplayName)

		users := &userModel.User{
			ID:        sId,
			FirstName: SplittedName[0],
			LastName:  SplittedName[1],
			Email:     userInfo.Email,
			Password:  hashedPassword,
			Image:     imagePath,
		}

		if err := s.repository.CreateUser(c, users); err != nil {
			log.Println("Error creating new user: ", err)
			return nil, nil, nil, fmt.Errorf("failed to create user: %v", err)
		}

		user, err = s.repository.GetUserByEmail(c, users.Email)
		if err != nil {
			log.Println("Error fetching newly created user: ", err)
			return nil, nil, nil, fmt.Errorf("failed to fetch user after creation: %v", err)
		}
	}

	if req.FcmToken != "" {
		data := &userModel.User{
			ID:       user.ID,
			FcmToken: req.FcmToken,
		}
		if err := s.repository.UpdateUser(c, data); err != nil {
			log.Println("Failed to update fcm token ", err)
			return nil, nil, nil, fmt.Errorf("failed to update fcm token: %v", err)
		}
	}

	if user.IsBanned {
		return nil, nil, nil, fmt.Errorf("user is banned")
	}

	tokenData := gin.H{
		"id":        strconv.FormatInt(user.ID, 10),
		"role":      user.Role,
		"fcm_token": req.FcmToken,
	}

	accessToken := jwt.GenerateAccessToken(tokenData)
	refreshToken := jwt.GenerateRefreshToken(tokenData)

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Login(c context.Context, req *userRequest.Login) (*userModel.User, any, any, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Println("Failed to validate login request ", err)
		return nil, nil, nil, err
	}

	user, err := s.repository.GetUserByEmail(c, req.Email)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to find user with email: %s", req.Email)
	}

	if !auths.ComparePasswords(user.Password, []byte(req.Password)) {
		return nil, nil, nil, fmt.Errorf("invalid password")
	}

	if user.IsBanned {
		return nil, nil, nil, fmt.Errorf("user is banned")
	}

	if req.FcmToken != "" {
		user.FcmToken = req.FcmToken
		if err := s.repository.UpdateUser(c, user); err != nil {
			log.Println("Failed to update fcm token ", err)
			return nil, nil, nil, fmt.Errorf("failed to update fcm token: %v", err)
		}
	}

	tokenData := gin.H{
		"id":        strconv.FormatInt(user.ID, 10),
		"role":      user.Role,
		"fcm_token": req.FcmToken,
	}

	accessToken := jwt.GenerateAccessToken(tokenData)
	refreshToken := jwt.GenerateRefreshToken(tokenData)

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Register(c context.Context, req *userRequest.Register) (*userModel.User, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Println("Failed to validate register request ", err)
		return nil, err
	}

	user := &userModel.User{}

	_, err := s.repository.GetUserByEmail(c, req.Email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	sId, err := generate.SnowflakeId(1)
	if err != nil {
		log.Println("Failed to generate snowflake id ", err)
		return nil, err
	}

	hashedPassword, err := auths.HashPassword(req.Password)
	if err != nil {
		log.Println("Password failed to be encrypted: ", err)
		return nil, err
	}

	imagePath, err := generate.ImageFromUrl(
		"https://avatar.iran.liara.run/username?username=" + req.FirstName + "+" + req.LastName)
	if err != nil {
		log.Println("Error downloading profile image: ", err)
		return nil, err
	}

	utils.CopyTo(&req, &user)

	user.ID = sId
	user.Password = hashedPassword
	user.Image = imagePath

	if err := s.repository.CreateUser(c, user); err != nil {
		log.Println("Failed to create user ", err)
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (s *UserService) Logout(c context.Context, userId string) error {
	return nil
}

func (s *UserService) BanUser(c context.Context, userId string) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(c, userId)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return nil, fmt.Errorf("user not found: %v", userId)
	}

	if user.IsBanned {
		return nil, fmt.Errorf("user is already banned")
	}

	user.IsBanned = true

	if err := s.repository.UpdateUser(c, user); err != nil {
		log.Println("Failed to ban user ", err)
		return nil, fmt.Errorf("failed to ban user: %v", err)
	}

	return user, nil
}

func (s *UserService) UnbanUser(c context.Context, userId string) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(c, userId)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return nil, fmt.Errorf("user not found: %v", userId)
	}

	if !user.IsBanned {
		return nil, fmt.Errorf("user is not banned")
	}

	user.IsBanned = false

	if err := s.repository.UpdateUser(c, user); err != nil {
		log.Println("Failed to unban user ", err)
		return nil, fmt.Errorf("failed to unban user: %v", err)
	}

	return user, nil
}

func (s *UserService) UpdateProfile(c context.Context, userId string, req *userRequest.UpdateProfile) (*userModel.User, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Println("Failed to validate update profile request ", err)
		return nil, err
	}

	user, err := s.repository.GetUserByID(c, userId)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return nil, fmt.Errorf("user not found: %v", userId)
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
		log.Println("Failed to update user ", err)
		return nil, fmt.Errorf("failed to update user: %v", err)
	}
	return user, nil
}

func (s *UserService) UpdatePassword(c context.Context, userId string, req *userRequest.UpdatePassword) error {
	if err := s.validator.Struct(req); err != nil {
		log.Println("Failed to validate update password request ", err)
		return err
	}

	user, err := s.repository.GetUserByID(c, userId)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return fmt.Errorf("user not found: %v", userId)
	}

	if !auths.ComparePasswords(user.Password, []byte(req.OldPassword)) {
		return fmt.Errorf("invalid password")
	}

	hashedPassword, err := auths.HashPassword(req.NewPassword)
	if err != nil {
		log.Println("Failed to encrypt password ", err)
		return err
	}

	user.Password = hashedPassword

	if err := s.repository.UpdateUser(c, user); err != nil {
		log.Println("Failed to update user ", err)
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func (s *UserService) GetMe(c context.Context, userId string) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(c, userId)
	if err != nil {
		log.Println("Failed to get profile information", err)
		return nil, fmt.Errorf("user not found: %v", userId)
	}

	return user, nil
}

func (s *UserService) GetUsers(c context.Context) ([]*userModel.User, error) {
	user, err := s.repository.GetUsers(c)
	if err != nil {
		log.Println("Failed to get users ", err)
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	return user, nil
}

func (s *UserService) GetBannedUsers(c context.Context) ([]*userModel.User, error) {
	user, err := s.repository.GetBannedUsers(c)
	if err != nil {
		log.Println("Failed to get banned users ", err)
		return nil, fmt.Errorf("failed to get banned users: %v", err)
	}

	return user, nil
}

func (s *UserService) GetUserByID(c context.Context, userId string) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(c, userId)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return nil, fmt.Errorf("user not found: %v", userId)
	}

	return user, nil
}
