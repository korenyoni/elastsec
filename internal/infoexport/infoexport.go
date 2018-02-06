package infoexport

import (
    "reflect"
    "fmt"
    "log"
    "../event"
    "../evalpath"
    "../constants"
    "encoding/json"
    "github.com/tidwall/gjson"
)

func GetTitle(e event.Event) string {
    if e.Type != constants.AggregationEvent {
        title := fmt.Sprintf("New `%s` event on host `%s`",
        e.Type, e.Beat.Host)
        return title
    }
    return "Previous suppresion of events:"
}

func GetFileEventData(e event.Event) string {
    jsonData := *e.Source

    metricset := constants.KernelMetricSet
    var data map[string]interface{}
    if gjson.GetBytes(jsonData,constants.KernelMetricSet).String() == "" {
        metricset = constants.FileMetricSet
        owner := gjson.GetBytes(jsonData,metricset + ".owner")
        group := gjson.GetBytes(jsonData,metricset + ".group")
        path := gjson.GetBytes(jsonData,metricset + ".path")
        action := gjson.GetBytes(jsonData,metricset + ".action")
        mode := gjson.GetBytes(jsonData,metricset + ".mode")

        data = map[string]interface{}{ "owner": owner.String(),
        "group": group.String(),
        "path": path.String(),
        "action": action.String(),
        "mode": mode.String()}

    } else {
        user := gjson.GetBytes(jsonData,metricset + ".actor.attrs.uid")
        action := gjson.GetBytes(jsonData,metricset + ".action")
        cwd := gjson.GetBytes(jsonData,metricset + ".data.cwd")
        parentPath := gjson.GetBytes(jsonData,metricset + ".paths.0.name")
        childPath := gjson.GetBytes(jsonData,metricset + ".paths.1.name")
        path := evalpath.Eval(cwd.String(),parentPath.String(), childPath.String())
        how :=  gjson.GetBytes(jsonData,metricset + ".how")

        data = map[string]interface{}{ "user": user.String(), "action": action.String(),
        "path": path, "how": how.String()}
    }

    for _,v := range reflect.ValueOf(data).MapKeys() {
        k := v.Interface().(string)
        if data[k] == "" {
            delete(data,k)
        }
    }

    jsonDataConcise, err := json.MarshalIndent(data,"  ","")
    if err != nil {
        log.Fatal("Error encoding FileEvent data")
    }
    return fmt.Sprintf("%s",jsonDataConcise)
}
