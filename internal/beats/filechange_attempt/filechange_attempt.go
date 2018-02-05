package filechange_attempt

import (
    "../../looper"
    "../../event"
    "../../evalpath"
    "github.com/olivere/elastic"
    "github.com/tidwall/gjson"
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

        jsonData := *event.Source
        user := gjson.GetBytes(jsonData,"audit.kernel.actor.attrs.uid")
        action := gjson.GetBytes(jsonData,"audit.kernel.action")
        cwd := gjson.GetBytes(jsonData,"audit.kernel.data.cwd")
        parentPath := gjson.GetBytes(jsonData,"audit.kernel.paths.0.name")
        childPath := gjson.GetBytes(jsonData,"audit.kernel.paths.1.name")
        path := evalpath.Eval(cwd.String(),parentPath.String(), childPath.String())
        data := fmt.Sprintf("%s\n%s\n%s\n",user.String(),action.String(),path)
        event.Message = fmt.Sprintf("%s\n%s",event.Time,data)
        event.Type = "Failed File Change"
        events <- event
    }
}
