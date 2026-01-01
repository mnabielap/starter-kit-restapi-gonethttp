package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"starter-kit-restapi-gonethttp/internal/services"
	"starter-kit-restapi-gonethttp/pkg/response"
	"starter-kit-restapi-gonethttp/pkg/utils"

	"github.com/google/uuid"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req services.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if errs := utils.ValidateStruct(req); errs != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "Validation error", "errors": errs})
		return
	}
	user, err := h.service.CreateUser(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, user)
}

// Implemented GetUsers with Query Params
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 {
		limit = 10
	}

	// Extract Sort Param (e.g., "created_at:desc")
	sortBy := query.Get("sortBy")

	// Extract Search Params
	search := query.Get("search")
	scope := query.Get("scope")
	role := query.Get("role")

	filters := map[string]interface{}{
		"search": search,
		"scope":  scope,
		"role":   role,
	}

	result, err := h.service.GetUsers(filters, page, limit, sortBy)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}
	response.Success(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	var req services.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "Validation error", "errors": errs})
		return
	}

	user, err := h.service.UpdateUser(id, req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	if err := h.service.DeleteUser(id); err != nil {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}
	response.Success(w, http.StatusNoContent, nil)
}