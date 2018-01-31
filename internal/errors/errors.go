package errors

import (
    "regexp"
    "fmt"
)

type MatchError struct {
    Message string
    Field string
}

func (e *MatchError) Error() string {
    return fmt.Sprintf("Error: cannot match: '%s' in field '%s'\n", e.Message, e.Field)
}

type ConnectionError struct {
}

func (e *ConnectionError) Error() string {
    return "Error: cannot connect to ElasticSearch. Is the env var ES_ADDR set in the format 'http(s)://[address][port]'?\n"
}


func CreateMatchError(r regexp.Regexp, field string) error {
    return &MatchError{Message: r.String(), Field: field}
}

func CreateConnectionError() error {
    return &ConnectionError{}
}
