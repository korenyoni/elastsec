package auth

import (
    "../looper"
    "../errors"
    "../event"
    "fmt"
    "log"
    "regexp"
)

const (
    indexName   = "filebeat*"
)

func Loop(events chan<- event.Event) {
    eventBus := make(chan event.Event)

    go looper.Loop(eventBus, "filebeat*")

    for event := range eventBus {
        replaceMessage(events, event)
    }
}

func replaceMessage(events chan<- event.Event, e event.Event) {
    authFailure := regexp.MustCompile("authentication failure")
    notInSudoers := regexp.MustCompile("NOT in sudoers")
    command := regexp.MustCompile("COMMAND=.*$")

    regExpressions := []regexp.Regexp{*authFailure,*notInSudoers,*command}
    _ = regExpressions

    // matches
    matchCommand := command.FindString(e.Message)
    if matchCommand != "" {
        user := regexp.MustCompile(`sudo:\s+\w+`)
        matchUser := user.FindString(e.Message)
        if matchUser == "" {
            log.Fatal(errors.CreateMatchError(*user,"message"))
        }
        e.Message = fmt.Sprintf("command: %s, user: %s\n", matchCommand,matchUser)
        e.Type = "Priviledge Escalation"
        events <- e
    }
}
