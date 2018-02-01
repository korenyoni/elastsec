package event

import (
    "time"
    "encoding/json"
)

type Beat struct {
    Host    string      `json:"hostname"`
}

type Event struct {
    Message string      `json:"message"`
    Time    time.Time   `json:"@timestamp"`
    Beat    Beat        `json:"beat"`
    Type    string
    Source  *json.RawMessage
}
