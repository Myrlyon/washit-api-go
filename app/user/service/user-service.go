package userService

import (
	"context"
	"errors"
	"log"
	"strconv"

	userRequest "washit-api/app/user/dto/request"
	userModel "washit-api/app/user/model"
	userRepository "washit-api/app/user/repository"
	jwt "washit-api/token"
	"washit-api/utils"
)

type UserServiceInterface interface {
	Register(ctx context.Context, req *userRequest.Register) (*userModel.User, error)
	Login(ctx context.Context, req *userRequest.Login) (*userModel.User, any, error)
	GetUserByID(ctx context.Context, id string) (*userModel.User, error)
	GetUsers(ctx context.Context) ([]*userModel.User, error)
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

func (s *UserService) Login(ctx context.Context, req *userRequest.Login) (*userModel.User, any, error) {
	user, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	if !utils.ComparePasswords(user.Password, []byte(req.Password)) {
		return nil, nil, errors.New("invalid email or password")
	}

	tokenData := map[string]interface{}{"id": strconv.Itoa(user.ID), "role": user.Role}

	accessToken := jwt.GenerateAccessToken(tokenData)
	return user, accessToken, nil
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

	user := &userModel.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		log.Println("Failed to create user ", err)
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

func (s *UserService) GetUserByID(ctx context.Context, id string) (*userModel.User, error) {
	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return nil, err
	}

	return user, nil
}
