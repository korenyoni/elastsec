package errors

import (
    "../constants"
    "fmt"
)

type ConnectionError struct {
}

func (e *ConnectionError) Error() string {
    return fmt.Sprintf("Error: cannot connect to ElasticSearch. Is the env var %s set in the format 'http(s)://[address][port]'?\n",constants.ElasticAddressEnv)
}

func CreateConnectionError() error {
    return &ConnectionError{}
}
