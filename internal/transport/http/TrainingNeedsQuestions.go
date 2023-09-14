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
	"github.com/zzenonn/trainocate-tna/internal/TrainingNeedsQuestions"
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

type QuestionSetService interface {
	GetQuestionSet(ctx context.Context, technologyName string) (TrainingNeedsQuestions.QuestionSet, error)
	GetQuestionSetByTechName(ctx context.Context, technologyName string) (TrainingNeedsQuestions.QuestionSet, error)
	GetAllQuestionSets(ctx context.Context, page int, pageSize int) ([]TrainingNeedsQuestions.QuestionSet, error)
	PostQuestionSet(ctx context.Context, questionSet TrainingNeedsQuestions.QuestionSet) (TrainingNeedsQuestions.QuestionSet, error)
	UpdateQuestionSet(ctx context.Context, questionSet TrainingNeedsQuestions.QuestionSet) (TrainingNeedsQuestions.QuestionSet, error)
	DeleteQuestionSet(ctx context.Context, id string) error
}

type QuestionSetHandler struct {
	questionSetService QuestionSetService
}

func NewQuestionSetHandler(s QuestionSetService) *QuestionSetHandler {

	h := &QuestionSetHandler{
		questionSetService: s,
	}

	return h
}

func (h *QuestionSetHandler) PostQuestionSet(w http.ResponseWriter, r *http.Request) {
	var qSet TrainingNeedsQuestions.QuestionSet

	if err := json.NewDecoder(r.Body).Decode(&qSet); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qSet, err := h.questionSetService.PostQuestionSet(r.Context(), qSet)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(qSet); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *QuestionSetHandler) GetQuestionSet(w http.ResponseWriter, r *http.Request) {
	qSetId := chi.URLParam(r, "id")

	qSet, err := h.questionSetService.GetQuestionSet(r.Context(), qSetId)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(qSet); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *QuestionSetHandler) GetQuestionSetByTechName(w http.ResponseWriter, r *http.Request) {
	technologyName := r.URL.Query().Get("technology-name")

	qSet, err := h.questionSetService.GetQuestionSetByTechName(r.Context(), technologyName)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(qSet); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *QuestionSetHandler) GetAllQuestionSets(w http.ResponseWriter, r *http.Request) {
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

	// Now call GetAllQuestionSets with page and pageSize as parameters
	qSets, err := h.questionSetService.GetAllQuestionSets(r.Context(), page, pageSize)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(qSets); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *QuestionSetHandler) UpdateQuestionSet(w http.ResponseWriter, r *http.Request) {
	var qSet TrainingNeedsQuestions.QuestionSet

	if err := json.NewDecoder(r.Body).Decode(&qSet); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qSet, err := h.questionSetService.UpdateQuestionSet(r.Context(), qSet)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(qSet); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *QuestionSetHandler) DeleteQuestionSet(w http.ResponseWriter, r *http.Request) {
	qSetId := chi.URLParam(r, "id")

	err := h.questionSetService.DeleteQuestionSet(r.Context(), qSetId)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *QuestionSetHandler) mapRoutes(router chi.Router) {
	router.Route("/api/v1/question-sets", func(r chi.Router) {
		r.Post("/", h.PostQuestionSet)

		r.With(h.queryParamMiddleware).Get("/", h.GetAllQuestionSets)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetQuestionSet)
			r.Put("/", h.UpdateQuestionSet)
			r.Delete("/", h.DeleteQuestionSet)
		})
	})
}
