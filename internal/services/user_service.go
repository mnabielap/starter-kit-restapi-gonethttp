package services

import (
	"errors"

	"starter-kit-restapi-gonethttp/internal/models"
	"starter-kit-restapi-gonethttp/internal/repository"
	"starter-kit-restapi-gonethttp/pkg/utils"

	"github.com/google/uuid"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req CreateUserRequest) (*models.User, error) {
	if exists, _ := s.repo.ExistsByEmail(req.Email); exists {
		return nil, errors.New("email already taken")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetUsers(filters map[string]interface{}, page, limit int, sort string) (*utils.PaginationResult, error) {
	paginationScope := &utils.PaginationScope{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	users, totalRows, err := s.repo.FindAll(filters, paginationScope)
	if err != nil {
		return nil, err
	}

	result := utils.GetPaginationResult(totalRows, page, limit, users)
	return &result, nil
}

func (s *userService) UpdateUser(id uuid.UUID, req UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Email != "" && req.Email != user.Email {
		if exists, _ := s.repo.ExistsByEmail(req.Email); exists {
			return nil, errors.New("email already taken")
		}
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Password != "" {
		user.Password = req.Password
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(id)
}