package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	mapRoutes(router chi.Router)
}

type MainHandler struct {
	Router   chi.Router
	Handlers []Handler
	Server   *http.Server
}

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
	} // User Routes

	// Comment Routes

}

func NewMainHandler() *MainHandler {
	h := &MainHandler{
		Handlers: []Handler{},
	}

	h.Router = chi.NewRouter()

	h.Router.Use()

	// h.Router.Use(CorsMiddleware)

	h.Router.Use(logger.Logger("router", log.New()))
	// h.Router.Use(JSONMiddleware)
	// h.Router.Use(TimeoutMiddleware)

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}

	return h
}

// Each rest object handler appends itself to the main handler
func (h *MainHandler) AddHandler(handler Handler) {
	h.Handlers = append(h.Handlers, handler)
}

func (h *MainHandler) MapRoutes() {
	// h.Router.Options("/", func(w http.ResponseWriter, r *http.Request) {})

	// corsMiddleware := cors.Handler(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{
	// 		http.MethodHead,
	// 		http.MethodGet,
	// 		http.MethodPost,
	// 		http.MethodPut,
	// 		http.MethodPatch,
	// 		http.MethodDelete,
	// 		http.MethodOptions,
	// 	},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// })

	h.Router.Use(CorsMiddleware)
	h.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	for _, handler := range h.Handlers {
		handler.mapRoutes(h.Router)
		chi.Walk(h.Router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			log.Debugf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
			return nil
		})
	}

}

func (h *MainHandler) Serve() error {

	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	h.Server.Shutdown(ctx)

	log.Warn("shutting down gracefully")

	return nil
}
