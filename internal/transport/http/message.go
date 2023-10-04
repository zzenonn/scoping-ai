package http

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	tnamessage "github.com/zzenonn/trainocate-tna/internal/message"
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

type MessageServiceInterface interface {
	PostMessage(ctx context.Context, message tnamessage.Message) (tnamessage.Message, error)
	PostAnswers(ctx context.Context, messages []tnamessage.Message) (tnamessage.Message, error)
	GetMessage(ctx context.Context, messageId string, userId string) (tnamessage.Message, error)
	GetAllUserMessages(ctx context.Context, userId string, page int, pageSize int) ([]tnamessage.Message, error)
	UpdateMessage(ctx context.Context, message tnamessage.Message) (tnamessage.Message, error)
	DeleteMessage(ctx context.Context, messageId string, userId string) error
}

type MessageHandler struct {
	messageService MessageServiceInterface
}

func NewMessageHandler(s MessageServiceInterface) *MessageHandler {
	return &MessageHandler{
		messageService: s,
	}
}

func (h *MessageHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	var message tnamessage.Message

	message.UserId = &userId

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message, err := h.messageService.PostMessage(r.Context(), message)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *MessageHandler) PostAnswers(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	// Change to an array of messages
	var messages []tnamessage.Message

	// Set userId for all messages
	for i := range messages {
		messages[i].UserId = &userId
	}

	// Decode the request body into the messages slice
	if err := json.NewDecoder(r.Body).Decode(&messages); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseMessage, err := h.messageService.PostAnswers(r.Context(), messages)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode and return the processed messages
	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	messageId := chi.URLParam(r, "messageId")

	if userId == "" || messageId == "" {
		http.Error(w, "User ID and Message ID are required", http.StatusBadRequest)
		return
	}

	message, err := h.messageService.GetMessage(r.Context(), messageId, userId)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *MessageHandler) GetAllUserMessages(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get page and pageSize from query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Convert them to integers with some default values
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	messages, err := h.messageService.GetAllUserMessages(r.Context(), userId, page, pageSize)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	var message tnamessage.Message

	userId := chi.URLParam(r, "userId")
	messageId := chi.URLParam(r, "messageId")

	if userId == "" || messageId == "" {
		http.Error(w, "User ID and Message ID are required", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Setting userId and messageId from the URL parameters
	message.UserId = &userId
	message.Id = messageId

	updatedMessage, err := h.messageService.UpdateMessage(r.Context(), message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(updatedMessage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	messageId := chi.URLParam(r, "messageId")

	if userId == "" || messageId == "" {
		http.Error(w, "User ID and Message ID are required", http.StatusBadRequest)
		return
	}

	if err := h.messageService.DeleteMessage(r.Context(), messageId, userId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandler) mapRoutes(router chi.Router) {
	router.Route("/api/v1/users/{userId}/messages", func(r chi.Router) {
		r.Post("/", h.PostMessage)
		r.Post("/answers", h.PostMessage)
		r.Get("/", h.GetAllUserMessages)

		r.Route("/{messageId}", func(r chi.Router) {
			r.Get("/", h.GetMessage)
			r.Put("/", h.UpdateMessage)
			r.Delete("/", h.DeleteMessage)
		})
	})
}
