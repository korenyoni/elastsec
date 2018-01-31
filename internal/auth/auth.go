package auth

import (
    "../looper"
    "fmt"
    "regexp"
)

const (
    indexName   = "filebeat*"
)

func Loop(events chan<- string) {
    eventBus := make(chan looper.Event)

    go looper.Loop(eventBus, "filebeat*")

    for event := range eventBus {
        parseEvent(events, event)
    }
}

func parseEvent(events chan<- string, e looper.Event) {
    r := regexp.MustCompile("COMMAND=.*$")
    match := r.FindString(e.Message)
    if match != "" {
        events <- fmt.Sprintf("time: %s command: %s\n", e.Time, match)
    }
}
