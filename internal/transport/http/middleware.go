package http

import "net/http"

func (h *QuestionSetHandler) queryParamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name != "" {
			h.GetQuestionSetByTechName(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
