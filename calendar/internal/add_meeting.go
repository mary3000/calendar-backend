package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AddMeetingRequest struct {
	Hostname    string
	Guests      []string
	MeetingName string
	StartDate   time.Time
	EndDate     time.Time
	Frequency   string
}

// param: username, []guests_names, meeting_name, time, length, repeated: [no, d, w, m, a]
func AddMeeting(w http.ResponseWriter, r *http.Request) {
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

	var m AddMeetingRequest
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*users := make([]User, len(m.Guests) + 1)
	users[0].Username = m.Hostname
	for i, _ := range m.Guests {
		users[i+1].Username = m.Guests[i]
	}*/

	var retrievedUsers []User
	allUsers := append(m.Guests, m.Hostname)
	Db.Where("username IN (?)", allUsers).Find(&retrievedUsers)

	createdMeeting := Meeting{
		MeetingName: m.MeetingName,
		HostName:    m.Hostname,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		Frequency:   MeetingFrequency(m.Frequency),
	}
	Db.Create(&createdMeeting)

	// Append users
	err = Db.Model(&createdMeeting).Association("Guests").Append(&retrievedUsers).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Added meeting: %v", createdMeeting)

	_, _ = w.Write([]byte(fmt.Sprint(createdMeeting)))
}
