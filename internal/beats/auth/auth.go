package auth

import (
    "../../looper"
    "../../event"
    "../../constants"
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
    acceptedPassword := regexp.MustCompile("Accepted")
    acceptedPublickey := regexp.MustCompile("Accepted publickey")
    disconnect := regexp.MustCompile("disconnected by")
    failedPassword := regexp.MustCompile("Failed password")
    failedConnection := regexp.MustCompile("Connection from")
    invalidUser := regexp.MustCompile("Invalid user.*from")
    authFailure := regexp.MustCompile("authentication failure")
    notInSudoers := regexp.MustCompile("NOT in sudoers")
    command := regexp.MustCompile("COMMAND=.*$")

    // matches
    matchSsh := ssh.FindString(e.Message)
    matchAcceptedPassword := acceptedPassword.FindString(e.Message)
    matchAcceptedPublickey := acceptedPublickey.FindString(e.Message)
    matchDisconnect := disconnect.FindString(e.Message)
    matchFailedPassword := failedPassword.FindString(e.Message)
    matchFailedConnection := failedConnection.FindString(e.Message)
    matchInvalidUser := invalidUser.FindString(e.Message)
    matchAuthFailure := authFailure.FindString(e.Message)
    matchNotInSudoers := notInSudoers.FindString(e.Message)
    matchCommand := command.FindString(e.Message)

    if matchSsh != "" && (matchAcceptedPassword != "" || matchAcceptedPublickey != "") {
        e.Type = constants.SSHAcceptedConnection
        events <- e
    } else if matchSsh != "" && matchDisconnect != "" {
        e.Type = constants.SSHDisconnect
        events <- e
    } else if matchSsh != "" && (matchFailedPassword != "" || (matchFailedConnection != "" && matchAcceptedPublickey == "")) {
        e.Type = constants.SSHFailedAuth
        events <- e
    } else if matchSsh != "" && matchInvalidUser != "" {
        e.Type = constants.SSHInvalidUser
        events <- e
    } else if matchAuthFailure != "" && matchSsh == "" {
        e.Type = constants.AuthFailure
        events <- e
    } else if matchNotInSudoers != "" {
        e.Type = constants.NotSudoer
        events <- e
    } else if matchCommand != "" {
        e.Type = constants.PrivEscalation
        events <- e
    }
}
