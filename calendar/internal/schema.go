package internal

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type User struct {
	gorm.Model

	Meetings []*Meeting `gorm:"many2many:user_meeting;"`

	Username string `gorm:"unique"`
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
	Unrepeated MeetingFrequency = "u"
	Daily                       = "d"
	Weekly                      = "w"
	Monthly                     = "m"
	Annually                    = "a"
)

type Meeting struct {
	gorm.Model

	Guests []*User `gorm:"many2many:user_meeting;"` // including host
	Slots  []MeetingSlot

	MeetingName string
	HostName    string
	StartDate   time.Time
	EndDate     time.Time
	Frequency   MeetingFrequency
}

type Decision uint

const (
	Unknown Decision = iota
	Accepted
	Declined
)

type MeetingSlot struct {
	gorm.Model

	MeetingID             uint // foreign key
	UserID                uint
	DefaultDecision       Decision
	OppositeDecisionDates []time.Time
}

var Db *gorm.DB
