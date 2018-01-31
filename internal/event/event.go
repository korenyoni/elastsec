package event

import (
    "time"
)

type Event struct {
    Host    string    `json:"beat.hostname"`
    Message string    `json:"message"`
    Time    time.Time `json:"@timestamp"`
    EventType string
}
