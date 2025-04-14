package handler

import (
	userservice "AP-1/pb/userService"
	"AP-1/userService/internal/usecase"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type UserServiceServer struct {
	userservice.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserServiceServer(u usecase.UserUsecase) *UserServiceServer {
	return &UserServiceServer{
		usecase: u,
	}
}

func (s *UserServiceServer) RegisterUser(ctx context.Context, req *userservice.RegisterUserRequest) (*userservice.RegisterUserResponse, error) {
	log.Println("RegisterUser: received request: ", req)
	userID, err := s.usecase.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &userservice.RegisterUserResponse{
		UserID:  userID,
		Message: "User registered successfully",
	}, nil
}

func (s *UserServiceServer) AuthenticateUser(ctx context.Context, req *userservice.AuthenticateUserRequest) (*userservice.AuthenticateUserResponse, error) {
	log.Println("AuthenticateUser: received request: ", req)
	token, err := s.usecase.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &userservice.AuthenticateUserResponse{
		Token:   token,
		UserId:  token,
		Message: "User authenticated successfully",
	}, nil
}

func (s *UserServiceServer) GetUserProfile(ctx context.Context, req *userservice.GetUserProfileRequest) (*userservice.UserProfile, error) {
	log.Println("GetUserProfile: received request: ", req)
	profile, err := s.usecase.GetUserProfile(ctx, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &userservice.UserProfile{
		UserID:   profile.UserID,
		Email:    profile.Email,
		Username: profile.Username,
	}, nil
}
