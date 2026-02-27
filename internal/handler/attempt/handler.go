package attempt

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/maximrakov/ai-quizzes-backend/internal/handler/attempt/dto"
	"github.com/maximrakov/ai-quizzes-backend/internal/model"
	attemptS "github.com/maximrakov/ai-quizzes-backend/internal/service/attempt"
)

type Service interface {
	Create(ctx context.Context, assignmentId int, answers []attemptS.UserAnswer) (*model.Attempt, error)
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный json", http.StatusBadRequest)
		return
	}

	answers := make([]attemptS.UserAnswer, 0, len(input.Answers))
	for _, a := range input.Answers {
		answers = append(answers, attemptS.UserAnswer{
			QuestionId:   a.QuestionId,
			AnswerNumber: a.AnswerNumber,
		})
	}

	attempt, err := h.service.Create(context.Background(), input.AssignmentId, answers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(dto.ToAttemptResponse(attempt)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}
