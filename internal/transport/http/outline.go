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
	outline "github.com/zzenonn/trainocate-tna/internal/outline"
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

type CourseOutlineService interface {
	PostCourseOutline(ctx context.Context, courseOutline outline.CourseOutline) (outline.CourseOutline, error)
	GetCourseOutline(ctx context.Context, id string) (outline.CourseOutline, error)
	GetCourseOutlinesByFilter(ctx context.Context, page int, pageSize int, filterName string, filterValue string) ([]outline.CourseOutline, error)
	GetAllCourseOutlines(ctx context.Context, page int, pageSize int) ([]outline.CourseOutline, error)
	UpdateCourseOutline(ctx context.Context, courseOutline outline.CourseOutline) (outline.CourseOutline, error)
	DeleteCourseOutline(ctx context.Context, id string) error
}

type CourseOutlineHandler struct {
	courseOutlineService CourseOutlineService
}

func NewCourseOutlineHandler(s CourseOutlineService) *CourseOutlineHandler {
	return &CourseOutlineHandler{
		courseOutlineService: s,
	}
}

func (h *CourseOutlineHandler) PostCourseOutline(w http.ResponseWriter, r *http.Request) {
	var courseOutline outline.CourseOutline

	if err := json.NewDecoder(r.Body).Decode(&courseOutline); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courseOutline, err := h.courseOutlineService.PostCourseOutline(r.Context(), courseOutline)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(courseOutline); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *CourseOutlineHandler) GetCourseOutline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	courseOutline, err := h.courseOutlineService.GetCourseOutline(r.Context(), id)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(courseOutline); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *CourseOutlineHandler) GetCourseOutlinesByFilter(w http.ResponseWriter, r *http.Request, filterName string, filterValue string) {
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

	courseOutlines, err := h.courseOutlineService.GetCourseOutlinesByFilter(r.Context(), page, pageSize, filterName, filterValue)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(courseOutlines); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *CourseOutlineHandler) GetAllCourseOutlines(w http.ResponseWriter, r *http.Request) {
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

	courseOutlines, err := h.courseOutlineService.GetAllCourseOutlines(r.Context(), page, pageSize)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(courseOutlines); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *CourseOutlineHandler) UpdateCourseOutline(w http.ResponseWriter, r *http.Request) {
	courseOutlineId := chi.URLParam(r, "id")

	var courseOutline outline.CourseOutline

	if err := json.NewDecoder(r.Body).Decode(&courseOutline); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courseOutline.Id = courseOutlineId

	courseOutline, err := h.courseOutlineService.UpdateCourseOutline(r.Context(), courseOutline)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(courseOutline); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *CourseOutlineHandler) DeleteCourseOutline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.courseOutlineService.DeleteCourseOutline(r.Context(), id)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *CourseOutlineHandler) mapRoutes(router chi.Router) {
	router.Route("/api/v1/course-outlines", func(r chi.Router) {

		r.Use(JwtMiddleware)

		r.Post("/", h.PostCourseOutline)

		r.With(h.outlineQueryParamMiddleware).Get("/", h.GetAllCourseOutlines)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetCourseOutline)
			r.Put("/", h.UpdateCourseOutline)
			r.Delete("/", h.DeleteCourseOutline)
		})

	})
}
