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

	var req AddMeetingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var retrievedUsers []User
	allUsers := append(req.Guests, req.Hostname)
	Db.Where("username IN (?)", allUsers).Find(&retrievedUsers)

	createdMeeting := Meeting{
		MeetingName: req.MeetingName,
		HostName:    req.Hostname,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Frequency:   MeetingFrequency(req.Frequency),
		Slots:       []MeetingSlot{},
	}
	Db.Create(&createdMeeting)

	for _, user := range retrievedUsers {
		Db.Model(&user).Association("Meetings").Append(&createdMeeting)
		// Db.Model(&createdMeeting).Association("Guests").Append(&retrievedUsers)

		slot := MeetingSlot{
			MeetingID:             createdMeeting.ID,
			UserID:                user.ID,
			DefaultDecision:       Unknown,
			OppositeDecisionDates: []time.Time{},
		}
		Db.Create(&slot)
		Db.Model(&createdMeeting).Association("MeetingSlots").Append(&slot)
	}

	log.Printf("Added meeting: %+v", createdMeeting)

	_, _ = w.Write([]byte(fmt.Sprint(createdMeeting.ID)))
}
