package filechange_attempt

import (
    "../../looper"
    "../../event"
    "github.com/olivere/elastic"
    "encoding/json"
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

        var source map[string]interface{}
        json.Unmarshal(*event.Source, &source)
        audit := source["audit"]
        auditPretty,_ := json.MarshalIndent(audit,"","  ")
        event.Message = fmt.Sprintf("%s",auditPretty)
        event.Type = "Failed File Change"
        events <- event
    }
}
