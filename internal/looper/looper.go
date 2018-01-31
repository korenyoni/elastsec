package looper

import (
    "../errors"
    "log"
    "os"
    "encoding/json"
    "context"
    "reflect"
    "time"
    "strings"
    "github.com/olivere/elastic"
)

type Event struct {
    Host    string    `json:"beat.hostname"`
    Message string    `json:"message"`
    Time    time.Time `json:"@timestamp"`
}

func Loop(events chan<- Event, indexName string) {
    ctx := context.Background()

    // create client
    client, err := elastic.NewClient(elastic.SetURL(os.Getenv("ES_ADDR")), elastic.SetSniff(false))
    if err != nil {
        if strings.Contains(err.Error(), "no Elasticsearch node available") {
            log.Fatal(errors.CreateConnectionError())
        }
      panic(err)
    }
    defer client.Stop()

    var e Event
    // Get latest
    searchResult, err := client.Search(indexName).
    Index(indexName).
    Query(elastic.NewMatchAllQuery()).
    Sort("@timestamp", false).
    Size(1).
    Do(ctx)

    lastItem := searchResult.Each(reflect.TypeOf(e))[0]
    // range for last 10 seconds
    for c := time.Tick(10 * time.Second);; <- c {
          query := elastic.NewRangeQuery("@timestamp").From(lastItem.(Event).Time.Add(time.Millisecond)).To("now")
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
            events <- e
            lastItem = e
        }
    }
}

