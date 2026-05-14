package response

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// JSON writes a JSON response with the given status code and data.
// It sets the Content-Type header to application/json.
// If marshaling fails, it writes a 500 Internal Server Error response.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// At this point the header may already be partially sent,
		// so we attempt to write the error body only if possible.
		http.Error(w, `{"error":"internal server error"}`+"\n", http.StatusInternalServerError)
	}
}

// Error writes a JSON error response with the given status code and message.
// The response body is {"error": "message"}.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

// DecodeJSON decodes the JSON request body into the provided destination.
// It returns an error if the body is empty or if decoding fails.
func DecodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return errors.New("request body is empty")
	}

	if err := json.Unmarshal(body, dst); err != nil {
		return err
	}

	return nil
}
