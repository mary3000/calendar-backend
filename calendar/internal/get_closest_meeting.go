package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type GetClosestMeetingRequest struct {
	Names    []string
	Duration string
}

func GetClosestMeeting(w http.ResponseWriter, r *http.Request) {
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

	var req GetClosestMeetingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var users []User
	Db.Where("username IN (?)", req.Names).Find(&users)

	meetingSlots := make([][]MeetingSlot, len(users))
	for i := range users {
		Db.Where("user_id = ?", users[i].ID).Find(&meetingSlots[i])
	}

	curTime := time.Now().Add(time.Minute)
	length, err := time.ParseDuration(req.Duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// potentially infinite!
	// todo: make loop finite
Again:
	for i := range users {
		for _, ms := range meetingSlots[i] {
			var m Meeting
			Db.Find(&m, ms.MeetingID)
			slotsInInterval := GetMeetings(ms, m, curTime, curTime.Add(length))
			for _, mts := range slotsInInterval {
				if mts.ConcreteTimeEnd.After(curTime) {
					curTime = mts.ConcreteTimeEnd
					goto Again
				}
			}
		}
	}

	log.Printf("Closest time: %v", curTime)

	payloadBytes, err := json.Marshal(curTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(payloadBytes)
}
