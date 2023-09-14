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

type QuestionSetRepository struct {
	client *firestore.Client
}

func NewQuestionSetRepository(client *firestore.Client) *QuestionSetRepository {
	return &QuestionSetRepository{
		client: client,
	}
}

func PostQuestionSet(ctx context.Context, client *firestore.Client, qSet TrainingNeedsQuestions.QuestionSet) (string, error) {
	ref, _, err := client.Collection("questionSets").Add(ctx, qSet)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

func GetQuestionSet(ctx context.Context, client *firestore.Client, docID string) (*TrainingNeedsQuestions.QuestionSet, error) {
	doc, err := client.Collection("questionSets").Doc(docID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var qSet TrainingNeedsQuestions.QuestionSet
	doc.DataTo(&qSet)
	return &qSet, nil
}

func GetQuestionSetByTechName(ctx context.Context, client *firestore.Client, techName string) (*TrainingNeedsQuestions.QuestionSet, error) {
	iter := client.Collection("questionSets").Where("TechnologyName", "==", techName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var qSet TrainingNeedsQuestions.QuestionSet
		err = doc.DataTo(&qSet)
		if err != nil {
			return nil, err
		}
		return &qSet, nil
	}
	return nil, nil
}

func GetAllQuestionSets(ctx context.Context, client *firestore.Client) ([]TrainingNeedsQuestions.QuestionSet, error) {
	iter := client.Collection("questionSets").Documents(ctx)
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

func UpdateQuestionSet(ctx context.Context, client *firestore.Client, docID string, updates []firestore.Update) error {
	_, err := client.Collection("questionSets").Doc(docID).Update(ctx, updates)
	return err
}
