package auth

import (
    "../looper"
    "../errors"
    "../event"
    "fmt"
    "log"
    "strings"
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

    // matches
    matchAuthFailure := authFailure.FindString(e.Message)
    matchNotInSudoers := notInSudoers.FindString(e.Message)
    matchCommand := command.FindString(e.Message)
    if matchAuthFailure != "" {
        e.Type = "Authentication Failure"
        user := regexp.MustCompile(`user=\w+`)
        matchUser := strings.Trim(user.FindString(e.Message),"user=")
        if matchUser == "" {
            log.Fatal(errors.CreateMatchError(*user,"message"))
        }
        e.Message = fmt.Sprintf("user `%s` failed to authenticate.\n",matchUser)
        events <- e
    } else if matchNotInSudoers != "" {
        e.Type = "Unauthorized sudo attempt"
        user := regexp.MustCompile(`sudo:\s+\w+`)
        matchUser := strings.Split(user.FindString(e.Message),"sudo:")[1]
        matchUser = strings.Trim(matchUser," ")
        if matchUser == "" {
            log.Fatal(errors.CreateMatchError(*user,"message"))
        }
        e.Message = fmt.Sprintf("user `%s` attempted to use sudo, but isn't in sudoers.\n",matchUser)
        events <- e
    } else if matchCommand != "" {
        e.Type = "Priviledge Escalation"
        user := regexp.MustCompile(`sudo:\s+\w+`)
        matchUser := strings.Split(user.FindString(e.Message),"sudo:")[1]
        matchUser = strings.Trim(matchUser," ")
        matchCommand = strings.Trim(matchCommand, "COMMAND=")
        if matchUser == "" {
            log.Fatal(errors.CreateMatchError(*user,"message"))
        }
        e.Message = fmt.Sprintf("user `%s` executed command: `%s`\n", matchUser, matchCommand)
        events <- e
    }
}
