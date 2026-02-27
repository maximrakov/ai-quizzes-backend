package assignment

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/maximrakov/ai-quizzes-backend/internal/handler/assignment/dto"
	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

type Service interface {
	Create(ctx context.Context, quizId, studentId, mentorId int) (*model.Assignment, error)
	FindById(ctx context.Context, id int) (*model.Assignment, error)
	FindByStudentId(ctx context.Context, studentId int) ([]*model.Assignment, error)
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный json", http.StatusBadRequest)
		return
	}

	assignment, err := h.service.Create(context.Background(), input.QuizId, input.StudentId, input.MentorId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(dto.ToAssignmentResponse(assignment)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) FindById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "некорректный id назначения", http.StatusBadRequest)
		return
	}

	assignment, err := h.service.FindById(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dto.ToAssignmentResponse(assignment)); err != nil {
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

	assignments, err := h.service.FindByStudentId(context.Background(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dto.ToAssignmentResponses(assignments)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}
