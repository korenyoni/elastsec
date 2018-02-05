package filechange_attempt

import (
    "../../looper"
    "../../event"
    "../../infoexport"
    "github.com/olivere/elastic"
    "fmt"
)

const (
    indexName   = "auditbeat*"
)

func Loop(events chan<- event.Event) {
    eventBus := make(chan event.Event)

    filter := elastic.NewBoolQuery()
    filter.Must(elastic.NewTermQuery("metricset.name","kernel"))
    filter.Must(elastic.NewTermQuery("audit.kernel.result","fail"))
    filter.Must(elastic.NewTermQuery("audit.kernel.thing.what", "file"))
    go looper.Loop(eventBus, indexName, filter)

    for event := range eventBus {

        data := infoexport.GetFileEventData(event)
        event.Message = fmt.Sprintf("%s\n%s",event.Time,data)
        event.Type = "Failed File Change"
        events <- event
    }
}
