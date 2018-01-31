package kernel

import (
    "../../looper"
    "../../event"
    "github.com/olivere/elastic"
)

const (
    indexName   = "auditbeat*"
)

func Loop(events chan<- event.Event) {
    eventBus := make(chan event.Event)

    go looper.Loop(eventBus, indexName, elastic.NewTermQuery("metricset.name","kernel"))

    for event := range eventBus {
        events <- event
    }
}
