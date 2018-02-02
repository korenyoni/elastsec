package file_integrity

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

    fsFilter := elastic.NewTermQuery("metricset.name","file")
    go looper.Loop(eventBus, indexName, fsFilter)

    for event := range eventBus {

        var source map[string]interface{}
        json.Unmarshal(*event.Source, &source)
        audit := source["audit"]
        auditPretty,_ := json.MarshalIndent(audit,"","  ")
        event.Message = fmt.Sprintf("%s",auditPretty)
        event.Type = "File Integrity Change"
        events <- event
    }
}
