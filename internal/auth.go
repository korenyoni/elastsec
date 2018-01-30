package auth

import (
    "os"
    "fmt"
    "context"
    "reflect"
    "strings"
    "time"
    "github.com/olivere/elastic"
)

type Event struct {
    Host    string    `json:"beat.hostname"`
    Message string    `json:"message"`
    Time    time.Time `json:"@timestamp"`
}


func Loop(events chan<- string) {
    ctx := context.Background()

    // create client
    client, err := elastic.NewClient(elastic.SetURL(os.Getenv("ES_ADDR")), elastic.SetSniff(false))
    if err != nil {
      panic(err)
    }
    defer client.Stop()

    // range for last 10 seconds
    for _ = range time.Tick(10 * time.Second) {
          query := elastic.NewRangeQuery("@timestamp").From("now-10s").To("now")
          searchResult, err := client.Search().
          Index("filebeat*").
          Query(query).   // specify the query
          Pretty(true).       // pretty print request and response JSON
          Do(ctx)             // execute
        if err != nil {
        // Handle error
          panic(err)
        }
        var e Event
        for _, item := range searchResult.Each(reflect.TypeOf(e)) {
            parseEvent(events, item.(Event))
        }
    }
}

func parseEvent(events chan<- string, e Event) {
    if strings.Contains(e.Message, "COMMAND") {
        events <- fmt.Sprintf("time: %s message: %s\n", e.Time, e.Message)
    }
}
