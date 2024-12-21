package userRepository

import (
	"context"

	userModel "washit-api/app/user/dto/model"
	dbs "washit-api/db"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *userModel.User) error
	GetUserByID(ctx context.Context, userId string) (*userModel.User, error)
	GetUserByEmail(ctx context.Context, email string) (*userModel.User, error)
	PutFcmToken(ctx context.Context, userId int64, fcmToken string) error
	GetUsers(ctx context.Context) ([]*userModel.User, error)
	GetBannedUsers(ctx context.Context) ([]*userModel.User, error)
	UpdateUser(ctx context.Context, user *userModel.User) error
}

type UserRepository struct {
	db dbs.DatabaseInterface
}

func NewUserRepository(db dbs.DatabaseInterface) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *userModel.User) error {
	return r.db.Create(ctx, user)
}

func (r *UserRepository) PutFcmToken(ctx context.Context, userId int64, fcmToken string) error {
	return r.db.Update(ctx, &userModel.User{ID: userId, FcmToken: fcmToken})
}

func (r *UserRepository) GetUserByID(ctx context.Context, userId string) (*userModel.User, error) {
	var user userModel.User
	if err := r.db.FindById(ctx, userId, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
	var user userModel.User
	query := dbs.NewQuery("email = ?", email)
	if err := r.db.FindOne(ctx, &user, dbs.WithQuery(query)); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]*userModel.User, error) {
	var users []*userModel.User
	if err := r.db.Find(ctx, &users, dbs.WithLimit(10), dbs.WithOrder("id")); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetBannedUsers(ctx context.Context) ([]*userModel.User, error) {
	var users []*userModel.User
	query := dbs.NewQuery("is_banned = ?", true)
	if err := r.db.Find(ctx, &users, dbs.WithQuery(query)); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *userModel.User) error {
	return r.db.Update(ctx, user)
}
