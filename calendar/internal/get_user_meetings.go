package internal

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"time"
)

type GetUserMeetingsRequest struct {
	Username  string
	BeginDate time.Time
	EndDate   time.Time
}

type MeetingTimeSlot struct {
	Slot              MeetingSlot
	ConcreteTimeStart time.Time
	ConcreteTimeEnd   time.Time
}

func GetUserMeetings(w http.ResponseWriter, r *http.Request) {
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

	var req GetUserMeetingsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u := User{Username: req.Username}
	Db. /*.Preload("Meetings").Preload("Slots")*/ First(&u)

	var meetingSlots []MeetingSlot
	Db.Find("UserID = ?", u.ID).Find(&meetingSlots)

	mts := GetMeetingsInInterval(meetingSlots, req.BeginDate, req.EndDate)

	log.Printf("Got user meetings: %+v", mts)

	payloadBytes, err := json.Marshal(mts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(payloadBytes)
}

func GetMeetingsInInterval(meetingSlots []MeetingSlot, begin time.Time, end time.Time) []MeetingTimeSlot {
	mts := make([]MeetingTimeSlot, 0)

	for _, ms := range meetingSlots {
		var m Meeting
		Db.Find(&m, ms.MeetingID)

		mts = append(mts, GetMeetings(ms, m, begin, end)...)
	}

	return mts
}

func GetUnrepeatedMeetings(ms MeetingSlot, m Meeting, begin time.Time, end time.Time) []MeetingTimeSlot {
	if ms.DefaultDecision != Accepted {
		return []MeetingTimeSlot{}
	}

	if !(m.StartDate.After(end) || m.EndDate.Before(begin)) {
		return []MeetingTimeSlot{{
			Slot:              ms,
			ConcreteTimeStart: m.StartDate,
			ConcreteTimeEnd:   m.EndDate,
		}}
	} else {
		return []MeetingTimeSlot{}
	}
}

func (f MeetingFrequency) Beginning(begin time.Time, end time.Time) (years, months, days int) {
	switch f {
	case Unrepeated:
		panic("unsupported")
	case Daily:
		return 0, 0, end.Day() - begin.Day()
	case Weekly:
		return 0, 0, (end.Day() - begin.Day()) / 7
	case Monthly:
		return 0, int(end.Month() - begin.Month()), 0
	case Annually:
		return end.Year() - begin.Year(), 0, 0
	}
	panic("unreachable")
}

func (f MeetingFrequency) Next(years, months, days int) (years2, months2, days2 int) {
	switch f {
	case Unrepeated:
		panic("unsupported")
	case Daily:
		return years, months, days + 1
	case Weekly:
		return years, months, days + 7
	case Monthly:
		return years, months + 1, days
	case Annually:
		return years + 1, months, days
	}
	panic("unreachable")
}

func GetMeetings(ms MeetingSlot, m Meeting, begin time.Time, end time.Time) []MeetingTimeSlot {
	if m.Frequency == Unrepeated {
		return GetUnrepeatedMeetings(ms, m, begin, end)
	}

	cur := m.StartDate
	years, months, days := 0, 0, 0
	if cur.Before(begin) {
		years, months, days = m.Frequency.Beginning(begin, cur)
		cur = m.StartDate.AddDate(years, months, days)
		if cur.Before(begin) {
			years, months, days = m.Frequency.Next(years, months, days)
			cur = m.StartDate.AddDate(years, months, days)
		}
	}

	acceptedTimes := make([]MeetingTimeSlot, 0)

	for !(cur.After(end)) {
		if ms.DefaultDecision == Accepted && slices.Index(ms.OppositeDecisionDates, cur) == -1 ||
			ms.DefaultDecision != Accepted && slices.Index(ms.OppositeDecisionDates, cur) != -1 {
			acceptedTimes = append(acceptedTimes, MeetingTimeSlot{
				Slot:              ms,
				ConcreteTimeStart: cur,
				ConcreteTimeEnd:   cur.Add(m.EndDate.Sub(m.StartDate)),
			})
		}
		years, months, days = m.Frequency.Next(years, months, days)
		cur = m.StartDate.AddDate(years, months, days)
	}

	return acceptedTimes
}

/*func GetDailyMeetings(ms MeetingSlot, m Meeting, begin time.Time, end time.Time) {
	cur := m.StartDate
	if cur.Before(begin) {
		diff := cur.Sub(begin)
		cur.Add(Day * time.Duration(diff.Hours() / 24))
		if cur.Before(begin) {
			cur.Add(Day)
		}
	}

	acceptedTimes := make([]MeetingTimeSlot, 0)

	for !(cur.After(end)) {
		acceptedTimes = append(acceptedTimes, MeetingTimeSlot{
			Slot:              ms,
			ConcreteTimeStart: cur,
			ConcreteTimeEnd:   cur.Add(m.EndDate.Sub(m.StartDate)),
		})
		cur.Add(Day)
	}
}*/
