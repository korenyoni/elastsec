package constants

const (
    SlackHookEnv = "ESEC_SLACK_WEBHOOK"
    ElasticAddressEnv = "ES_ADDR"
    SSHAcceptedConnection = "Accepted SSH connection"
    SSHDisconnect = "SSH Disconnect"
    SSHFailedPass = "Failed SSH connection (invalid password)"
    SSHInvalidUser = "Failed SSH connection (invalid user)"
    AuthFailure = "Authentication Failure"
    NotSudoer = "Unauthorized sudo attempt"
    PrivEscalation = "Priviledge Escalation"
    FileIntegrityChange = "File Integrity Change"
    FailedFileChange = "Failed File Change"
    AggregationEvent = "Aggregation Event"
)
