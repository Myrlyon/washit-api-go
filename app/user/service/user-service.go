package userService

import (
	"context"
	"errors"
	"log"
	"strconv"

	userModel "washit-api/app/user/dto/model"
	userRequest "washit-api/app/user/dto/request"
	userRepository "washit-api/app/user/repository"
	jwt "washit-api/token"
	"washit-api/utils"
)

type UserServiceInterface interface {
	RefreshToken(ctx context.Context, userId string) (string, error)
	Register(ctx context.Context, req *userRequest.Register) (*userModel.User, error)
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
		return "", err
	}

	tokenData := map[string]interface{}{"id": strconv.FormatInt(user.ID, 10), "role": user.Role}

	accessToken := jwt.GenerateAccessToken(tokenData)
	return accessToken, nil
}

func (s *UserService) Login(ctx context.Context, req *userRequest.Login) (*userModel.User, any, any, error) {
	user, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, nil, errors.New("invalid email or password")
	}

	if !utils.ComparePasswords(user.Password, []byte(req.Password)) {
		return nil, nil, nil, errors.New("invalid email or password")
	}

	tokenData := map[string]interface{}{"id": strconv.FormatInt(user.ID, 10), "role": user.Role}

	accessToken := jwt.GenerateAccessToken(tokenData)
	refreshToken := jwt.GenerateRefreshToken(tokenData)
	return user, accessToken, refreshToken, nil
}

func (s *UserService) Register(ctx context.Context, req *userRequest.Register) (*userModel.User, error) {
	_, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already in use")
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetMe(ctx context.Context, userId string) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(ctx, userId)
	if err != nil {
		log.Println("Failed to get profile information", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]*userModel.User, error) {
	user, err := s.repository.GetUsers(ctx)
	if err != nil {
		log.Println("Failed to get users ", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userId string) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(ctx, userId)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return nil, err
	}

	return user, nil
}
