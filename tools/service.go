package tools

import (
	"encoding/json"
	"net/http"
)

// HandlerError write error back to response and record error with log
func HandlerError(w http.ResponseWriter, err error, msg string, statusCode int) {
	errorMsg := ErrorMsg{
		Msg: msg,
	}
	jsonMsg, err := json.Marshal(errorMsg)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Can not marshal json to string"))
		return
	}
	w.WriteHeader(statusCode)
	w.Write([]byte(jsonMsg))
	// Log the error
}