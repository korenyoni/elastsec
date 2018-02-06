package looper

import (
    "../errors"
    "../event"
    "../constants"
    "log"
    "os"
    "encoding/json"
    "context"
    "reflect"
    "time"
    "strings"
    "github.com/olivere/elastic"
)

func Loop(events chan<- event.Event, indexName string, filter elastic.Query) {
    ctx := context.Background()

    // create client
    client, err := elastic.NewClient(elastic.SetURL(os.Getenv(constants.ElasticAddressEnv)), elastic.SetSniff(false))
    if err != nil {
        if strings.Contains(err.Error(), "no Elasticsearch node available") {
            log.Fatal(errors.CreateConnectionError())
        }
      panic(err)
    }
    defer client.Stop()

    var e event.Event
    // Get latest
    query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery())
    if filter != nil {
        query = query.Filter(filter)
    }
    searchResult, err := client.Search(indexName).
    Index(indexName).
    Query(query).
    Sort("@timestamp", false).
    Size(1).
    Do(ctx)

    lastItem := searchResult.Each(reflect.TypeOf(e))[0]
    // range for last 10 seconds
    for c := time.Tick(10 * time.Second);; <- c {
          query := elastic.NewBoolQuery().Must(
              elastic.NewRangeQuery("@timestamp").
              From(lastItem.(event.Event).Time.Add(time.Millisecond)).
              To("now"))
          if filter != nil {
              query = query.Filter(filter)
          }
          searchResult, err = client.Search().
          Index(indexName).
          Query(query).   // specify the query
          Sort("@timestamp",true).
          Pretty(true).       // pretty print request and response JSON
          Do(ctx)             // execute
        if err != nil {
        // Handle error
          panic(err)
        }
        array := searchResult.Hits.Hits
        for _, hit := range array {
            json.Unmarshal(*hit.Source, &e)
            e.Source = hit.Source
            events <- e
            lastItem = e
        }
    }
}

