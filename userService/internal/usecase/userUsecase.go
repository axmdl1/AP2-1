package usecase

import (
	"AP-1/userService/internal/entity"
	"AP-1/userService/internal/repository"
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, username, email, password string) (string, error)
	AuthenticateUser(ctx context.Context, email, password string) (string, error)
	GetUserProfile(ctx context.Context, userID string) (*entity.UserProfile, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u *userUsecase) RegisterUser(ctx context.Context, username, email, password string) (string, error) {
	existing, _ := u.repo.FindByEmail(ctx, email)
	if existing != nil {
		return "", errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUUID := uuid.New().String()

	user := entity.User{
		ID:           newUUID,
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := u.repo.Create(ctx, &user); err != nil {
		return "", err
	}
	return user.ID, nil
}

func (u *userUsecase) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	//TODO: implement token methods.
	return "JWT-TOKEN", err
}

func (u *userUsecase) GetUserProfile(ctx context.Context, userID string) (*entity.UserProfile, error) {
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	profile := &entity.UserProfile{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	return profile, nil
}
