package errors

import "fmt"

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
