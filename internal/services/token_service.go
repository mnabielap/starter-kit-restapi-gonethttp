package services

import (
	"time"

	"starter-kit-restapi-gonethttp/config"
	"starter-kit-restapi-gonethttp/internal/models"
	"starter-kit-restapi-gonethttp/internal/repository"
	"starter-kit-restapi-gonethttp/pkg/utils"
)

type TokenService struct {
	repo repository.TokenRepository
	cfg  *config.Config
}

func NewTokenService(repo repository.TokenRepository, cfg *config.Config) *TokenService {
	return &TokenService{repo: repo, cfg: cfg}
}

func (s *TokenService) GenerateAuthTokens(user *models.User) (map[string]interface{}, error) {
	accessToken, refreshToken, accessExp, refreshExp, err := utils.GenerateAuthTokens(user.ID, s.cfg)
	if err != nil {
		return nil, err
	}

	// Save refresh token to database
	err = s.SaveToken(refreshToken, user.ID.String(), refreshExp, models.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"access": map[string]interface{}{
			"token":   accessToken,
			"expires": accessExp,
		},
		"refresh": map[string]interface{}{
			"token":   refreshToken,
			"expires": refreshExp,
		},
	}, nil
}

func (s *TokenService) SaveToken(token, userID string, expires time.Time, tokenType string) error {
	tokenModel := &models.Token{
		Token:   token,
		UserID:  userID,
		Expires: expires,
		Type:    tokenType,
	}
	return s.repo.Create(tokenModel)
}

func (s *TokenService) VerifyToken(token string, tokenType string) (*models.Token, error) {
	return s.repo.FindByToken(token, tokenType)
}