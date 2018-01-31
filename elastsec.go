package main

import (
    "fmt"
    "./internal/auth"
    "./internal/event"
)

func main() {
    eventBus := make(chan event.Event)

    go auth.Loop(eventBus)

    for event := range eventBus {
        fmt.Println(event.Message)
    }
}
