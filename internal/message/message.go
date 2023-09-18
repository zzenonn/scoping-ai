package messages

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	tnacommon "github.com/zzenonn/trainocate-tna/pkg/common"
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

var (
	ErrNotImplemented = errors.New("this function is not yet implemented")
)

type Answer struct {
	Question tnacommon.Question `json:"question,omitempty" firestore:"question,omitempty"`
	Answer   string             `json:"answer,omitempty" firestore:"answer,omitempty"`
}

// Message representation
type Message struct {
	Id          string     `json:"id" firestore:"id"`
	UserId      *string    `json:"user_id,omitempty" firestore:"user_id,omitempty"`
	MessageText *string    `json:"message_text,omitempty" firestore:"message_text,omitempty"`
	Answer      Answer     `json:"answer,omitempty" firestore:"answer,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty" firestore:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" firestore:"updated_at,omitempty"`
}

// Implements the message repository interface design pattern
type MessageRepository interface {
	GetMessage(ctx context.Context, messageId string, userId string) (Message, error)
	GetAllUserMessages(ctx context.Context, userId string, page int, pageSize int) ([]Message, error)
	PostMessage(ctx context.Context, message Message) (Message, error)
	UpdateMessage(ctx context.Context, message Message) (Message, error)
	DeleteMessage(ctx context.Context, messageId string, userId string) error
}

type MessageService struct {
	messageRepository MessageRepository
}

func NewMessageService(messageRepository MessageRepository) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
	}
}

func (service *MessageService) PostMessage(ctx context.Context, message Message) (Message, error) {
	log.Debug("Posting message . . .")

	message.Id = uuid.New().String()

	postedMessage, err := service.messageRepository.PostMessage(ctx, message)

	if err != nil {
		log.Error("Failed to post message")
		return Message{}, err
	}

	return postedMessage, nil
}
func (service *MessageService) GetMessage(ctx context.Context, messageId string, userId string) (Message, error) {
	log.Debugf("Retreiving message Id: %s for user %s . . .", messageId, userId)

	message, err := service.messageRepository.GetMessage(ctx, messageId, userId)

	if err != nil {
		log.Errorf("Failed to retrieve message %s for user", messageId, userId)
		return Message{}, err
	}

	return message, nil
}
func (service *MessageService) GetAllUserMessages(ctx context.Context, userId string, page int, pageSize int) ([]Message, error) {
	log.Debug("Retreiving all course messages . . .")

	messages, err := service.messageRepository.GetAllUserMessages(ctx, userId, page, pageSize)

	if err != nil {
		log.Error("Failed to retrieve all messages")
		return nil, err
	}

	return messages, nil
}
func (service *MessageService) UpdateMessage(ctx context.Context, message Message) (Message, error) {
	log.Debugf("Updating message %s", message.Id)

	updatedMessage, err := service.messageRepository.UpdateMessage(ctx, message)

	if err != nil {
		log.Errorf("Failed to update message %s", message.Id)
		return Message{}, err
	}

	return updatedMessage, nil
}

func (service *MessageService) DeleteMessage(ctx context.Context, messageId string, userId string) error {
	log.Debugf("Deleting message %s from user %s. . .", messageId, userId)

	err := service.messageRepository.DeleteMessage(ctx, messageId, userId)

	if err != nil {
		log.Errorf("Failed to delete message %s from user %s", messageId, userId)
		return err
	}

	return nil
}
