package constants

const (
    SlackHookEnv = "ESEC_SLACK_WEBHOOK"
    ElasticAddressEnv = "ES_ADDR"
    SmtpSendAddress = "SMTP_SEND_ADDR"
    SSHAcceptedConnection = "Accepted SSH connection"
    SSHDisconnect = "SSH Disconnect"
    SSHFailedPass = "Failed SSH connection (invalid password)"
    SSHInvalidUser = "Failed SSH connection (invalid user)"
    AuthFailure = "Authentication Failure"
    NotSudoer = "Unauthorized sudo attempt"
    PrivEscalation = "Priviledge Escalation"
    FileIntegrityChange = "File Integrity Change"
    FailedFileAccess = "Failed File Access"
    AggregationEvent = "Aggregation Event"
    FileMetricSet = "audit.file"
    KernelMetricSet = "audit.kernel"
)
