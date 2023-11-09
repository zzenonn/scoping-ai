package messages

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	scopingaicommon "github.com/zzenonn/scoping-ai/pkg/common"
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
	Question       *scopingaicommon.Question `json:"question,omitempty" firestore:"question,omitempty"`
	TechnologyName *string                   `json:"technology_name,omitempty" firestore:"technology_name,omitempty"`
	Answer         *string                   `json:"answer,omitempty" firestore:"answer,omitempty"`
}

// Message representation
type Message struct {
	Id          string     `json:"id" firestore:"id"`
	UserId      *string    `json:"user_id,omitempty" firestore:"-"`
	MessageText *string    `json:"message_text,omitempty" firestore:"message_text,omitempty"`
	Answer      *Answer    `json:"answer,omitempty" firestore:"answer,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty" firestore:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" firestore:"updated_at,omitempty"`
}

type ChatCompletion struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int           `json:"index"`
	Message      OpenAiMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

type OpenAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Implements the message repository interface design pattern
type MessageRepository interface {
	GetMessage(ctx context.Context, messageId string, userId string) (Message, error)
	GetAllUserMessages(ctx context.Context, userId string, page int, pageSize int) ([]Message, error)
	PostMessage(ctx context.Context, message Message) (Message, error)
	UpdateMessage(ctx context.Context, message Message) (Message, error)
	DeleteMessage(ctx context.Context, messageId string, userId string) error
}

type OpenAiRepository interface {
	PostPrompt(ctx context.Context, aiContext string, prompt string) (ChatCompletion, error)
}

type MessageService struct {
	messageRepository MessageRepository
	openAiRepository  OpenAiRepository
}

func NewMessageService(messageRepository MessageRepository, openAiRepository OpenAiRepository) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
		openAiRepository:  openAiRepository,
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

func (service *MessageService) promptOpenAi(postedMessages []Message, responseMessageId string) (Message, error) {
	log.Debug("Prompting the Open AI API . . .")

	var promptBuilder strings.Builder

	aiContext := `As a sales representative, you'll offer a range of mobile plans:

	Plans Without Device:
	
		Plan 599: Provides 2.5GB of data, unlimited calls and texts across all networks, and 25 minutes of calls to PLDT.
		Plan 599 with Netflix: Same as above but includes 5GB of data, unlimited landline calls, and a Netflix subscription.
		Plan 999: Offers 5GB data, unlimited calls and texts, and 50 minutes to PLDT.
		Plan 1,499: Upgrades to 10GB data, with unlimited calls and texts, plus 100 minutes to PLDT.
		Plan 1,999: Includes 15GB data, unlimited calls and texts, and 150 minutes to PLDT.
		Plan 2,499: Offers a large data allowance of 60GB, with unlimited calls, texts, and a Netflix subscription.
		Plan 2,999: Comes with 30GB data, unlimited calls and texts, and the highest PLDT call allowance of 300 minutes.
	
	Plans With Device:
	
		Plan 599: Includes 2.5GB data, unlimited calls and texts, 25 minutes to PLDT, with a device.
		Plan 999: Offers 20GB data, unlimited 5G (for 12 months), calls, texts, landline calls, Netflix, and a device.
		Plan 1,499 Option 1: Provides 40GB data, unlimited 5G access (12-month promo), calls, texts, landline calls, Netflix, and a device.
		Plan 1,999: Comes with 60GB data, unlimited 5G (for 12 months), calls, texts, landline calls, Netflix, and a device.
		Plan 2,999: Includes unlimited 4G/5G data, calls, texts, landline calls, Netflix, and a device.
	
	Remember that the "Unlimited 5G for 12 months" is a promotional feature and subject to change post-promotion. Your role is to understand customer needs and align them with the most suitable plan.
	Do not ask for any additional feedback or elaboration from the customer.`

	for _, msg := range postedMessages {
		if msg.Answer != nil && msg.Answer.Question != nil && msg.Answer.Answer != nil {
			promptBuilder.WriteString("Question: ")
			promptBuilder.WriteString(*msg.Answer.Question.Text)
			promptBuilder.WriteString("\nAnswer: ")
			promptBuilder.WriteString(*msg.Answer.Answer)
			promptBuilder.WriteString("\n\n")
		} else {
			log.Error("Message with ID ", msg.Id, " lacks either a question or an answer or both.")
		}
	}

	prompt := promptBuilder.String()

	chatCompletion, err := service.openAiRepository.PostPrompt(context.Background(), aiContext, prompt)

	if err != nil {
		log.Error("Failed to prompt Open AI API")
		return Message{}, err
	}

	var message Message

	message.Id = responseMessageId
	message.UserId = postedMessages[0].UserId

	jsonData, err := json.Marshal(chatCompletion)
	if err != nil {
		log.Error("Error marshaling struct")
		return Message{}, err
	}

	jsonString := string(jsonData)

	message.MessageText = &jsonString

	completionMessage, err := service.UpdateMessage(context.Background(), message)

	return completionMessage, nil
}

func (service *MessageService) PostAnswers(ctx context.Context, messages []Message) (Message, error) {
	log.Debug("Posting multiple answers...")

	postedMessages := make([]Message, 0, len(messages))

	for _, message := range messages {
		message.Id = uuid.New().String()
		postedMessage, err := service.PostMessage(ctx, message)
		if err != nil {
			log.Errorf("Failed to post message with ID %s. Error: %v", message.Id, err)

			continue
		}
		postedMessages = append(postedMessages, postedMessage)
	}

	messagePending := "Thank you for your message. Please wait for the AI Engine to generate a response."

	pendingMessage := Message{
		Id:          uuid.New().String(),
		UserId:      postedMessages[0].UserId,
		MessageText: &messagePending,
	}

	postedPendingMessage, err := service.PostMessage(ctx, pendingMessage)

	if err != nil {
		log.Errorf("Failed to post pending message with ID %s. Error: %v", pendingMessage.Id, err)
		return Message{}, err
	}

	defer func() {
		go service.promptOpenAi(postedMessages, postedPendingMessage.Id)
	}()

	log.Debug("Completed posting messages.")
	return postedPendingMessage, nil
}

func (service *MessageService) GetMessage(ctx context.Context, messageId string, userId string) (Message, error) {
	log.Debugf("Retreiving message Id: %s for user %s . . .", messageId, userId)

	message, err := service.messageRepository.GetMessage(ctx, messageId, userId)

	if err != nil {
		log.Errorf("Failed to retrieve message %s for user %s", messageId, userId)
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
