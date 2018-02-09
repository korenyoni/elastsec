package errors

import (
    "../constants"
    "fmt"
)

type ConnectionError struct {
    Err error
}

func (e *ConnectionError) Error() string {
    return fmt.Sprintf("Error: cannot connect to ElasticSearch. Is the env var %s set in the format 'http(s)://[address][port]'?\n%s\n",constants.ElasticAddressEnv,e.Err)
}

func CreateConnectionError(err error) error {
    return &ConnectionError{Err: err}
}
