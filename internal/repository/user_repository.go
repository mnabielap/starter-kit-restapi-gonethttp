package repository

import (
	"fmt"
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

func (r *userRepository) FindAll(filters map[string]interface{}, pagination *utils.PaginationScope) ([]models.User, int64, error) {
	var users []models.User
	var totalRows int64

	query := r.db.Model(&models.User{})

	// --- 1. SEARCH LOGIC ---
	if search, ok := filters["search"].(string); ok && search != "" {
		scope, _ := filters["scope"].(string)
		searchPattern := "%" + strings.ToLower(search) + "%"

		switch scope {
		case "name":
			query = query.Where("lower(name) LIKE ?", searchPattern)
		case "email":
			query = query.Where("lower(email) LIKE ?", searchPattern)
		case "id":
			// Strict ID search
			if _, err := uuid.Parse(search); err == nil {
				query = query.Where("id = ?", search)
			} else {
				// If scope is ID but invalid UUID provided, return nothing
				query = query.Where("1 = 0")
			}
		case "all":
			fallthrough
		default:
			// OR Logic: Name OR Email OR ID (if valid UUID)
			subQuery := r.db.Where("lower(name) LIKE ?", searchPattern).
				Or("lower(email) LIKE ?", searchPattern)

			if _, err := uuid.Parse(search); err == nil {
				subQuery = subQuery.Or("id = ?", search)
			}
			query = query.Where(subQuery)
		}
	}

	// --- 2. FILTER LOGIC ---
	if role, ok := filters["role"].(string); ok && role != "" {
		query = query.Where("role = ?", role)
	}

	// --- 3. COUNT TOTAL ---
	query.Count(&totalRows)

	// --- 4. SORTING LOGIC ---
	// Parse "field:order" (e.g., "created_at:desc")
	sortParam := pagination.Sort
	orderClause := "created_at desc" // Default

	if sortParam != "" {
		parts := strings.Split(sortParam, ":")
		field := parts[0]
		direction := "asc"
		if len(parts) > 1 && strings.ToLower(parts[1]) == "desc" {
			direction = "desc"
		}

		// Whitelist allowed fields to prevent SQL injection
		allowedFields := map[string]bool{
			"id":         true,
			"name":       true,
			"email":      true,
			"role":       true,
			"created_at": true,
		}

		if allowedFields[field] {
			// Standard sorting
			orderClause = fmt.Sprintf("%s %s", field, direction)
			query = query.Order(orderClause)
		} else {
			// Fallback if invalid field
			query = query.Order(orderClause)
		}
	} else {
		query = query.Order(orderClause)
	}

	// --- 5. PAGINATION ---
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