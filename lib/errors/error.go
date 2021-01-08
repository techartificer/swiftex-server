package errors

import (
	"fmt"
	"strings"
)

type UndefinedError struct {
	Err string `json:"message"`
}

func NewError(msg string) *UndefinedError {
	return &UndefinedError{
		Err: msg,
	}
}

func (uf *UndefinedError) Error() string {
	return fmt.Sprintf("%s", uf.Err)
}

func IsMongoDupError(err error) bool {
	isDup := strings.Contains(err.Error(), "E11000")
	return isDup
}
