package http

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (h *QuestionSetHandler) qSetQueryParamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name != "" {
			h.GetQuestionSetByTechName(w, r, name)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *CourseOutlineHandler) outlineQueryParamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filterName := r.URL.Query().Get("filterName")
		filterValue := r.URL.Query().Get("filterValue")
		if filterName != "" && filterValue != "" {
			h.GetCourseOutlinesByFilter(w, r, filterName, filterValue)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			log.Error("invalid authorization header")
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")

		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			log.Error("invalid authorization header")
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		token := authHeaderParts[1]
		if VerifyFirebaseToken(r.Context(), token) {
			next.ServeHTTP(w, r) // this will call the next handler in the chain
		} else {
			log.Error("unauthorized authorization header")
			http.Error(w, "not authorized", http.StatusUnauthorized)
		}
	})
}
