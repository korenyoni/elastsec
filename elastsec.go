package main

import (
    "fmt"
    "./internal/beats/auth"
    "./internal/beats/file_integrity"
    "./internal/beats/filechange_attempt"
    "./internal/event"
)

func main() {
    eventBus := make(chan event.Event)

    go auth.Loop(eventBus)
    go file_integrity.Loop(eventBus)
    go filechange_attempt.Loop(eventBus)

    for event := range eventBus {
        fmt.Printf("New `%s` event on host `%s`: %s\n",
        event.Type, event.Beat.Host, event.Message)
    }
}
