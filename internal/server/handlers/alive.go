package handlers

import (
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

//Welcome test handler
func Alive(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("ID").(uuid.UUID)
	outObj := struct {
		Message string `json:"message"`
	}{
		requestID.String(),
	}
	resp, _ := json.Marshal(outObj)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
