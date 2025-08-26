package helpers

import (
	"fmt"
)

func ErrWithCause(
	msg string,
	err error,
) error {
	return fmt.Errorf("%s | cause: `%w`", msg, err)
}