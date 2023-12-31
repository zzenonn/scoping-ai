package db

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
	questionSet "github.com/zzenonn/scoping-ai/internal/question-set"
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
	client         *firestore.Client
	CollectionName string
}

func NewQuestionSetRepository(client *firestore.Client, collectionName string) QuestionSetRepository {
	return QuestionSetRepository{
		client:         client,
		CollectionName: collectionName,
	}
}

func convertQuestionSetToMap(qSet questionSet.QuestionSet) map[string]interface{} {
	qSetMap := make(map[string]interface{})

	// Omitting the ID field to avoid redundant Ids
	if qSet.TechnologyName != nil {
		qSetMap["technology_name"] = *qSet.TechnologyName
	}

	questions := make([]map[string]interface{}, len(qSet.Questions))
	for i, question := range qSet.Questions {
		questionMap := make(map[string]interface{})
		if question.Category != nil {
			questionMap["category"] = *question.Category
		}
		if question.Text != nil {
			questionMap["text"] = *question.Text
		}
		if question.Options != nil {
			optionsMap := make(map[string]interface{})
			optionsMap["multi_answer"] = question.Options.MultiAnswer
			optionsMap["possible_options"] = question.Options.PossibleOptions
			questionMap["options"] = optionsMap
		}
		questions[i] = questionMap
	}

	qSetMap["questions"] = questions

	return qSetMap
}

func (repo *QuestionSetRepository) PostQuestionSet(ctx context.Context, qSet questionSet.QuestionSet) (questionSet.QuestionSet, error) {

	qSetMap := convertQuestionSetToMap(qSet)

	_, err := repo.client.Collection(repo.CollectionName).Doc(qSet.Id).Set(ctx, qSetMap)
	if err != nil {
		return questionSet.QuestionSet{}, err
	}

	return qSet, nil
}

func (repo *QuestionSetRepository) GetQuestionSet(ctx context.Context, docID string) (questionSet.QuestionSet, error) {
	doc, err := repo.client.Collection(repo.CollectionName).Doc(docID).Get(ctx)
	if err != nil {
		return questionSet.QuestionSet{}, err
	}

	var qSet questionSet.QuestionSet
	doc.DataTo(&qSet)
	qSet.Id = doc.Ref.ID

	return qSet, nil
}

func (repo *QuestionSetRepository) GetQuestionSetByTechName(ctx context.Context, techName string) (questionSet.QuestionSet, error) {
	iter := repo.client.Collection(repo.CollectionName).Where("technology_name", "==", techName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return questionSet.QuestionSet{}, err
		}

		var qSet questionSet.QuestionSet
		err = doc.DataTo(&qSet)
		qSet.Id = doc.Ref.ID

		if err != nil {
			return questionSet.QuestionSet{}, err
		}
		return qSet, nil
	}
	return questionSet.QuestionSet{}, nil
}

func (repo *QuestionSetRepository) GetAllQuestionSets(ctx context.Context, page int, pageSize int) ([]questionSet.QuestionSet, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	iter := repo.client.Collection(repo.CollectionName).OrderBy("technology_name", firestore.Asc).Offset(offset).Limit(pageSize).Documents(ctx)
	var qSets []questionSet.QuestionSet

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var qSet questionSet.QuestionSet
		err = doc.DataTo(&qSet)
		qSet.Id = doc.Ref.ID

		if err != nil {
			return nil, err
		}

		qSets = append(qSets, qSet)
	}
	return qSets, nil
}

func (repo *QuestionSetRepository) UpdateQuestionSet(ctx context.Context, qSet questionSet.QuestionSet) (questionSet.QuestionSet, error) {
	qSetMap := convertQuestionSetToMap(qSet)
	log.Debugf("Updating question set: %v", qSet.Id)
	_, err := repo.client.Collection(repo.CollectionName).Doc(qSet.Id).Set(ctx, qSetMap, firestore.MergeAll)
	return qSet, err
}

func (repo *QuestionSetRepository) DeleteQuestionSet(ctx context.Context, docID string) error {
	_, err := repo.client.Collection(repo.CollectionName).Doc(docID).Delete(ctx)
	return err
}
