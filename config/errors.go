package config

import "fmt"

// ErrParsingTimeout indicates that something went wrong while parsing a timeout field
type ErrParsingTimeout struct {
	Value string
	Err   error
}

func (e ErrParsingTimeout) Error() string {
	return fmt.Sprintf("parsing timeout value %s: %s", e.Value, e.Err)
}
