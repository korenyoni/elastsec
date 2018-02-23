package notify

import (
    "../event"
    "../env"
    "../constants"
    "log"
    "fmt"
    "time"
    "os/exec"
    "github.com/ashwanthkumar/slack-go-webhook"
)

type Email struct {
    Msg []byte
    SendAddress string
    Window time.Duration
    Host string
}

func SendSlack(e event.Event, title string) {
    webhookUrl := env.GetSlackWebhook()

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

func EmailInit(events chan event.Event, window time.Duration) *Email {
    sendAddress := env.GetSendEmailAddress()

    msg := make([]byte,0)

    return &Email{Msg:msg,SendAddress:sendAddress,Window:window}
}

func (em *Email) Consume(e event.Event, title string) {
    em.Host = e.Beat.Host
    eventMessage := []byte(fmt.Sprintf("%s\n%s\n\n",title,e.Message))
    em.Msg = append(em.Msg,eventMessage...)
    if isUrgent(e) {
        em.sendMailUrgent(e,title)
    }
}

func (em *Email) Loop() {
    for range time.Tick(em.Window) {
        if len(em.Msg) > 0 {
            fmt.Println("Sending Email...")
            err := em.sendMail()
            if err != nil {
                log.Fatal(fmt.Sprintf("Error sending email: %s",err))
            }
        }
        em.Msg = make([]byte,0)
    }
}

func (em *Email) sendMail() error {
    address := em.SendAddress
    cmd := exec.Command(constants.SendMailCommand, fmt.Sprintf("-s Elastsec security notifications on host %s, env %s", em.Host, env.GetEnvName()),address)
    stdin,err := cmd.StdinPipe()
    if err != nil {
        panic(err)
    }
    stdin.Write([]byte(em.Msg))
    stdin.Close()
    err = cmd.Start()
    return err
}

func (em *Email) sendMailUrgent(e event.Event, title string) error {
    address := em.SendAddress
    cmd := exec.Command(constants.SendMailCommand, fmt.Sprintf("-s Urgent Elastsec security notification on host %s, env %s", em.Host, env.GetEnvName()),address)
    stdin,err := cmd.StdinPipe()
    if err != nil {
        panic(err)
    }
    stdin.Write([]byte(fmt.Sprintf("%s\n%s\n",title,e.Message)))
    stdin.Close()
    err = cmd.Start()
    return err
}

func isUrgent(e event.Event) bool {
    urgentRegex := env.GetUrgentRegex()
    envName := env.GetEnvName()
    if urgentRegex.FindString(envName) != "" {
        return true && e.Type != constants.AggregationEvent
    }
    return false
}
