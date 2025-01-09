package userService

import (
	"context"
	"errors"
	"testing"
	userModel "washit-api/internal/user/dto/model"
	userRequest "washit-api/internal/user/dto/request"
	mocks "washit-api/internal/user/repository/mock"
	auths "washit-api/pkg/auth"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.IUserRepository
	service  IUserService
}

func (suite *UserServiceTestSuite) SetupTest() {
	validator := validator.New()
	suite.mockRepo = new(mocks.IUserRepository)
	suite.service = NewUserService(suite.mockRepo, validator)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// Login
// =================================================================

func (suite *UserServiceTestSuite) TestLoginGetUserByEmailFail() {
	req := &userRequest.Login{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(nil, errors.New("error")).Times(1)

	user, accessToken, refreshToken, err := suite.service.Login(context.Background(), req)
	suite.Nil(user)
	suite.Empty(accessToken)
	suite.Empty(refreshToken)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestLoginInvalidEmailFormat() {
	req := &userRequest.Login{
		Email:    "email",
		Password: "test123456",
	}

	user, accessToken, refreshToken, err := suite.service.Login(context.Background(), req)
	suite.Nil(user)
	suite.Empty(accessToken)
	suite.Empty(refreshToken)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestLoginWrongPassword() {
	req := &userRequest.Login{
		Email:    "test@test.com",
		Password: "test123456",
	}

	hashedPassword, err := auths.HashPassword("password")
	suite.NoError(err)

	suite.mockRepo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(&userModel.User{
			Email:    "test@test.com",
			Password: hashedPassword,
		}, nil).Times(1)

	user, accessToken, refreshToken, err := suite.service.Login(context.Background(), req)
	suite.Nil(user)
	suite.Empty(accessToken)
	suite.Empty(refreshToken)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestLoginSuccess() {
	req := &userRequest.Login{
		Email:    "test@test.com",
		Password: "test123456",
	}
	hashedPassword, err := auths.HashPassword("test123456")
	suite.NoError(err)

	suite.mockRepo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(
			&userModel.User{
				Email:    "test@test.com",
				Password: hashedPassword,
			},
			nil,
		).Times(1)

	user, accessToken, refreshToken, err := suite.service.Login(context.Background(), req)
	suite.NotNil(user)
	suite.Equal(req.Email, user.Email)
	suite.NotEmpty(accessToken)
	suite.NotEmpty(refreshToken)
	suite.Nil(err)
}

// Register
// =================================================================
func (suite *UserServiceTestSuite) TestRegisterSuccess() {
	req := &userRequest.Register{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@test.com",
		Password:  "test123456",
	}

	suite.mockRepo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(nil, errors.New("user not found")).Times(1)

	suite.mockRepo.On("CreateUser", mock.Anything, mock.Anything).
		Return(nil).Times(1)

	user, err := suite.service.Register(context.Background(), req)
	suite.NotNil(user)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestRegisterCreateUserFail() {
	req := &userRequest.Register{
		// FirstName: "John",
		// LastName: "Doe",
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	user, err := suite.service.Register(context.Background(), req)
	suite.Nil(user)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestRegisterInvalidEmailFormat() {
	req := &userRequest.Register{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "email",
		Password:  "test123456",
	}
	user, err := suite.service.Register(context.Background(), req)
	suite.Nil(user)
	suite.NotNil(err)
}

// UpdateProfile
// =================================================================

func (suite *UserServiceTestSuite) TestUpdateProfileSuccess() {
	req := &userRequest.UpdateProfile{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@test.com",
	}

	suite.mockRepo.On("GetUserByID", mock.Anything, mock.Anything).
		Return(&userModel.User{}, nil).Times(1)

	suite.mockRepo.On("UpdateUser", mock.Anything, mock.Anything).
		Return(nil).Times(1)

	user, err := suite.service.UpdateProfile(context.Background(), 0, req)
	suite.NotNil(user)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestUpdateProfileGetUserByIDFail() {
	req := &userRequest.UpdateProfile{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@test.com",
	}

	suite.mockRepo.On("GetUserByID", mock.Anything, mock.Anything).
		Return(nil, errors.New("error")).Times(1)

	user, err := suite.service.UpdateProfile(context.Background(), 0, req)
	suite.Nil(user)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestUpdateProfileUpdateUserFail() {
	req := &userRequest.UpdateProfile{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@test.com",
	}

	suite.mockRepo.On("GetUserByID", mock.Anything, mock.Anything).
		Return(&userModel.User{}, nil).Times(1)

	suite.mockRepo.On("UpdateUser", mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	user, err := suite.service.UpdateProfile(context.Background(), 0, req)
	suite.Nil(user)
	suite.NotNil(err)
}

// UpdatePassword
// =================================================================

func (suite *UserServiceTestSuite) TestUpdatePasswordSuccess() {
	req := &userRequest.UpdatePassword{
		OldPassword:     "test123456",
		NewPassword:     "test1234567",
		ConfirmPassword: "test1234567",
	}

	hashedPassword, err := auths.HashPassword("test123456")
	suite.NoError(err)

	suite.mockRepo.On("GetUserByID", mock.Anything, mock.Anything).
		Return(&userModel.User{Password: hashedPassword}, nil).Times(1)

	suite.mockRepo.On("UpdateUser", mock.Anything, mock.Anything).
		Return(nil).Times(1)

	err = suite.service.UpdatePassword(context.Background(), 0, req)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestUpdatePasswordMinPasswordLength() {
	req := &userRequest.UpdatePassword{
		OldPassword:     "test123456",
		NewPassword:     "test",
		ConfirmPassword: "test",
	}

	hashedPassword, err := auths.HashPassword("test123456")
	suite.NoError(err)

	suite.mockRepo.On("GetUserByID", mock.Anything, mock.Anything).
		Return(&userModel.User{Password: hashedPassword}, nil).Times(1)

	suite.mockRepo.On("UpdateUser", mock.Anything, mock.Anything).
		Return(nil).Times(1)

	err = suite.service.UpdatePassword(context.Background(), 0, req)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestUpdatePasswordGetUserByIDFail() {
	req := &userRequest.UpdatePassword{
		OldPassword:     "test123456",
		NewPassword:     "test1234567",
		ConfirmPassword: "test1234567",
	}

	suite.mockRepo.On("GetUserByID", mock.Anything, mock.Anything).
		Return(nil, errors.New("error")).Times(1)

	err := suite.service.UpdatePassword(context.Background(), 0, req)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestUpdatePasswordWrongOldPassword() {
	req := &userRequest.UpdatePassword{
		OldPassword:     "test123456",
		NewPassword:     "test1234567",
		ConfirmPassword: "test1234567",
	}

	hashedPassword, err := auths.HashPassword("oldpassword")
	suite.NoError(err)

	suite.mockRepo.On("GetUserByID", mock.Anything, mock.Anything).
		Return(&userModel.User{Password: hashedPassword}, nil).Times(1)

	err = suite.service.UpdatePassword(context.Background(), 0, req)
	suite.NotNil(err)
}
