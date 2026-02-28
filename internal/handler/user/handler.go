package user

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/maximrakov/ai-quizzes-backend/internal/handler/user/dto"
	"github.com/maximrakov/ai-quizzes-backend/internal/model"
	userS "github.com/maximrakov/ai-quizzes-backend/internal/service/user"
)

type Service interface {
	Create(ctx context.Context, nickname, password string, role model.Role) (*model.User, error)
	Login(ctx context.Context, nickname, password string) (string, error)
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный json", http.StatusBadRequest)
		return
	}

	user, err := h.service.Create(context.Background(), input.Username, input.Password, input.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(dto.ToUserResponse(user)); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный json", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(context.Background(), input.Username, input.Password)
	if err != nil {
		if errors.Is(err, userS.ErrInvalidCredentials) {
			http.Error(w, "неверный логин или пароль", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dto.LoginResponse{Token: token}); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
