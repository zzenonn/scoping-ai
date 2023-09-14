package db

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
	"github.com/zzenonn/trainocate-tna/internal/TrainingNeedsQuestions"
	"google.golang.org/api/iterator"
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

const COLLECTION_NAME = "question_sets"

type QuestionSetRepository struct {
	client *firestore.Client
}

func NewQuestionSetRepository(client *firestore.Client) QuestionSetRepository {
	return QuestionSetRepository{
		client: client,
	}
}

func (repo *QuestionSetRepository) PostQuestionSet(ctx context.Context, qSet TrainingNeedsQuestions.QuestionSet) (TrainingNeedsQuestions.QuestionSet, error) {
	ref, _, err := repo.client.Collection(COLLECTION_NAME).Add(ctx, qSet)
	if err != nil {
		return TrainingNeedsQuestions.QuestionSet{}, err
	}

	postedQSet, err := repo.GetQuestionSet(ctx, ref.ID)
	if err != nil {
		return TrainingNeedsQuestions.QuestionSet{}, err
	}

	return postedQSet, nil
}

func (repo *QuestionSetRepository) GetQuestionSet(ctx context.Context, docID string) (TrainingNeedsQuestions.QuestionSet, error) {
	doc, err := repo.client.Collection(COLLECTION_NAME).Doc(docID).Get(ctx)
	if err != nil {
		return TrainingNeedsQuestions.QuestionSet{}, err
	}

	var qSet TrainingNeedsQuestions.QuestionSet
	doc.DataTo(&qSet)
	return qSet, nil
}

func (repo *QuestionSetRepository) GetQuestionSetByTechName(ctx context.Context, techName string) (TrainingNeedsQuestions.QuestionSet, error) {
	iter := repo.client.Collection(COLLECTION_NAME).Where("TechnologyName", "==", techName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return TrainingNeedsQuestions.QuestionSet{}, err
		}

		var qSet TrainingNeedsQuestions.QuestionSet
		err = doc.DataTo(&qSet)
		if err != nil {
			return TrainingNeedsQuestions.QuestionSet{}, err
		}
		return qSet, nil
	}
	return TrainingNeedsQuestions.QuestionSet{}, nil
}

func (repo *QuestionSetRepository) GetAllQuestionSets(ctx context.Context, page int, pageSize int) ([]TrainingNeedsQuestions.QuestionSet, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	iter := repo.client.Collection(COLLECTION_NAME).OrderBy("TechnologyName", firestore.Asc).Offset(offset).Limit(pageSize).Documents(ctx)
	var qSets []TrainingNeedsQuestions.QuestionSet

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var qs TrainingNeedsQuestions.QuestionSet
		err = doc.DataTo(&qs)
		if err != nil {
			return nil, err
		}

		qSets = append(qSets, qs)
	}
	return qSets, nil
}

func (repo *QuestionSetRepository) UpdateQuestionSet(ctx context.Context, qSet TrainingNeedsQuestions.QuestionSet) (TrainingNeedsQuestions.QuestionSet, error) {
	_, err := repo.client.Collection(COLLECTION_NAME).Doc(qSet.Id).Set(ctx, qSet, firestore.MergeAll)
	return qSet, err
}

func (repo *QuestionSetRepository) DeleteQuestionSet(ctx context.Context, docID string) error {
	_, err := repo.client.Collection(COLLECTION_NAME).Doc(docID).Delete(ctx)
	return err
}
