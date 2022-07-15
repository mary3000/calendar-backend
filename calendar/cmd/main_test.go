package main

import (
	"bytes"
	"calendar-backend/calendar/internal"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

var testDbName = "test_db"

func Server() {
	dbName = testDbName
	_ = os.Remove(testDbName)
	main()
}

// Methods below are based on curl-to-Go: https://mholt.github.io/curl-to-go

func testAddUser(request internal.AddUserRequest) error {
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/add-user", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func testAddMeeting(request internal.AddMeetingRequest) error {
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/add-meeting", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func testGetMeeting(request internal.GetMeetingRequest) (*internal.Meeting, error) {
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/get-meeting", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m internal.Meeting
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func testGetUserMeetings(request internal.GetUserMeetingsRequest) (*[]internal.MeetingTimeSlot, error) {
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/get-user-meetings", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mts []internal.MeetingTimeSlot
	err = json.NewDecoder(resp.Body).Decode(&mts)
	if err != nil {
		return nil, err
	}

	return &mts, nil
}

func TestAllAPI(t *testing.T) {
	go Server()

	time.Sleep(time.Millisecond * 300) // wait for db startup

	if err := testAddUser(internal.AddUserRequest{Username: "Mary"}); err != nil {
		t.Errorf("err = %v", err)
	}

	if err := testAddUser(internal.AddUserRequest{Username: "Susan"}); err != nil {
		t.Errorf("err = %v", err)
	}

	if err := testAddUser(internal.AddUserRequest{Username: "Tom"}); err != nil {
		t.Errorf("err = %v", err)
	}

	now := time.Now()
	if err := testAddMeeting(internal.AddMeetingRequest{
		Hostname:    "Mary",
		Guests:      []string{"Susan"},
		MeetingName: "meeting1",
		StartDate:   now.Add(time.Hour),
		EndDate:     now.Add(time.Hour * 2),
		Frequency:   "u",
	}); err != nil {
		t.Errorf("err = %v", err)
	}

	m, err := testGetMeeting(internal.GetMeetingRequest{
		MeetingID: 1,
	})
	if err != nil {
		t.Errorf("err = %v", err)
	}

	// add m checks?
	_ = m

}
