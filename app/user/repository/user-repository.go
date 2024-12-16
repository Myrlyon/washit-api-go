package repository

import (
	"context"

	"washit-api/app/user/model"
	dbs "washit-api/db"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepository struct {
	db dbs.DatabaseInterface
}

func NewUserRepository(db dbs.DatabaseInterface) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.Create(ctx, user)
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.FindById(ctx, id, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := dbs.NewQuery("email = ?", email)
	if err := r.db.FindOne(ctx, &user, dbs.WithQuery(query)); err != nil {
		return nil, err
	}

	return &user, nil
}
