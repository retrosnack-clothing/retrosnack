package httputil

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorBody struct {
	Error string `json:"error"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Error writes a JSON error response, hiding internal details for 5xx errors.
func Error(w http.ResponseWriter, status int, err error) {
	if status >= 500 {
		slog.Error("internal error", "status", status, "error", err)
		JSON(w, status, errorBody{Error: "internal server error"})
		return
	}
	JSON(w, status, errorBody{Error: err.Error()})
}

// ErrorMsg writes a JSON error response from a string message.
func ErrorMsg(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, errorBody{Error: msg})
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
