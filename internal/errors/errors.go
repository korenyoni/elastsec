package errors

type ConnectionError struct {
}

func (e *ConnectionError) Error() string {
    return "Error: cannot connect to ElasticSearch. Is the env var ES_ADDR set in the format 'http(s)://[address][port]'?\n"
}

func CreateConnectionError() error {
    return &ConnectionError{}
}
