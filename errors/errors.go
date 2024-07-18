package errors

import (
	"errors"
	"fmt"
)

func baseError(errorName string, details string) error {
	errorMsg := fmt.Sprintf("%s: %s", errorName, details)
	return errors.New(errorMsg)
}

func IllegalCharError(details string) error {
	return baseError("Illegal Character", details)
}

func InvalidSyntaxError(details string) error {
	return baseError("Invalid Syntax", details)
}
