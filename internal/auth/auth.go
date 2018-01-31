package auth

import (
    "../looper"
    "../errors"
    "fmt"
    "log"
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
    authFailure := regexp.MustCompile("authentication failure")
    notInSudoers := regexp.MustCompile("NOT in sudoers")
    command := regexp.MustCompile("COMMAND=.*$")

    regExpressions := []regexp.Regexp{*authFailure,*notInSudoers,*command}
    _ = regExpressions

    // matches
    matchCommand := command.FindString(e.Message)
    if matchCommand != "" {
        user := regexp.MustCompile("sudo:")
        matchUser := user.FindString(e.Message)
        if matchUser == "" {
            log.Fatal(errors.CreateMatchError(*user,"message"))
        }
        events <- fmt.Sprintf("time: %s command: %s\n", e.Time, matchCommand)
    }
}
