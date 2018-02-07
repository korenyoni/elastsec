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
    aggregatedEventBus := make(chan event.Event)
    var a aggregator.Aggregator
    a.SupressedCount = make(map[aggregator.Key]int)

    go auth.Loop(eventBus)
    go file_integrity.Loop(eventBus)
    go filechange_attempt.Loop(eventBus)
    go a.Loop(eventBus, time.Minute)

    go func() {
        for event := range eventBus {
            event, ok := a.Consume(event)
            if ok {
                aggregatedEventBus <- event
            }
        }
    }()

    email := notify.EmailInit(aggregatedEventBus, time.Minute)
    go email.Loop()
    for event := range aggregatedEventBus {
        title := infoexport.GetTitle(event)
        fmt.Printf("%s %s\n\n",title,event.Message)
        notify.SendSlack(event,title)
        email.Consume(event,title)
    }
}
