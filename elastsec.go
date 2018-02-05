package main

import (
    "fmt"
    "./internal/beats/auth"
    "./internal/beats/file_integrity"
    "./internal/beats/filechange_attempt"
    "./internal/event"
    "./internal/notify"
)

func main() {
    eventBus := make(chan event.Event)

    go auth.Loop(eventBus)
    go file_integrity.Loop(eventBus)
    go filechange_attempt.Loop(eventBus)

    for event := range eventBus {
        stdoutTitle := fmt.Sprintf("New `%s` event on host `%s`: %s",
        event.Type, event.Beat.Host, event.Message)
        notifyTitle := fmt.Sprintf("New `%s` event on host `%s`",
        event.Type, event.Beat.Host)
        fmt.Printf("%s\n\n",stdoutTitle)
        notify.SendSlack(event,notifyTitle)
    }
}
