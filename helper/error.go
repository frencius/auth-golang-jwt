package helper

import (
	"errors"
	"strings"
)

func ErrStringsToErr(errStrs []string) error {
	if len(errStrs) > 0 {
		return errors.New(strings.Join(errStrs, ", "))
	}

	return nil
}
