package userService

import (
	"context"
	"fmt"
	"log"
	"strconv"

	userModel "washit-api/app/user/dto/model"
	userRequest "washit-api/app/user/dto/request"
	userRepository "washit-api/app/user/repository"
	jwt "washit-api/token"
	"washit-api/utils"

	"firebase.google.com/go/auth"
	"github.com/fatih/camelcase"
	"github.com/gin-gonic/gin"
)

type UserServiceInterface interface {
	RefreshToken(ctx context.Context, userId string) (string, error)
	Register(ctx context.Context, req *userRequest.Register) (*userModel.User, error)
	LoginWithGoogle(ctx context.Context, req *userRequest.Google, userInfo *auth.UserInfo) (*userModel.User, any, any, error)
	Login(ctx context.Context, req *userRequest.Login) (*userModel.User, any, any, error)
	Logout(ctx context.Context, userId string) error
	GetMe(ctx context.Context, userId string) (*userModel.User, error)
	GetUserByID(ctx context.Context, userId string) (*userModel.User, error)
	GetUsers(ctx context.Context) ([]*userModel.User, error)
	UpdateMe(ctx context.Context, userId string, req *userRequest.Update) (*userModel.User, error)
}

type UserService struct {
	repository userRepository.UserRepositoryInterface
}

func NewUserService(
	repository userRepository.UserRepositoryInterface) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) RefreshToken(ctx context.Context, userId string) (string, error) {
	user, err := s.repository.GetUserByID(ctx, userId)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return "", fmt.Errorf("user not found: %v", userId)
	}

	tokenData := gin.H{"id": strconv.FormatInt(user.ID, 10), "role": user.Role}

	accessToken := jwt.GenerateAccessToken(tokenData)
	return accessToken, nil
}

func (s *UserService) LoginWithGoogle(ctx context.Context, req *userRequest.Google, userInfo *auth.UserInfo) (*userModel.User, any, any, error) {
	var user *userModel.User
	user, err := s.repository.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		log.Println("Error fetching user: ", err)
		return nil, nil, nil, fmt.Errorf("unable to find user with email: %s", userInfo.Email)
	}

	if user == nil {
		hashedPassword, err := utils.HashPassword(userInfo.UID)
		if err != nil {
			log.Println("Error encrypting password: ", err)
			return nil, nil, nil, err
		}

		sId, err := utils.SnowflakeId(1)
		if err != nil {
			log.Println("Error generating snowflake ID: ", err)
			return nil, nil, nil, err
		}

		SplittedName := camelcase.Split(userInfo.DisplayName)

		imagePath, err := utils.TakeGoogleImage(userInfo.PhotoURL)
		if err != nil {
			log.Println("Error downloading Google profile image: ", err)
			return nil, nil, nil, err
		}

		users := &userModel.User{
			ID:        sId,
			FirstName: SplittedName[0],
			LastName:  SplittedName[1],
			Email:     userInfo.Email,
			Password:  hashedPassword,
			Image:     imagePath,
		}

		if err := s.repository.CreateUser(ctx, users); err != nil {
			log.Println("Error creating new user: ", err)
			return nil, nil, nil, fmt.Errorf("failed to create user: %v", err)
		}

		user, err = s.repository.GetUserByEmail(ctx, users.Email)
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
		if err := s.repository.UpdateUser(ctx, data); err != nil {
			log.Println("Failed to update fcm token ", err)
			return nil, nil, nil, fmt.Errorf("failed to update fcm token: %v", err)
		}
	}

	tokenData := gin.H{"id": strconv.FormatInt(user.ID, 10), "role": user.Role}
	accessToken := jwt.GenerateAccessToken(tokenData)
	refreshToken := jwt.GenerateRefreshToken(tokenData)
	return user, accessToken, refreshToken, nil
}

func (s *UserService) Login(ctx context.Context, req *userRequest.Login) (*userModel.User, any, any, error) {
	user, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to find user with email: %s", req.Email)
	}

	if !utils.ComparePasswords(user.Password, []byte(req.Password)) {
		return nil, nil, nil, err
	}

	if req.FcmToken != "" {
		data := &userModel.User{
			ID:       user.ID,
			FcmToken: req.FcmToken,
		}
		if err := s.repository.UpdateUser(ctx, data); err != nil {
			log.Println("Failed to update fcm token ", err)
			return nil, nil, nil, fmt.Errorf("failed to update fcm token: %v", err)
		}
	}

	tokenData := gin.H{"id": strconv.FormatInt(user.ID, 10), "role": user.Role}

	accessToken := jwt.GenerateAccessToken(tokenData)
	refreshToken := jwt.GenerateRefreshToken(tokenData)
	return user, accessToken, refreshToken, nil
}

func (s *UserService) Register(ctx context.Context, req *userRequest.Register) (*userModel.User, error) {
	_, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println("Password failed to be encrypted: ", err)
		return nil, err
	}

	sId, err := utils.SnowflakeId(1)
	if err != nil {
		log.Println("Failed to generate snowflake id ", err)
		return nil, err
	}

	imagePath, err := utils.MakeProfileImage(req.FirstName, req.LastName)
	if err != nil {
		log.Println("Failed to create profile image ", err)
		return nil, err
	}

	user := &userModel.User{
		ID:        sId,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
		Image:     imagePath,
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		log.Println("Failed to create user ", err)
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (s *UserService) Logout(ctx context.Context, userId string) error {
	return nil
}

func (s *UserService) UpdateMe(ctx context.Context, userId string, req *userRequest.Update) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(ctx, userId)
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
	if req.Password != "" {
		user.Password, err = utils.HashPassword(req.Password)
		if err != nil {
			log.Println("Password failed to be encrypted: ", err)
			return nil, err
		}
	}

	if err := s.repository.UpdateUser(ctx, user); err != nil {
		log.Println("Failed to update user ", err)
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}

func (s *UserService) GetMe(ctx context.Context, userId string) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(ctx, userId)
	if err != nil {
		log.Println("Failed to get profile information", err)
		return nil, fmt.Errorf("user not found: %v", userId)
	}

	return user, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]*userModel.User, error) {
	user, err := s.repository.GetUsers(ctx)
	if err != nil {
		log.Println("Failed to get users ", err)
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userId string) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(ctx, userId)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return nil, fmt.Errorf("user not found: %v", userId)
	}

	return user, nil
}
