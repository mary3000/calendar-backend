package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AddUserRequest struct {
	Username string
}

// param: name
func AddUser(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		http.Error(w, fmt.Sprintf("Content-type: expected %v, got %v", expectedContentType, contentType),
			http.StatusBadRequest)
		return
	}

	expectedMethod := "POST"
	if r.Method != expectedMethod {
		http.Error(w, fmt.Sprintf("Method: expected %v, got %v", expectedMethod, r.Method), http.StatusBadRequest)
	}

	var req AddUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &User{Username: req.Username}
	res := Db.Create(user)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Added user: %+v", user)

	payloadBytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(payloadBytes)
}
