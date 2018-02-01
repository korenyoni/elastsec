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
    filter.Must(elastic.NewTermQuery("audit.kernel.key", "access"))
    go looper.Loop(eventBus, indexName, filter)

    for event := range eventBus {
        replaceMessage(events, event)
    }
}

func replaceMessage(events chan<- event.Event, e event.Event) {
    var source map[string]interface{}
    json.Unmarshal(*e.Source, &source)
    audit := source["audit"]
    kernel := audit.(map[string]interface{})["kernel"]
    actor := kernel.(map[string]interface{})["actor"]
    //
    username := actor.(map[string]interface{})["primary"]
    paths := kernel.(map[string]interface{})["paths"]
    data := kernel.(map[string]interface{})["data"]
    thing := kernel.(map[string]interface{})["thing"]
    dir := data.(map[string]interface{})["cwd"]
    filename := thing.(map[string]interface{})["primary"]
    var file_action string
    if paths.([]interface{})[0].(map[string]interface{})["name"] == "NORMAL" {
        file_action = "open"
    } else {
        file_action = "create"
    }
    e.Message = fmt.Sprintf("%s tried to %s file `%s` in `%s`",
    username, file_action, filename, dir)
    e.Type = fmt.Sprintf("File %s attempt",file_action)
    events <- e
}
