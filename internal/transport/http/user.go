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
	tnauser "github.com/zzenonn/trainocate-tna/internal/user"
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

type UserServiceInterface interface {
	GetUser(ctx context.Context, id string) (tnauser.User, error)
	GetAllUsers(ctx context.Context, page int, pageSize int) ([]tnauser.User, error)
	CreateUser(ctx context.Context, user tnauser.User) (tnauser.User, error)
	UpdateUser(ctx context.Context, user tnauser.User) (tnauser.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserHandler struct {
	userService UserServiceInterface
}

func NewUserHandler(s UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

func (h *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	var user tnauser.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), user)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "id")
	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(r.Context(), uid)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
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

	// Now call GetAllUsers with page and pageSize as parameters
	qSets, err := h.userService.GetAllUsers(r.Context(), page, pageSize)

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

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "id")

	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user tnauser.User

	user.ID = uid

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), user)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "id")

	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.userService.DeleteUser(r.Context(), uid)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) mapRoutes(router chi.Router) {
	router.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/", h.PostUser)
		r.Get("/", h.GetAllUsers)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetUser)
			r.Put("/", h.UpdateUser)
			r.Delete("/", h.DeleteUser)
		})
	})
}
