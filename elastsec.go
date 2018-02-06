package main

import (
    "fmt"
    "time"
    "./internal/beats/auth"
    "./internal/beats/file_integrity"
    "./internal/beats/filechange_attempt"
    "./internal/event"
    "./internal/notify"
    "./internal/infoexport"
    "./internal/aggregator"
)

func main() {
    eventBus := make(chan event.Event)
    var a aggregator.Aggregator
    a.SupressedCount= make(map[string]int)

    go auth.Loop(eventBus)
    go file_integrity.Loop(eventBus)
    go filechange_attempt.Loop(eventBus)
    go a.Loop(eventBus, time.Minute)

    for event := range eventBus {
        event, ok := a.Consume(event)
        if ok {
            title := infoexport.GetTitle(event)
            fmt.Printf("%s %s\n\n",title,event.Message)
            notify.SendSlack(event,title)
        }
    }
}
