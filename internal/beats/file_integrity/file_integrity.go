package file_integrity

import (
    "../../looper"
    "../../event"
    "github.com/olivere/elastic"
    "github.com/tidwall/gjson"
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

        jsonData := *event.Source
        user := gjson.GetBytes(jsonData,"actor")
        data := fmt.Sprintf("%s\n",user.String())
        event.Message = fmt.Sprintf("%s\n%s",event.Time,data)
        event.Type = "File Integrity Change"
        events <- event
    }
}
