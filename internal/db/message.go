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

const MESSAGE_COLLECTION_NAME = "messages"

type MessageRepository struct {
	client *firestore.Client
}

func NewMessageRepository(client *firestore.Client) MessageRepository {
	return MessageRepository{
		client: client,
	}
}

func convertMessageToMap(message tnamessage.Message) (map[string]interface{}, error) {
	if message.UserId == nil || (message.MessageText == nil && message.Answer.Question.Category == nil) {
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

	if message.Answer.Question.Category != nil {
		answerMap := map[string]interface{}{
			"question": map[string]interface{}{
				"category": message.Answer.Question.Category,
				"text":     message.Answer.Question.Text,
				"options": map[string]interface{}{
					"multi_answer":     message.Answer.Question.Options.MultiAnswer,
					"possible_options": message.Answer.Question.Options.PossibleOptions,
				},
			},
			"answer": message.Answer.Answer,
		}
		messageMap["answers"] = answerMap
	}

	if message.CreatedAt != nil {
		messageMap["created_at"] = message.CreatedAt.Format(time.RFC3339) // Will be overwritted by firestore.ServerTimestamp
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

	_, err = repo.client.Collection(USER_COLLECTION_NAME).Doc(*userId).Collection(MESSAGE_COLLECTION_NAME).Doc(message.Id).Set(ctx, messageMap)
	if err != nil {
		return tnamessage.Message{}, err
	}

	return message, nil
}

func (repo *MessageRepository) GetMessage(ctx context.Context, messageId string, userId string) (tnamessage.Message, error) {

	doc, err := repo.client.Collection(USER_COLLECTION_NAME).Doc(userId).Collection(MESSAGE_COLLECTION_NAME).Doc(messageId).Get(ctx)
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

	iter := repo.client.Collection(USER_COLLECTION_NAME).Doc(userId).Collection(MESSAGE_COLLECTION_NAME).OrderBy("created_at", firestore.Asc).Offset(offset).Limit(pageSize).Documents(ctx)
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

	_, err = repo.client.Collection(MESSAGE_COLLECTION_NAME).Doc(message.Id).Set(ctx, messageMap, firestore.MergeAll)
	if err != nil {
		return tnamessage.Message{}, err
	}

	return message, nil
}

func (repo *MessageRepository) DeleteMessage(ctx context.Context, messageId string, userId string) error {
	_, err := repo.client.Collection(USER_COLLECTION_NAME).Doc(userId).Collection(MESSAGE_COLLECTION_NAME).Doc(messageId).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
