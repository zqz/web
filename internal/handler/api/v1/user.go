package v1

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zqz/web/backend/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userSvc *service.UserService
	fileSvc *service.FileService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userSvc *service.UserService, fileSvc *service.FileService) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
		fileSvc: fileSvc,
	}
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	DisplayTag string `json:"display_tag,omitempty"`
	Colour     string `json:"colour,omitempty"`
}

// UserWithFilesResponse represents a user with their files
type UserWithFilesResponse struct {
	UserResponse
	FileCount int `json:"file_count"`
}

// ListUsers returns a list of all users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := int32(min(100, maxListLimit))
	users, err := h.userSvc.ListUsers(r.Context(), limit, 0)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	response := make([]UserWithFilesResponse, len(users))
	fileLimit := int32(min(1000, maxListLimit))
	for i, user := range users {
		// Get file count for each user
		files, _ := h.fileSvc.ListFilesByUserID(r.Context(), user.ID, fileLimit, 0)
		fileCount := 0
		if files != nil {
			fileCount = len(files)
		}

		response[i] = UserWithFilesResponse{
			UserResponse: UserResponse{
				ID:         user.ID,
				Name:       user.Name,
				Email:      user.Email,
				Role:       user.Role,
				DisplayTag: user.DisplayTag,
				Colour:     user.Colour,
			},
			FileCount: fileCount,
		}
	}

	JSON(w, http.StatusOK, response)
}

// ListUserFiles returns files for a specific user
func (h *UserHandler) ListUserFiles(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID64, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		ErrorMessage(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	userID := int32(userID64)

	limit := int32(min(1000, maxListLimit))
	files, err := h.fileSvc.ListFilesByUserID(r.Context(), userID, limit, 0)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	response := make([]FileResponse, len(files))
	for i, f := range files {
		response[i] = toFileResponse(f)
	}

	JSON(w, http.StatusOK, response)
}
