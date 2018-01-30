package auth

import (
    "os"
    "fmt"
    "context"
    "time"
    "github.com/olivere/elastic"
)

func Loop(events chan<- string) {
    ctx := context.Background()

    // create client
    client, err := elastic.NewClient(elastic.SetURL(os.Getenv("ES_ADDR")), elastic.SetSniff(false))
    if err != nil {
      panic(err)
    }
    defer client.Stop()

    // range for last 10 seconds
    for x := range time.Tick(10 * time.Second) {
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
        events <- fmt.Sprintf("%d hits, ", searchResult.TotalHits()) + x.String()
    }
}

