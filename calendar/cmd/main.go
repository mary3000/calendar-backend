package main

import (
	"calendar-backend/calendar/internal"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

var port = "9000"

var dbName = "db"

func main() {
	var err error
	internal.Db, err = gorm.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer internal.Db.Close()

	log.Print("internal.Db opened successfully")

	internal.Db.LogMode(true)
	internal.Db.AutoMigrate(&internal.User{})
	internal.Db.AutoMigrate(&internal.Meeting{})
	internal.Db.AutoMigrate(&internal.MeetingSlot{})
	mux := http.NewServeMux()

	mux.HandleFunc("/add-user", internal.AddUser)
	mux.HandleFunc("/add-meeting", internal.AddMeeting)
	mux.HandleFunc("/get-meeting", internal.GetMeeting)
	mux.HandleFunc("/decide-on-meeting", internal.DecideOnMeeting)
	mux.HandleFunc("/get-user-meetings", internal.GetUserMeetings)
	mux.HandleFunc("/get-closest-meeting", internal.GetClosestMeeting)

	// Helper methods
	mux.HandleFunc("/get-users", internal.GetUsers)

	chat := http.Server{Addr: ":" + port, Handler: mux}
	log.Fatal(chat.ListenAndServe())
}
