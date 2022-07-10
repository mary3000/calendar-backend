package internal

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type User struct {
	gorm.Model

	Username string     `gorm:"unique"`
	Meetings []*Meeting `gorm:"many2many:user_meeting;"`
}

// Types of meetings:
// 1) Unaccepted meeting
// 2) No-repeated meeting
// 3) Daily m.
// 4) Weekly
// 5) Monthly
// 6) Annually

type MeetingFrequency string

const (
	Unrepeated = "u"
	Daily      = "d"
	Weekly     = "w"
	Monthly    = "m"
	Annually   = 'a'
)

type Meeting struct {
	gorm.Model

	MeetingName string
	Guests      []*User `gorm:"many2many:user_meeting;"`
	HostName    string
	StartDate   time.Time
	EndDate     time.Time
	Frequency   MeetingFrequency
}

var Db *gorm.DB
