package event

import (
    "time"
)

type Beat struct {
    Host    string      `json:"hostname"`
}

type Event struct {
    Message string      `json:"message"`
    Time    time.Time   `json:"@timestamp"`
    Beat    Beat        `json:"beat"`
    Type    string
}
