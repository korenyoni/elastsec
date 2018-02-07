package notify

import (
    "../event"
    "../constants"
    "os"
    "os/exec"
    "log"
    "fmt"
    "time"
    "github.com/ashwanthkumar/slack-go-webhook"
)

type Email struct {
    Msg []byte
    SendAddress string
    Window time.Duration
}

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

func EmailInit(events chan event.Event, window time.Duration) *Email {
    smtpSendAddress := os.Getenv(constants.SmtpSendAddress)
    checkEmailEnv([]string{smtpSendAddress,constants.SmtpSendAddress})

    msg := make([]byte,0)

    return &Email{Msg:msg,SendAddress:smtpSendAddress,Window:window}
}

func (em *Email) Consume(e event.Event, title string) {
    eventMessage := []byte(fmt.Sprintf("%s\n%s\n",title,e.Message))
    em.Msg = append(em.Msg,eventMessage...)
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

func checkEmailEnv(e ...[]string) {
    for _,envVar := range e {
        if envVar[0] == "" {
            log.Fatal(fmt.Sprintf("Error: %s not set",envVar[1]))
        }
    }
}

func (em *Email) sendMail() error {
    address := em.SendAddress
    cmd := exec.Command("sendmail", address)
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return err
    }
    stdin.Write([]byte(fmt.Sprintf("To: %s\n", address)))
    stdin.Write([]byte(fmt.Sprintf("Subject: %s\n", "Elastsec security notifications")))
    stdin.Write(em.Msg)
    err = cmd.Start()
    return err
}
