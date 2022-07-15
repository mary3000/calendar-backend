package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		http.Error(w, fmt.Sprintf("Content-type: expected %v, got %v", expectedContentType, contentType),
			http.StatusBadRequest)
		return
	}

	expectedMethod := "GET"
	if r.Method != expectedMethod {
		http.Error(w, fmt.Sprintf("Method: expected %v, got %v", expectedMethod, r.Method), http.StatusBadRequest)
	}

	var users []User

	res := Db.Find(&users)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Got users: %+v", users)

	payloadBytes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(payloadBytes)
}
