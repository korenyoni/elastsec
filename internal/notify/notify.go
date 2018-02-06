package notify

import (
    "../event"
    "os"
    "log"
    "fmt"
    "github.com/ashwanthkumar/slack-go-webhook"
)

func SendSlack(e event.Event, title string) {
    webhookUrl := os.Getenv("ESEC_SLACK_WEBHOOK")
    if webhookUrl == "" {
        log.Fatal("Error: No ESEC_SLACK_WEBHOOK set")
    }

    attachment1 := slack.Attachment {}
    attachment1.AddField(slack.Field { Title: "Content", Value: e.Message })
    payload := slack.Payload {
      Text: title,
      Username: "SwiftCop",
      Channel: "#platform",
      IconEmoji: ":cop:",
      Attachments: []slack.Attachment{attachment1},
    }
    err := slack.Send(webhookUrl, "", payload)
    if len(err) > 0 {
      fmt.Printf("error: %s\n", err)
    }
}
