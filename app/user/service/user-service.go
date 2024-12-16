package service

import (
	"context"
	"errors"
	"log"
	"strconv"

	"washit-api/app/user/dto/request"
	"washit-api/app/user/model"
	"washit-api/app/user/repository"
	jwt "washit-api/token"
	"washit-api/utils"
)

type UserServiceInterface interface {
	Register(ctx context.Context, req *request.Register) (*model.User, error)
	Login(ctx context.Context, req *request.Login) (*model.User, any, any, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

type UserService struct {
	repository repository.UserRepositoryInterface
}

func NewUserService(
	repository repository.UserRepositoryInterface) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Login(ctx context.Context, r *request.Login) (*model.User, any, any, error) {
	user, err := s.repository.GetUserByEmail(ctx, r.Email)
	if err != nil {
		return nil, nil, nil, errors.New("invalid email or password")
	}

	same := utils.ComparePasswords(user.Password, []byte(r.Password))
	if !same {
		return nil, nil, nil, errors.New("invalid email or password")
	}

	tokenData := map[string]interface{}{
		"id":    strconv.Itoa(user.ID),
		"email": user.Email,
	}
	accessToken := jwt.GenerateAccessToken(tokenData)
	refreshToken := jwt.GenerateRefreshToken(tokenData)
	return user, accessToken, refreshToken, nil
}

func (s *UserService) Register(ctx context.Context, r *request.Register) (*model.User, error) {
	_, err := s.repository.GetUserByEmail(ctx, r.Email)
	if err == nil {
		return nil, errors.New("email already in use")
	}

	hashedPassword, err := utils.HashPassword(r.Password)
	if err != nil {
		log.Println("Password failed to be encrypted: ", err)
		return nil, err
	}

	user := &model.User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  hashedPassword,
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		log.Println("Failed to create user ", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		log.Println("Failed to get user by id ", err)
		return nil, err
	}

	return user, nil
}
