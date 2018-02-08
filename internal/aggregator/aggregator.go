package aggregator

import (
    "log"
    "fmt"
    "time"
    "regexp"
    "strings"
    "encoding/json"
    "../env"
    "../event"
    "../constants"
)

type Aggregator struct {
    SupressedCount map[Key]*Info
}

type Info struct {
    Count int `json:"count"`
    Things []string `json:"things"`
}

type Key struct {
    /* Obligatory */
    Type string `json:"type"`
    Host string `json:"host"`
    /* Optional */
    User string `json:"user"`
}

type Instance struct {
    Key Key `json:"key"`
    Info Info `json:"info"`
    Env string `json:"env"`
}

func (a Aggregator) Loop(events chan<- event.Event, window time.Duration) {
    timer := time.NewTimer(window)
    defer timer.Stop()

    for c := time.Tick(window);; <- c {
        for k,i := range a.SupressedCount {
            instance := Instance{Key:k,Info:*i,Env:env.GetEnvName()}
            js, err := json.MarshalIndent(&instance, "","\t")
            if err != nil {
                log.Fatal("Error parsing aggregator event data")
            }
            if i.Count > 0 {
                e := event.Event{
                    Beat: event.Beat{
                        Host: k.Host},
                    Type: constants.AggregationEvent,
                    Message: fmt.Sprintf("supression of instance(s):\n%s",
                    js)}
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
    key,thing := genKeyThing(e)
    info, ok := a.SupressedCount[key]
    if ok {
        count := (*info).Count
        count = count + 1
        (*info).Count = count
        if !thingsContain((*info).Things, thing) {
            (*info).Things = append((*info).Things,thing)
        }
    } else {
        a.SupressedCount[key] = &Info{Count: 0,Things:[]string{thing}}
        return e, true
    }
    _ = thing
    return e, false
}

func genKeyThing(e event.Event) (Key, string) {
    var k Key
    var thing string
    k.Type = e.Type
    k.Host = e.Beat.Host
    quoteEscape := `\"`
    if e.Type == constants.FailedFileAccess || e.Type == constants.FileIntegrityChange {
        userRegex := regexp.MustCompile(`"user":\s+"\w+"`)
        userMatch := userRegex.FindString(e.Message)
        splitRegex := regexp.MustCompile(`:\s+`)
        userSplitMatch := splitRegex.FindString(userMatch)
        pathRegex := regexp.MustCompile(`"path":\s+".*"`)
        pathMatch := pathRegex.FindString(e.Message)
        pathSplitMatch := splitRegex.FindString(pathMatch)
        if userSplitMatch != "" {
            k.User = splitRegex.Split(userMatch,2)[1]
            k.User = strings.Trim(k.User,quoteEscape)
        }
        if pathSplitMatch != "" {
            thing = splitRegex.Split(pathMatch,2)[1]
            thing = strings.Trim(thing,quoteEscape)
        }
        return k,thing
    } else if e.Type == constants.AuthFailure {
        userRegex := regexp.MustCompile(`user=\w+`)
        userMatch := userRegex.FindString(e.Message)
        splitRegex := regexp.MustCompile(`=`)
        userSplitMatch := splitRegex.FindString(userMatch)
        if userSplitMatch != "" {
            k.User = splitRegex.Split(userMatch,2)[1]
        }
    } else if e.Type == constants.SSHAcceptedConnection || e.Type == constants.SSHFailedPass {
        userRegex := regexp.MustCompile(`for\s+\w+`)
        userMatch := userRegex.FindString(e.Message)
        splitRegex := regexp.MustCompile(`\s`)
        userSplitMatch := splitRegex.FindString(userMatch)
        invalidUserRegex := regexp.MustCompile(`invalid\s+user`)
        invalidUserMatch := invalidUserRegex.FindString(e.Message)
        if userSplitMatch != "" {
            k.User = splitRegex.Split(userMatch,2)[1]
        }
        if invalidUserMatch != "" {
            userRegex := regexp.MustCompile(`user\s+\w+`)
            userMatch := userRegex.FindString(e.Message)
            if userSplitMatch != "" {
                k.User = splitRegex.Split(userMatch,2)[1]
            }
        }
    } else if e.Type == constants.SSHInvalidUser {
        userRegex := regexp.MustCompile(`user\s+\w+`)
        userMatch := userRegex.FindString(e.Message)
        splitRegex := regexp.MustCompile(`\s`)
        userSplitMatch := splitRegex.FindString(userMatch)
        if userSplitMatch != "" {
            k.User = splitRegex.Split(userMatch,2)[1]
        }
    }  else if e.Type == constants.PrivEscalation {
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
            thing = commandSplitRegex.Split(commandMatch,2)[1]
        }
    }
    return k,thing
}

func thingsContain(things []string, thing string) bool {
    for _,t := range things {
        if t == thing {
            return true
        }
    }
    return false
}
