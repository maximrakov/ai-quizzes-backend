package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/maximrakov/ai-quizzes-backend/internal/app"
	"github.com/maximrakov/ai-quizzes-backend/internal/handler/user/dto"
	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

type Service interface {
	Create(ctx context.Context, nickname, password string, role model.Role) (*model.User, error)
}
type handler struct {
	service Service
	ctx     app.Context
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный json", http.StatusBadRequest)
	}

	user, err := h.service.Create(context.Background(), input.Username, input.Password, input.Role)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	//h.ctx.Logger.Info("aboba ", user.Id, " ", user.Nickname, " ", user.Password, " ", user.Role)

	if err = json.NewEncoder(w).Encode(dto.ToUserResponse(user)); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
