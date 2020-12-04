package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Error string      `json:"error"`
	Key   string      `json:"key"`
	Data  interface{} `json:"data,omitempty"`
}

type Payload map[string]interface{}

func Err(error, key string, data ...interface{}) Error {
	if len(data) == 0 {
		return Error{Error: error, Key: key}
	}
	return Error{Error: error, Key: key, Data: data[0]}
}

func Ok(w http.ResponseWriter, payload interface{}) int {
	return respondWithPayload(http.StatusOK, w, payload)
}

func BadRequest(w http.ResponseWriter, e Error) int {
	return respondWithError(http.StatusBadRequest, w, e)
}

func Unauthorized(w http.ResponseWriter, e Error) int {
	return respondWithError(http.StatusUnauthorized, w, e)
}

func Forbidden(w http.ResponseWriter, e Error) int {
	return respondWithError(http.StatusForbidden, w, e)
}

func NotFound(w http.ResponseWriter, e Error) int {
	return respondWithError(http.StatusNotFound, w, e)
}

func InternalError(w http.ResponseWriter, e Error) int {
	return respondWithError(http.StatusInternalServerError, w, e)
}

func Teapot(w http.ResponseWriter, e Error) int {
	return respondWithError(http.StatusTeapot, w, e)
}

func WithErr(code int, w http.ResponseWriter, e Error) int {
	return respondWithError(code, w, e)
}

func With(code int, w http.ResponseWriter, payload interface{}) int {
	return respondWithPayload(code, w, payload)
}

func respondWithError(status int, w http.ResponseWriter, e Error) int {
	resp, err := json.Marshal(e)
	if err != nil {
		resp = []byte(fmt.Sprintf(`{"error": "%s", "key": "%s"}`, e.Error, e.Key))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
	return status
}

func respondWithPayload(status int, w http.ResponseWriter, payload interface{}) int {
	resp, err := json.Marshal(payload)
	if err != nil {
		status = http.StatusInternalServerError
		resp = []byte(`{"error": "something went wrong"}`)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
	return status
}
