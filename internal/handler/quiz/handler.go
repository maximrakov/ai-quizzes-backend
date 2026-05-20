package quiz

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	roothandler "github.com/maximrakov/ai-quizzes-backend/internal/handler"
	"github.com/maximrakov/ai-quizzes-backend/internal/handler/quiz/dto"
	"github.com/maximrakov/ai-quizzes-backend/internal/model"
	pkgjwt "github.com/maximrakov/ai-quizzes-backend/pkg/jwt"
)

type Service interface {
	Create(ctx context.Context, title string, creatorId int, questions []model.Question) (*model.Quiz, error)
	Assign(ctx context.Context, quizId, studentId, mentorId int) error
	FindAll(ctx context.Context) ([]*model.Quiz, error)
	FindById(ctx context.Context, id int) (*model.Quiz, error)
	FindByCreatorId(ctx context.Context, creatorId int) ([]*model.Quiz, error)
	FindByUserId(ctx context.Context, userId int) ([]*model.Quiz, error)
	GenerateQuestions(ctx context.Context, topic string, count int) ([]model.Question, error)
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func claimsFromCtx(r *http.Request) *pkgjwt.Claims {
	return r.Context().Value(roothandler.UserClaimsKey).(*pkgjwt.Claims)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateQuizRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный json", http.StatusBadRequest)
		return
	}

	creatorId := claimsFromCtx(r).UserID
	questions := dto.ToQuestions(input.Questions)
	quiz, err := h.service.Create(context.Background(), input.Title, creatorId, questions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(dto.ToQuizResponse(quiz)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) Assign(w http.ResponseWriter, r *http.Request) {
	quizIdStr := r.PathValue("quizId")
	quizId, err := strconv.Atoi(quizIdStr)
	if err != nil {
		http.Error(w, "некорректный id квиза", http.StatusBadRequest)
		return
	}

	var input dto.AssignQuizRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный json", http.StatusBadRequest)
		return
	}

	if err = h.service.Assign(context.Background(), quizId, input.StudentId, input.MentorId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) FindAll(w http.ResponseWriter, r *http.Request) {
	quizzes, err := h.service.FindAll(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dto.ToQuizResponses(quizzes)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) FindById(w http.ResponseWriter, r *http.Request) {
	quizIdStr := r.PathValue("quizId")
	quizId, err := strconv.Atoi(quizIdStr)
	if err != nil {
		http.Error(w, "некорректный id квиза", http.StatusBadRequest)
		return
	}

	quiz, err := h.service.FindById(context.Background(), quizId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dto.ToQuizResponse(quiz)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) FindCreated(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "некорректный id пользователя", http.StatusBadRequest)
		return
	}

	quizzes, err := h.service.FindByCreatorId(context.Background(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dto.ToQuizResponses(quizzes)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) FindAssigned(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "некорректный id пользователя", http.StatusBadRequest)
		return
	}

	quizzes, err := h.service.FindByUserId(context.Background(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dto.ToQuizResponses(quizzes)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) Generate(w http.ResponseWriter, r *http.Request) {
	var input dto.GenerateQuestionsRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный json", http.StatusBadRequest)
		return
	}

	if input.Topic == "" {
		http.Error(w, "topic обязателен", http.StatusBadRequest)
		return
	}
	if input.Title == "" {
		http.Error(w, "title обязателен", http.StatusBadRequest)
		return
	}
	if input.Count <= 0 {
		input.Count = 5
	}

	questions, err := h.service.GenerateQuestions(r.Context(), input.Topic, input.Count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	creatorId := claimsFromCtx(r).UserID
	quiz, err := h.service.Create(r.Context(), input.Title, creatorId, questions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(dto.ToQuizResponse(quiz)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) FindByUserId(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "некорректный id пользователя", http.StatusBadRequest)
		return
	}

	quizzes, err := h.service.FindByUserId(context.Background(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dto.ToQuizResponses(quizzes)); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}
