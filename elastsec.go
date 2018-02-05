package main

import (
    "fmt"
    "./internal/beats/auth"
    "./internal/beats/file_integrity"
    "./internal/beats/filechange_attempt"
    "./internal/event"
    "./internal/notify"
    "./internal/infoexport"
)

func main() {
    eventBus := make(chan event.Event)

    go auth.Loop(eventBus)
    go file_integrity.Loop(eventBus)
    go filechange_attempt.Loop(eventBus)

    for event := range eventBus {
        title := infoexport.GetTitle(event)
        fmt.Printf("%s %s\n\n",title,event.Message)
        notify.SendSlack(event,title)
    }
}
