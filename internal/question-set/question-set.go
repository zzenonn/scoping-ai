package TrainingNeedsQuestions

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	scopingaicommon "github.com/zzenonn/scoping-ai/pkg/common"
)

var (
	ErrFetchingQuestions = errors.New("failed to fetch scoping questions by technology")
	ErrNotImplemented    = errors.New("this function is not yet implemented")
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

// Question set representation
type QuestionSet struct {
	Id             string                     `json:"id,omitempty" firestore:"id,omitempty"`
	TechnologyName *string                    `json:"technology_name,omitempty" firestore:"technology_name,omitempty"`
	Questions      []scopingaicommon.Question `json:"questions,omitempty" firestore:"questions,omitempty"`
}

// Implements the question set repository interface design pattern
type QuestionSetRepository interface {
	GetQuestionSet(ctx context.Context, technologyName string) (QuestionSet, error)
	GetAllQuestionSets(ctx context.Context, page int, pageSize int) ([]QuestionSet, error)
	GetQuestionSetByTechName(ctx context.Context, technologyName string) (QuestionSet, error)
	PostQuestionSet(ctx context.Context, questionSet QuestionSet) (QuestionSet, error)
	UpdateQuestionSet(ctx context.Context, questionSet QuestionSet) (QuestionSet, error)
	DeleteQuestionSet(ctx context.Context, id string) error
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

	return questionSet, nil
}

func (q *QuestionSetService) GetAllQuestionSets(ctx context.Context, page int, pageSize int) ([]QuestionSet, error) {
	log.Debug("Retreiving all question sets . . .")

	questionSets, err := q.questionSetRepository.GetAllQuestionSets(ctx, page, pageSize)

	if err != nil {
		log.Error("Failed to retrieve all question sets")
		return nil, err
	}

	return questionSets, nil
}

func (q *QuestionSetService) GetQuestionSetByTechName(ctx context.Context, technologyName string) (QuestionSet, error) {
	log.Debug("Retreiving question set by technology name . . .")

	qSet, err := q.questionSetRepository.GetQuestionSetByTechName(ctx, technologyName)

	if err != nil {
		log.Error("Failed to retrieve question set by technology name")
		return QuestionSet{}, ErrFetchingQuestions
	}

	return qSet, nil
}

func (q *QuestionSetService) PostQuestionSet(ctx context.Context, qSet QuestionSet) (QuestionSet, error) {
	log.Debug("Posting question set . . .")
	qSet.Id = uuid.New().String()

	postedQSet, err := q.questionSetRepository.PostQuestionSet(ctx, qSet)

	if err != nil {
		log.Error("Failed to post question set")
		return QuestionSet{}, err
	}

	return postedQSet, nil
}

func (q *QuestionSetService) UpdateQuestionSet(ctx context.Context, qSet QuestionSet) (QuestionSet, error) {
	log.Debug("Updating question set . . .")

	updatedQSet, err := q.questionSetRepository.UpdateQuestionSet(ctx, qSet)

	if err != nil {
		log.Error("Failed to update question set")
		return QuestionSet{}, err
	}

	return updatedQSet, nil
}

func (q *QuestionSetService) DeleteQuestionSet(ctx context.Context, id string) error {
	err := q.questionSetRepository.DeleteQuestionSet(ctx, id)

	return err
}
