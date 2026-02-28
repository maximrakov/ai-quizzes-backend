package quiz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

var (
	ErrNotMentor  = errors.New("only mentor can create quizzes")
	ErrNotStudent = errors.New("quiz can only be assigned to a student")
)

type AIClient interface {
	Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

// aiQuestion is the JSON structure expected from the AI response.
type aiQuestion struct {
	Text                string     `json:"text"`
	CorrectAnswerNumber int        `json:"correct_answer_number"`
	Options             []aiOption `json:"options"`
}

type aiOption struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

type Repository interface {
	Create(ctx context.Context, quiz *model.Quiz) (*model.Quiz, error)
	FindById(ctx context.Context, id int) (*model.Quiz, error)
	FindAll(ctx context.Context) ([]*model.Quiz, error)
	AssignToUser(ctx context.Context, quizId, userId int) error
	FindByUserId(ctx context.Context, userId int) ([]*model.Quiz, error)
	FindByCreatorId(ctx context.Context, creatorId int) ([]*model.Quiz, error)
}

type UserRepository interface {
	FindById(ctx context.Context, id int) (*model.User, error)
}

type service struct {
	repo     Repository
	userRepo UserRepository
	aiClient AIClient
}

func NewService(repo Repository, userRepo UserRepository, aiClient AIClient) *service {
	return &service{repo: repo, userRepo: userRepo, aiClient: aiClient}
}

func (s *service) Create(ctx context.Context, title string, creatorId int, questions []model.Question) (*model.Quiz, error) {
	creator, err := s.userRepo.FindById(ctx, creatorId)
	if err != nil {
		return nil, err
	}

	if creator.Role != model.Mentor {
		return nil, ErrNotMentor
	}

	quiz := model.NewQuiz(title, creatorId, questions)
	return s.repo.Create(ctx, quiz)
}

func (s *service) Assign(ctx context.Context, quizId, studentId, mentorId int) error {
	mentor, err := s.userRepo.FindById(ctx, mentorId)
	if err != nil {
		return err
	}
	if mentor.Role != model.Mentor {
		return ErrNotMentor
	}

	student, err := s.userRepo.FindById(ctx, studentId)
	if err != nil {
		return err
	}
	if student.Role != model.Student {
		return ErrNotStudent
	}

	return s.repo.AssignToUser(ctx, quizId, studentId)
}

func (s *service) FindAll(ctx context.Context) ([]*model.Quiz, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) FindById(ctx context.Context, id int) (*model.Quiz, error) {
	return s.repo.FindById(ctx, id)
}

func (s *service) FindByCreatorId(ctx context.Context, creatorId int) ([]*model.Quiz, error) {
	return s.repo.FindByCreatorId(ctx, creatorId)
}

func (s *service) FindByUserId(ctx context.Context, userId int) ([]*model.Quiz, error) {
	return s.repo.FindByUserId(ctx, userId)
}

func (s *service) GenerateQuestions(ctx context.Context, topic string, count int) ([]model.Question, error) {
	systemPrompt := `Ты генератор вопросов для викторины.
Отвечай только валидным JSON-массивом — без markdown, без блоков кода, без пояснений.
Все вопросы и варианты ответов должны быть на русском языке.
Каждый элемент массива должен содержать: "text" (string), "correct_answer_number" (int 1-4), "options" (массив из 4 объектов с полями "number" (int 1-4) и "text" (string)).`

	userPrompt := fmt.Sprintf(
		`Сгенерируй %d вопросов по теме "%s". Каждый вопрос должен иметь ровно 4 варианта ответа с номерами от 1 до 4.`,
		count, topic,
	)

	raw, err := s.aiClient.Complete(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("ai request failed: %w", err)
	}

	var aiQuestions []aiQuestion
	if err := json.Unmarshal([]byte(raw), &aiQuestions); err != nil {
		return nil, fmt.Errorf("failed to parse ai response: %w", err)
	}

	questions := make([]model.Question, 0, len(aiQuestions))
	for _, q := range aiQuestions {
		options := make([]model.AnswerOption, 0, len(q.Options))
		for _, o := range q.Options {
			options = append(options, model.AnswerOption{
				Number: o.Number,
				Text:   o.Text,
			})
		}
		questions = append(questions, model.Question{
			Text:                q.Text,
			CorrectAnswerNumber: q.CorrectAnswerNumber,
			Options:             options,
		})
	}

	return questions, nil
}
