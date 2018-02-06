package aggregator

import (
    "fmt"
    "time"
    "../event"
    "../constants"
)

type Aggregator struct {
    SupressedCount map[string]int
}

func (a Aggregator) Loop(events chan<- event.Event, window time.Duration) {
    timer := time.NewTimer(window)
    defer timer.Stop()

    for c := time.Tick(window);; <- c {
        for k,v := range a.SupressedCount {
            e := event.Event{
                Type: constants.AggregationEvent,
                Message: fmt.Sprintf("`%d` supressed instance(s) of `%s`",v,k)}
            events <- e
            delete(a.SupressedCount,k)
        }
    }
}

func (a Aggregator) Consume(e event.Event) (event.Event, bool) {
    count, ok := a.SupressedCount[e.Type]
    if e.Type == constants.AggregationEvent {
        return e, true
    }
    if ok {
        count = count + 1
        a.SupressedCount[e.Type] = count
    } else {
        a.SupressedCount[e.Type] = 0
        return e, true
    }
    return e, false
}
