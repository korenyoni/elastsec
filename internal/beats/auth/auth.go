package auth

import (
    "../../looper"
    "../../event"
    "regexp"
)

const (
    indexName   = "filebeat*"
)

func Loop(events chan<- event.Event) {
    eventBus := make(chan event.Event)

    go looper.Loop(eventBus, indexName, nil)

    for event := range eventBus {
        replaceMessage(events, event)
    }
}

func replaceMessage(events chan<- event.Event, e event.Event) {
    ssh := regexp.MustCompile("sshd")
    acceptedPassword := regexp.MustCompile("Accepted password")
    disconnect := regexp.MustCompile("session closed")
    failedPassword := regexp.MustCompile("Failed password")
    invalidUser := regexp.MustCompile("Invalid user.*from")
    authFailure := regexp.MustCompile("authentication failure")
    notInSudoers := regexp.MustCompile("NOT in sudoers")
    command := regexp.MustCompile("COMMAND=.*$")

    // matches
    matchSsh := ssh.FindString(e.Message)
    matchAcceptedPassword := acceptedPassword.FindString(e.Message)
    matchDisconnect := disconnect.FindString(e.Message)
    matchFailedPassword := failedPassword.FindString(e.Message)
    matchInvalidUser := invalidUser.FindString(e.Message)
    matchAuthFailure := authFailure.FindString(e.Message)
    matchNotInSudoers := notInSudoers.FindString(e.Message)
    matchCommand := command.FindString(e.Message)

    if matchSsh != "" && matchAcceptedPassword != "" {
        e.Type = "Accepted SSH connection"
        events <- e
    } else if matchSsh != "" && matchDisconnect != "" {
        e.Type = "SSH Disconnect"
        events <- e
    } else if matchSsh != "" && matchFailedPassword != "" {
        e.Type = "Failed SSH connection (invalid password)"
        events <- e
    } else if matchSsh != "" && matchInvalidUser != "" {
        e.Type = "Failed SSH connection (invalid user)"
        events <- e
    } else if matchAuthFailure != "" && matchSsh == "" {
        e.Type = "Authentication Failure"
        events <- e
    } else if matchNotInSudoers != "" {
        e.Type = "Unauthorized sudo attempt"
        events <- e
    } else if matchCommand != "" {
        e.Type = "Priviledge Escalation"
        events <- e
    }
}
