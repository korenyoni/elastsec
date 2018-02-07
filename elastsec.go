package main

import (
    "fmt"
    "time"
    "os"
    "log"
    "./internal/beats/auth"
    "./internal/beats/file_integrity"
    "./internal/beats/filechange_attempt"
    "./internal/event"
    "./internal/notify"
    "./internal/infoexport"
    "./internal/aggregator"
    "./internal/constants"
)

func main() {
    eventBus := make(chan event.Event)
    aggregatedEventBus := make(chan event.Event)
    var a aggregator.Aggregator
    a.SupressedCount = make(map[aggregator.Key]int)

    go auth.Loop(eventBus)
    go file_integrity.Loop(eventBus)
    go filechange_attempt.Loop(eventBus)
    go a.Loop(eventBus, getAggDuration())

    go func() {
        for event := range eventBus {
            event, ok := a.Consume(event)
            if ok {
                aggregatedEventBus <- event
            }
        }
    }()

    email := notify.EmailInit(aggregatedEventBus, getEmailDuration())
    go email.Loop()
    for event := range aggregatedEventBus {
        title := infoexport.GetTitle(event)
        fmt.Printf("%s %s\n\n",title,event.Message)
        notify.SendSlack(event,title)
        email.Consume(event,title)
    }
}

func getAggDuration() time.Duration {
    durationString := os.Getenv(constants.AggDurationEnv)
    if durationString != "" {
        duration, err := time.ParseDuration(durationString)
        if err != nil {
            log.Println("Invalid aggregator duration.")
        } else {
            return duration
        }
    }
    return time.Hour
}

func getEmailDuration() time.Duration {
    durationString := os.Getenv(constants.EmailDurationEnv)
    if durationString != "" {
    duration, err := time.ParseDuration(durationString)
        if err != nil {
            log.Println("Invalid email duration.")
        } else {
            return duration
        }
    }
    return time.Hour
}
