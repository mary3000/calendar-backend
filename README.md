# calendar-backend
Backend for a google-like calendar

Server lives on a `9000` port.

## Run a server

1. `cd calendar/cmd`
2. `go run main.go` 

## HTTP API

_I prefer to use [Postman](https://www.postman.com/). It is more convenient than raw curl requests. You can put in Postman one of the curl request below and then play with it._

1. Add a user named "Susan"
```bash
curl --location --request POST 'localhost:9000/add-user' \
--header 'Content-Type: application/json' \
--data-raw '{"username":"Susan"}'
```

Sample answer:
```json
{"ID":1,"CreatedAt":"2022-07-15T19:36:31.800278606+03:00","UpdatedAt":"2022-07-15T19:36:31.800278606+03:00","DeletedAt":null,"Meetings":null,"Username":"Susan"}
```

2. Add meeting with specified name, guests and dates
```bash
curl --location --request POST 'localhost:9000/add-meeting' \
--header 'Content-Type: application/json' \
--data-raw '{"hostname":"Mary",
"guests":["Susan"],
"meetingname":"MarySusanMeeting",
"startdate":"2022-07-15T19:42:31.000+03:00",
"enddate":"2022-07-15T20:42:31.000+03:00",
"frequency":"u"
}'
```

Meeting frequency - see [code](calendar/internal/schema.go:17).

Sample answer:
```json
{"ID":3,"CreatedAt":"2022-07-15T19:44:27.225446291+03:00","UpdatedAt":"2022-07-15T19:44:27.225446291+03:00","DeletedAt":null,"Guests":[{"ID":1,"CreatedAt":"2022-07-15T19:36:31.800278606+03:00","UpdatedAt":"2022-07-15T19:36:31.800278606+03:00","DeletedAt":null,"Meetings":null,"Username":"Susan"},{"ID":2,"CreatedAt":"2022-07-15T19:37:13.874735869+03:00","UpdatedAt":"2022-07-15T19:37:13.874735869+03:00","DeletedAt":null,"Meetings":null,"Username":"Mary"}],"Slots":[{"ID":5,"CreatedAt":"2022-07-15T19:44:27.23855388+03:00","UpdatedAt":"2022-07-15T19:44:27.23855388+03:00","DeletedAt":null,"MeetingID":3,"UserID":2,"DefaultDecision":0,"OppositeDecisionDates":null},{"ID":6,"CreatedAt":"2022-07-15T19:44:27.255597625+03:00","UpdatedAt":"2022-07-15T19:44:27.255597625+03:00","DeletedAt":null,"MeetingID":3,"UserID":1,"DefaultDecision":0,"OppositeDecisionDates":null}],"MeetingName":"MarySusanMeeting","HostName":"Mary","StartDate":"2022-07-15T19:42:31+03:00","EndDate":"2022-07-15T20:42:31+03:00","Frequency":"u"}
```

3. Get meeting info by it's ID
```bash
curl --location --request POST 'localhost:9000/get-meeting' \
--header 'Content-Type: application/json' \
--data-raw '{"meetingid": 1}'
```

Sample answer:
```json
{"ID":1,"CreatedAt":"2022-07-15T19:37:27.549806415+03:00","UpdatedAt":"2022-07-15T19:37:27.549806415+03:00","DeletedAt":null,"Guests":[{"ID":1,"CreatedAt":"2022-07-15T19:36:31.800278606+03:00","UpdatedAt":"2022-07-15T19:36:31.800278606+03:00","DeletedAt":null,"Meetings":null,"Username":"Susan"},{"ID":2,"CreatedAt":"2022-07-15T19:37:13.874735869+03:00","UpdatedAt":"2022-07-15T19:37:13.874735869+03:00","DeletedAt":null,"Meetings":null,"Username":"Mary"}],"Slots":[{"ID":1,"CreatedAt":"2022-07-15T19:37:27.569429049+03:00","UpdatedAt":"2022-07-15T19:37:27.569429049+03:00","DeletedAt":null,"MeetingID":1,"UserID":2,"DefaultDecision":0,"OppositeDecisionDates":null},{"ID":2,"CreatedAt":"2022-07-15T19:37:27.586223338+03:00","UpdatedAt":"2022-07-15T19:37:27.586223338+03:00","DeletedAt":null,"MeetingID":1,"UserID":1,"DefaultDecision":0,"OppositeDecisionDates":null}],"MeetingName":"MarySusanMeeting","HostName":"Mary","StartDate":"2022-07-15T19:42:31+03:00","EndDate":"2022-07-15T20:42:31+03:00","Frequency":"u"}
```

4. Accept or decline the request (if date is null, then decision is made on all of the meetings at once).

```bash
curl --location --request POST 'localhost:9000/decide-on-meeting' \
--header 'Content-Type: application/json' \
--data-raw '{"meetingid": 1,
"userid": 3,
"accepted": true,
"date": null}'
```

Sample answer:
```json
{"ID":1,"CreatedAt":"2022-07-15T19:37:27.569429049+03:00","UpdatedAt":"2022-07-15T19:37:27.569429049+03:00","DeletedAt":null,"MeetingID":1,"UserID":2,"DefaultDecision":1,"OppositeDecisionDates":[]}
```

5. Get all user meetings, that occur in the specified time interval.

```bash
curl --location --request POST 'localhost:9000/get-user-meetings' \
--header 'Content-Type: application/json' \
--data-raw '{"username": "Mary",
"begindate": "2021-07-15T19:42:31.000+03:00",
"enddate": "2023-07-15T19:42:31.000+03:00"}'
```

Sample answer:
```json
[{"Slot":{"ID":1,"CreatedAt":"2022-07-15T19:37:27.569429049+03:00","UpdatedAt":"2022-07-15T21:35:06.990447896+03:00","DeletedAt":null,"MeetingID":1,"UserID":2,"DefaultDecision":1,"OppositeDecisionDates":null},"ConcreteTimeStart":"2022-07-15T19:42:31+03:00","ConcreteTimeEnd":"2022-07-15T20:42:31+03:00"},{"Slot":{"ID":3,"CreatedAt":"2022-07-15T19:43:32.853959694+03:00","UpdatedAt":"2022-07-15T21:36:51.638948535+03:00","DeletedAt":null,"MeetingID":2,"UserID":2,"DefaultDecision":1,"OppositeDecisionDates":null},"ConcreteTimeStart":"2022-07-16T19:42:31+03:00","ConcreteTimeEnd":"2022-07-16T20:42:31+03:00"},{"Slot":{"ID":5,"CreatedAt":"2022-07-15T19:44:27.23855388+03:00","UpdatedAt":"2022-07-15T21:37:10.896243389+03:00","DeletedAt":null,"MeetingID":3,"UserID":2,"DefaultDecision":1,"OppositeDecisionDates":null},"ConcreteTimeStart":"2022-07-15T19:42:31+03:00","ConcreteTimeEnd":"2022-07-15T20:42:31+03:00"}]
```

6. Get the nearest time in which all the given users a free for a given period of time.

```bash
curl --location --request POST 'localhost:9000/get-closest-meeting' \
--header 'Content-Type: application/json' \
--data-raw '{"names": ["Mary", "Susan"],
"duration": "1h"}'
```

Sample answer:
```json
"2022-07-15T22:03:26.778937897+03:00"
```
