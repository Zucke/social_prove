package response

import (
	"encoding/json"
	"net/http"
)

// Map is an alias for map[string]interface{}, this makes
// it easier to work with objects of undefined structure.
type Map map[string]interface{}

// ErrorMessage standardized error response.
type ErrorMessage struct {
	Message string `json:"message"`
}

// HTTPError standardized error response in JSON format.
func HTTPError(w http.ResponseWriter, statusCode int, message string) error {
	msg := ErrorMessage{
		Message: message,
	}

	return JSON(w, statusCode, msg)
}

// JSON standarized JSON response.
func JSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	if data == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(j)
	return nil
}
