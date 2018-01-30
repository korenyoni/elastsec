package auth

import (
    "os"
    "fmt"
    "context"
    "reflect"
    "regexp"
    "time"
    "github.com/olivere/elastic"
)

type Event struct {
    Host    string    `json:"beat.hostname"`
    Message string    `json:"message"`
    Time    time.Time `json:"@timestamp"`
}

const (
    indexName   = "filebeat*"
)

func Loop(events chan<- string) {
    ctx := context.Background()

    // create client
    client, err := elastic.NewClient(elastic.SetURL(os.Getenv("ES_ADDR")), elastic.SetSniff(false))
    if err != nil {
      panic(err)
    }
    defer client.Stop()

    var e Event
    // range for last 10 seconds
    for c := time.Tick(10 * time.Second);; <- c {
          // Get latest
          searchResult, err := client.Search(indexName).
          Index(indexName).
          Query(elastic.NewMatchAllQuery()).
          Sort("@timestamp", false).
          Size(1).
          Do(ctx)

          last := searchResult.Each(reflect.TypeOf(e))[0]
          events <- fmt.Sprintf("latest %s", last.(Event).Time)
          lastTime := last.(Event).Time
          lastTime = lastTime.Add(time.Millisecond)

          query := elastic.NewRangeQuery("@timestamp").From(lastTime).To("now")
          searchResult, err = client.Search().
          Index(indexName).
          Query(query).   // specify the query
          Pretty(true).       // pretty print request and response JSON
          Do(ctx)             // execute
        if err != nil {
        // Handle error
          panic(err)
        }
        for _, item := range searchResult.Each(reflect.TypeOf(e)) {
            parseEvent(events, item.(Event))
        }
    }
}

func parseEvent(events chan<- string, e Event) {
    r := regexp.MustCompile("COMMAND=.*$")
    match := r.FindString(e.Message)
    if match != "" {
        events <- fmt.Sprintf("time: %s command: %s\n", e.Time, match)
    }
}
