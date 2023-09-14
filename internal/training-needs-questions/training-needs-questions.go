package comment

import (
	"context"
	"errors"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	ErrFetchingTna    = errors.New("failed to fetch TNA questions by technology")
	ErrNotImplemented = errors.New("this function is not yet implemented")
)

func init() {

	// Set log level based on environment variables
	switch logLevel := strings.ToLower(os.Getenv("LOG_LEVEL")); logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}

}

// Question representation
type Question struct {
	Category *string `json:"category,omitempty"`
	Text     *string `json:"text,omitempty"`
}

// Question set representation
type QuestionSet struct {
	TechnologyName *string    `json:"technology_name,omitempty"`
	Questions      []Question `json:"questions,omitempty"`
}

// Implements the question set repository interface design pattern
type QuestionSetRepository interface {
	GetQuestionSet(ctx context.Context, technologyName string) (*QuestionSet, error)
	PostQuestionSet(ctx context.Context, questionSet *QuestionSet) (*QuestionSet, error)
	UpdateQuestionSet(ctx context.Context, questionSet *QuestionSet) (*QuestionSet, error)
	DeleteQuestionSet(ctx context.Context, questionSet *QuestionSet) error
}

type QuestionSetService struct {
	questionSetRepository QuestionSetRepository
}

func NewQuestionService(questionSetRepository QuestionSetRepository) *QuestionSetService {
	return &QuestionSetService{
		questionSetRepository: questionSetRepository,
	}
}

func (q *QuestionSetService) GetQuestionSet(ctx context.Context, technologyName string) (QuestionSet, error) {
	log.Debug("Retreiving question set . . .")

	questionSet, err := q.questionSetRepository.GetQuestionSet(ctx, technologyName)

	if err != nil {
		log.Error("Failed to retrieve question set")
		return QuestionSet{}, err
	}

	return *questionSet, nil
}

func (q *QuestionSetService) PostQuestionSet(ctx context.Context, questionSet *QuestionSet) (QuestionSet, error) {
	log.Debug("Posting question set . . .")

	postedQuestionSet, err := q.questionSetRepository.PostQuestionSet(ctx, questionSet)

	if err != nil {
		log.Error("Failed to post question set")
		return QuestionSet{}, err
	}

	return *postedQuestionSet, nil
}

func (q *QuestionSetService) UpdateQuestionSet(ctx context.Context, questionSet *QuestionSet) (QuestionSet, error) {
	log.Debug("Updating question set . . .")

	updatedQuestionSet, err := q.questionSetRepository.UpdateQuestionSet(ctx, questionSet)

	if err != nil {
		log.Error("Failed to update question set")
		return QuestionSet{}, err
	}

	return *updatedQuestionSet, nil
}

func (q *QuestionSetService) DeleteQuestionSet(ctx context.Context, questionSet *QuestionSet) error {
	err := q.questionSetRepository.DeleteQuestionSet(ctx, questionSet)

	return err
}
