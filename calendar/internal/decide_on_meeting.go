package internal

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"time"
)

type DecideOnMeetingRequest struct {
	MeetingID uint
	UserID    uint
	Accepted  bool
	Date      time.Time
}

func DecideOnMeeting(w http.ResponseWriter, r *http.Request) {
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

	var req DecideOnMeetingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var slot MeetingSlot
	Db.Where("meeting_id = ? AND user_id = ?", req.MeetingID, req.UserID).First(&slot)

	if req.Date.IsZero() {
		if req.Accepted {
			slot.DefaultDecision = Accepted
		} else {
			slot.DefaultDecision = Declined
		}
		slot.OppositeDecisionDates = make([]time.Time, 0)
	} else {
		// Date already inside oppdd
		idx := slices.IndexFunc(slot.OppositeDecisionDates, func(t time.Time) bool { return t == req.Date })
		if idx >= 0 {
			if slot.DefaultDecision == Accepted && req.Accepted || slot.DefaultDecision != Accepted && !req.Accepted {
				slot.OppositeDecisionDates = slices.Delete(slot.OppositeDecisionDates, idx, idx)
			} else {
				// do nothing, already correct decision
			}
		} else {
			// Date is not in oppdd
			if slot.DefaultDecision == Accepted && !req.Accepted || slot.DefaultDecision != Accepted && req.Accepted {
				slot.OppositeDecisionDates = append(slot.OppositeDecisionDates, req.Date)
			}
		}
	}

	res := Db.Save(&slot)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Decided on slot: %+v", slot)

	payloadBytes, err := json.Marshal(slot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(payloadBytes)
}
