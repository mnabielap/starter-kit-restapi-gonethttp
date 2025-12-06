package services

import (
	"errors"
	"time"

	"starter-kit-restapi-gonethttp/config"
	"starter-kit-restapi-gonethttp/internal/models"
	"starter-kit-restapi-gonethttp/internal/repository"
	"starter-kit-restapi-gonethttp/pkg/utils"

	"github.com/google/uuid"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenRepo    repository.TokenRepository
	tokenService *TokenService
	emailService EmailService
	cfg          *config.Config
}

func NewAuthService(uRepo repository.UserRepository, tRepo repository.TokenRepository, tService *TokenService, eService EmailService, cfg *config.Config) AuthService {
	return &authService{
		userRepo:     uRepo,
		tokenRepo:    tRepo,
		tokenService: tService,
		emailService: eService,
		cfg:          cfg,
	}
}


func (s *authService) Login(email, password string) (*models.User, map[string]interface{}, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || !user.ComparePassword(password) {
		return nil, nil, errors.New("incorrect email or password")
	}
	tokens, err := s.tokenService.GenerateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (s *authService) Register(req RegisterRequest) (*models.User, map[string]interface{}, error) {
	if exists, _ := s.userRepo.ExistsByEmail(req.Email); exists {
		return nil, nil, errors.New("email already taken")
	}
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}
	tokens, err := s.tokenService.GenerateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (s *authService) Logout(refreshToken string) error {
	tokenDoc, err := s.tokenService.VerifyToken(refreshToken, models.TokenTypeRefresh)
	if err != nil {
		return errors.New("not found")
	}
	return s.tokenRepo.Delete(tokenDoc)
}

func (s *authService) RefreshAuth(refreshToken string) (map[string]interface{}, error) {
	tokenDoc, err := s.tokenService.VerifyToken(refreshToken, models.TokenTypeRefresh)
	if err != nil {
		return nil, errors.New("please authenticate")
	}
	payload, err := utils.ValidateToken(refreshToken, s.cfg.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	userUUID, _ := uuid.Parse(payload.Sub)
	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	s.tokenRepo.Delete(tokenDoc)
	return s.tokenService.GenerateAuthTokens(user)
}

func (s *authService) ForgotPassword(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Return nil to avoid email enumeration
		return nil 
	}

	expires := time.Duration(s.cfg.JWT.ResetPasswordExpirationMinutes) * time.Minute
	resetToken, _, err := utils.GenerateToken(user.ID, expires, models.TokenTypeResetPassword, s.cfg.JWT.Secret)
	if err != nil {
		return err
	}

	// Save token
	err = s.tokenService.SaveToken(resetToken, user.ID.String(), time.Now().Add(expires), models.TokenTypeResetPassword)
	if err != nil {
		return err
	}

	return s.emailService.SendResetPasswordEmail(user.Email, resetToken)
}

func (s *authService) ResetPassword(tokenStr, newPassword string) error {
	tokenDoc, err := s.tokenService.VerifyToken(tokenStr, models.TokenTypeResetPassword)
	if err != nil {
		return errors.New("password reset failed")
	}

	userUUID, err := uuid.Parse(tokenDoc.UserID)
	if err != nil {
		return errors.New("invalid user data")
	}

	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return errors.New("password reset failed")
	}

	user.Password = newPassword
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Consume token (delete all reset tokens for this user)
	return s.tokenRepo.DeleteByUserIDAndType(user.ID.String(), models.TokenTypeResetPassword)
}

func (s *authService) SendVerificationEmail(user *models.User) error {
	expires := time.Duration(s.cfg.JWT.VerifyEmailExpirationMinutes) * time.Minute
	verifyToken, _, err := utils.GenerateToken(user.ID, expires, models.TokenTypeVerifyEmail, s.cfg.JWT.Secret)
	if err != nil {
		return err
	}

	err = s.tokenService.SaveToken(verifyToken, user.ID.String(), time.Now().Add(expires), models.TokenTypeVerifyEmail)
	if err != nil {
		return err
	}

	return s.emailService.SendVerificationEmail(user.Email, verifyToken)
}

func (s *authService) VerifyEmail(tokenStr string) error {
	tokenDoc, err := s.tokenService.VerifyToken(tokenStr, models.TokenTypeVerifyEmail)
	if err != nil {
		return errors.New("email verification failed")
	}

	userUUID, err := uuid.Parse(tokenDoc.UserID)
	if err != nil {
		return errors.New("invalid user data")
	}

	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return errors.New("email verification failed")
	}

	user.IsEmailVerified = true
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return s.tokenRepo.DeleteByUserIDAndType(user.ID.String(), models.TokenTypeVerifyEmail)
}