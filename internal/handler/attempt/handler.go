package attempt

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/maximrakov/ai-quizzes-backend/internal/handler/attempt/dto"
	"github.com/maximrakov/ai-quizzes-backend/internal/model"
	attemptS "github.com/maximrakov/ai-quizzes-backend/internal/service/attempt"
)

type Service interface {
	Create(ctx context.Context, assignmentId int, answers []attemptS.UserAnswer) (*model.Attempt, error)
	GetByStudentId(ctx context.Context, studentId int) ([]*model.Attempt, error)
	GetByQuizId(ctx context.Context, quizId int) ([]*model.Attempt, error)
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

func (h *handler) FindByQuizId(w http.ResponseWriter, r *http.Request) {
	quizIdStr := r.PathValue("quizId")
	quizId, err := strconv.Atoi(quizIdStr)
	if err != nil {
		http.Error(w, "некорректный id квиза", http.StatusBadRequest)
		return
	}

	attempts, err := h.service.GetByQuizId(context.Background(), quizId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]dto.AttemptResponse, 0, len(attempts))
	for _, a := range attempts {
		result = append(result, dto.ToAttemptResponse(a))
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) FindByStudentId(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "некорректный id пользователя", http.StatusBadRequest)
		return
	}

	attempts, err := h.service.GetByStudentId(context.Background(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]dto.AttemptResponse, 0, len(attempts))
	for _, a := range attempts {
		result = append(result, dto.ToAttemptResponse(a))
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}
