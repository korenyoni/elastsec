package env

import (
    "log"
    "os"
    "time"
    "strconv"
    "../constants"
)

func GetAggDuration() time.Duration {
    durationString := os.Getenv(constants.AggDurationEnv)
    if durationString != "" {
        duration, err := time.ParseDuration(durationString)
        if err != nil {
            log.Println("Invalid aggregator duration.")
        } else {
            return duration
        }
    }
    return time.Hour
}

func GetEmailDuration() time.Duration {
    durationString := os.Getenv(constants.EmailDurationEnv)
    if durationString != "" {
    duration, err := time.ParseDuration(durationString)
        if err != nil {
            log.Println("Invalid email duration.")
        } else {
            return duration
        }
    }
    return time.Hour
}

func GetElasticUrl() string {
    elasticUrl := os.Getenv(constants.ElasticAddressEnv)
    if elasticUrl == "" {
        log.Fatal("No elastic url set")
    }
    return elasticUrl
}

func GetElasticSniff() bool {
    boolStr := os.Getenv(constants.ElasticSniffEnv)
    if boolStr == "" {
        return false
    }
    b,err := strconv.ParseBool(boolStr)
    if err != nil {
        log.Fatal("Invalid ElasticSearch sniffer setting")
    }
    return b
}

func GetSlackWebhook() string {
    slackWebhook := os.Getenv(constants.SlackHookEnv)
    if slackWebhook == "" {
        log.Fatal("No slack webhook address set")
    }
    return slackWebhook
}

func GetSendEmailAddress() string {
    emailStr := os.Getenv(constants.SendAddressEnv)
    if emailStr == "" {
        log.Fatal("No send email set")
    }
    return emailStr
}
