package file_integrity

import (
    "../../looper"
    "../../event"
    "../../infoexport"
    "../../constants"
    "github.com/olivere/elastic"
    "fmt"
)

const (
    indexName   = "auditbeat*"
)

func Loop(events chan<- event.Event) {
    eventBus := make(chan event.Event)

    fsFilter := elastic.NewTermQuery("metricset.name","file")
    go looper.Loop(eventBus, indexName, fsFilter)

    for event := range eventBus {

        data := infoexport.GetFileEventData(event)
        event.Message = fmt.Sprintf("%s\n%s",event.Time,data)
        event.Type = constants.FileIntegrityChange
        events <- event
    }
}
