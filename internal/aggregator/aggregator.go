package aggregator

import (
    "log"
    "fmt"
    "time"
    "regexp"
    "strings"
    "encoding/json"
    "../event"
    "../constants"
)

type Aggregator struct {
    SupressedCount map[Key]int
}

type Key struct {
    /* Obligatory */
    Type string
    Host string
    /* Optional */
    User string
    Thing string
}

func (a Aggregator) Loop(events chan<- event.Event, window time.Duration) {
    timer := time.NewTimer(window)
    defer timer.Stop()

    for c := time.Tick(window);; <- c {
        for k,v := range a.SupressedCount {
            keyJS, err := json.MarshalIndent(&k, "  ","")
            if err != nil {
                log.Fatal("Error parsing aggregator event data")
            }
            if v > 0 {
                e := event.Event{
                    Beat: event.Beat{
                        Host: k.Host},
                    Type: constants.AggregationEvent,
                    Message: fmt.Sprintf("supression of `%d` instance(s):\n%s",
                    v,keyJS)}
                events <- e
            }
            delete(a.SupressedCount,k)
        }
    }
}

func (a Aggregator) Consume(e event.Event) (event.Event, bool) {
    if e.Type == constants.AggregationEvent {
        return e, true
    }
    key := genKey(e)
    count, ok := a.SupressedCount[key]
    if ok {
        count = count + 1
        a.SupressedCount[key] = count
    } else {
        a.SupressedCount[key] = 0
        return e, true
    }
    return e, false
}

func genKey(e event.Event) Key {
    var k Key
    k.Type = e.Type
    k.Host = e.Beat.Host
    quoteEscape := `\"`
    if e.Type == constants.FailedFileAccess || e.Type == constants.FileIntegrityChange {
        userRegex := regexp.MustCompile(`"user":\s+"\w+"`)
        userMatch := userRegex.FindString(e.Message)
        splitRegex := regexp.MustCompile(`:\s+`)
        userSplitMatch := splitRegex.FindString(userMatch)
        howRegex := regexp.MustCompile(`"how":\s+".*"`)
        howMatch := howRegex.FindString(e.Message)
        howSplitMatch := splitRegex.FindString(howMatch)
        if userSplitMatch != "" {
            k.User = splitRegex.Split(userMatch,2)[1]
            k.User = strings.Trim(k.User,quoteEscape)
        }
        if howSplitMatch != "" {
            k.Thing = splitRegex.Split(howMatch,2)[1]
            k.Thing = strings.Trim(k.Thing,quoteEscape)
        }
        return k
    } else if e.Type == constants.AuthFailure {
        userRegex := regexp.MustCompile(`user=\w+`)
        userMatch := userRegex.FindString(e.Message)
        splitRegex := regexp.MustCompile(`=`)
        userSplitMatch := splitRegex.FindString(userMatch)
        if userSplitMatch != "" {
            k.User = splitRegex.Split(userMatch,2)[1]
        }
    } else if e.Type == constants.PrivEscalation {
        userRegex := regexp.MustCompile(`sudo:\s+\w+`)
        userMatch := userRegex.FindString(e.Message)
        userSplitRegex := regexp.MustCompile(`\s+`)
        userSplitMatch := userSplitRegex.FindString(userMatch)
        commandRegex := regexp.MustCompile(`COMMAND=.*$`)
        commandMatch := commandRegex.FindString(e.Message)
        commandSplitRegex := regexp.MustCompile(`=`)
        commandSplitMatch := userSplitRegex.FindString(commandMatch)
        if userSplitMatch != "" {
            k.User = userSplitRegex.Split(userMatch,2)[1]
        }
        if commandSplitMatch != "" {
            k.Thing = commandSplitRegex.Split(commandMatch,2)[1]
        }
    }
    return k
}
