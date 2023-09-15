package http

import "net/http"

func (h *QuestionSetHandler) qSetQueryParamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name != "" {
			h.GetQuestionSetByTechName(w, r)
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
			h.GetCourseOutlinesByFilter(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
