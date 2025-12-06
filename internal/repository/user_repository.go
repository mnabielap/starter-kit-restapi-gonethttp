package repository

import (
	"strings"

	"starter-kit-restapi-gonethttp/internal/models"
	"starter-kit-restapi-gonethttp/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Implemented FindAll with pagination and filtering
func (r *userRepository) FindAll(filters map[string]interface{}, pagination *utils.PaginationScope) ([]models.User, int64, error) {
	var users []models.User
	var totalRows int64

	query := r.db.Model(&models.User{})

	// Apply Filters
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("lower(name) LIKE ?", "%"+strings.ToLower(name)+"%")
	}
	if role, ok := filters["role"].(string); ok && role != "" {
		query = query.Where("role = ?", role)
	}

	// Count total rows for pagination
	query.Count(&totalRows)

	// Apply Sorting
	if pagination.Sort != "" {
		query = query.Order(pagination.Sort)
	} else {
		query = query.Order("created_at desc")
	}

	// Apply Pagination
	err := query.Scopes(pagination.Paginate()).Find(&users).Error
	return users, totalRows, err
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}