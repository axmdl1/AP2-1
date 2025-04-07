package usecase

import (
	"AP-1/userService/internal/entity"
	"AP-1/userService/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, email, password string) (*entity.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u userUsecase) Register(ctx context.Context, user *entity.User) error {
	existing, _ := u.repo.FindByEmail(ctx, user.Email)
	if existing != nil {
		return errors.New("user with this email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return u.repo.Create(ctx, user)
}

func (u userUsecase) Login(ctx context.Context, email, password string) (*entity.User, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("wrong password")
	}

	return user, nil
}
