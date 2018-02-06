package notify

import (
    "../event"
    "../constants"
    "os"
    "log"
    "fmt"
    "github.com/ashwanthkumar/slack-go-webhook"
)

func SendSlack(e event.Event, title string) {
    webhookUrl := os.Getenv(constants.SlackHookEnv)
    if webhookUrl == "" {
        log.Fatal(fmt.Sprintf("Error: No %s set",constants.SlackHookEnv))
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
