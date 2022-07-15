package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GetMeetingRequest struct {
	// Hostname    string
	// MeetingName string
	MeetingID uint
}

func GetMeeting(w http.ResponseWriter, r *http.Request) {
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

	var req GetMeetingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var m Meeting
	res := Db.Preload("Guests").Preload("Slots").First(&m, req.MeetingID)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Got meeting: %+v", m)

	/*_, _ = w.Write([]byte(fmt.Sprintf("Meeting details: \n"+
	"{\n"+
	"MeetingName: %v, \n"+
	"Guests: %v, \n"+
	"HostName: %v, \n"+
	"StartDate: %v, \n"+
	"EndDate: %v, \n"+
	"Frequency: %v \n"+
	"} \n"+
	"obj: \n"+
	"%v", m.MeetingName, m.Guests, m.HostName, m.StartDate, m.EndDate, m.Frequency, m)))*/

	payloadBytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(payloadBytes)
}
