package db

import (
	"context"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
	tnamessage "github.com/zzenonn/trainocate-tna/internal/message"
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

// const repo.CollectionName = "messages"

type MessageRepository struct {
	client                *firestore.Client
	MessageCollectionName string
	UserCollectionName    string
}

func NewMessageRepository(client *firestore.Client, messageCollectionName string, userCollectionName string) MessageRepository {
	return MessageRepository{
		client:                client,
		MessageCollectionName: messageCollectionName,
		UserCollectionName:    userCollectionName,
	}
}

func convertMessageToMap(message tnamessage.Message) (map[string]interface{}, error) {
	if message.UserId == nil || (message.MessageText == nil && (message.Answer == nil || message.Answer.Question == nil || (message.Answer.Question.Category == nil && message.Answer.Question.Text == nil) && message.Answer.TechnologyName == nil)) {
		return nil, ErrMissingRequiredFields
	}

	messageMap := make(map[string]interface{})

	messageMap["id"] = message.Id

	if message.UserId != nil {
		messageMap["user_id"] = *message.UserId
	}

	if message.MessageText != nil {
		messageMap["message_text"] = *message.MessageText
	}

	if message.Answer != nil {
		answerMap := map[string]interface{}{}

		if message.Answer.TechnologyName != nil {
			answerMap["technology_name"] = *message.Answer.TechnologyName
		}

		if message.Answer.Question != nil {
			questionMap := map[string]interface{}{}

			if message.Answer.Question.Options != nil {
				questionMap["options"] = map[string]interface{}{
					"multi_answer":     message.Answer.Question.Options.MultiAnswer,
					"possible_options": message.Answer.Question.Options.PossibleOptions,
				}
			}

			if message.Answer.Question.Category != nil {
				questionMap["category"] = *message.Answer.Question.Category
			}

			if message.Answer.Question.Text != nil {
				questionMap["text"] = *message.Answer.Question.Text
			}

			answerMap["question"] = questionMap
		}

		if message.Answer.Answer != nil {
			answerMap["answer"] = *message.Answer.Answer
		}

		messageMap["answer"] = answerMap
	}

	if message.CreatedAt != nil {
		messageMap["created_at"] = message.CreatedAt.Format(time.RFC3339)
	}

	if message.UpdatedAt != nil {
		messageMap["updated_at"] = firestore.ServerTimestamp
	}

	return messageMap, nil
}

func (repo *MessageRepository) PostMessage(ctx context.Context, message tnamessage.Message) (tnamessage.Message, error) {
	messageMap, err := convertMessageToMap(message)

	if err != nil {
		return tnamessage.Message{}, err
	}

	messageMap["created_at"] = firestore.ServerTimestamp

	if err != nil {
		return tnamessage.Message{}, err
	}

	userId := message.UserId

	_, err = repo.client.Collection(repo.UserCollectionName).Doc(*userId).Collection(repo.MessageCollectionName).Doc(message.Id).Set(ctx, messageMap)
	if err != nil {
		return tnamessage.Message{}, err
	}

	return message, nil
}

func (repo *MessageRepository) GetMessage(ctx context.Context, messageId string, userId string) (tnamessage.Message, error) {

	doc, err := repo.client.Collection(repo.UserCollectionName).Doc(userId).Collection(repo.MessageCollectionName).Doc(messageId).Get(ctx)
	if err != nil {
		return tnamessage.Message{}, err
	}

	var message tnamessage.Message
	err = doc.DataTo(&message)
	if err != nil {
		return tnamessage.Message{}, err
	}

	return message, nil
}

func (repo *MessageRepository) GetAllUserMessages(ctx context.Context, userId string, page int, pageSize int) ([]tnamessage.Message, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	iter := repo.client.Collection(repo.UserCollectionName).Doc(userId).Collection(repo.MessageCollectionName).OrderBy("created_at", firestore.Asc).Offset(offset).Limit(pageSize).Documents(ctx)
	var messages []tnamessage.Message

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var message tnamessage.Message
		err = doc.DataTo(&message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (repo *MessageRepository) UpdateMessage(ctx context.Context, message tnamessage.Message) (tnamessage.Message, error) {
	messageMap, err := convertMessageToMap(message)
	if err != nil {
		return tnamessage.Message{}, err
	}

	_, err = repo.client.Collection(repo.UserCollectionName).Doc(*message.UserId).Collection(repo.MessageCollectionName).Doc(message.Id).Set(ctx, messageMap, firestore.MergeAll)
	if err != nil {
		return tnamessage.Message{}, err
	}

	return message, nil
}

func (repo *MessageRepository) DeleteMessage(ctx context.Context, messageId string, userId string) error {
	_, err := repo.client.Collection(repo.UserCollectionName).Doc(userId).Collection(repo.MessageCollectionName).Doc(messageId).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
