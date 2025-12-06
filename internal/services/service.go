package services

import (
	"starter-kit-restapi-gonethttp/internal/models"
	"starter-kit-restapi-gonethttp/pkg/utils"

	"github.com/google/uuid"
)

// AuthService defines the interface for authentication logic
type AuthService interface {
	Login(email, password string) (*models.User, map[string]interface{}, error)
	Register(req RegisterRequest) (*models.User, map[string]interface{}, error)
	RefreshAuth(refreshToken string) (map[string]interface{}, error)
	Logout(refreshToken string) error
	
	// Password Reset & Verification
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	SendVerificationEmail(user *models.User) error
	VerifyEmail(token string) error
}

// UserService defines the interface for user management logic
type UserService interface {
	CreateUser(req CreateUserRequest) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	// Updated signature to match the implementation using PaginationResult
	GetUsers(filters map[string]interface{}, page, limit int, sort string) (*utils.PaginationResult, error)
	UpdateUser(id uuid.UUID, req UpdateUserRequest) (*models.User, error)
	DeleteUser(id uuid.UUID) error
}

// DTOs (Data Transfer Objects) for Requests
type RegisterRequest struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

type CreateUserRequest struct {
	RegisterRequest
	Role string `validate:"required,oneof=user admin"`
}

type UpdateUserRequest struct {
	Name     string `validate:"omitempty"`
	Email    string `validate:"omitempty,email"`
	Password string `validate:"omitempty,min=8"`
}