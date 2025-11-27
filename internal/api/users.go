// user api for setting up user accounts and authentication
package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"cliscraper/internal/database"
)

type UserHandler struct {
	userRepo *database.UserRepository
}

func NewUserHandler(userRepo *database.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func writeJSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

// request body for user registration
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// request body for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// POST /api/register
// body: {"username": "user123", "email": "user@example.com", "password": "securepass"}
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Username, email, and password are required",
		})
		return
	}

	if len(req.Password) < 6 {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Password must be at least 6 characters",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to hash password",
		})
		return
	}

	user, err := h.userRepo.CreateUser(req.Username, req.Email, string(hashedPassword))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			status = http.StatusConflict
		}
		writeJSON(w, status, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	userData := map[string]interface{}{
		"id":        user.ID.Hex(),
		"username":  user.Username,
		"email":     user.Email,
		"created_at": user.CreatedAt,
	}
	dataJSON, _ := json.Marshal(userData)
	writeJSON(w, http.StatusCreated, Response{
		Status: "ok",
		Data:   dataJSON,
	})
}

// POST /api/login
// body: {"username": "user123", "password": "securepass"}
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Username and password are required",
		})
		return
	}

	user, err := h.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, Response{
			Status:  "error",
			Message: "Invalid username or password",
		})
		return
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, Response{
			Status:  "error",
			Message: "Invalid username or password",
		})
		return
	}

	userData := map[string]interface{}{
		"id":        user.ID.Hex(),
		"username":  user.Username,
		"email":     user.Email,
		"created_at": user.CreatedAt,
	}
	dataJSON, _ := json.Marshal(userData)
	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data:   dataJSON,
	})
}
