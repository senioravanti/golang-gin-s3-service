package helpers

import (
	"fmt"
	"os"
)

func Getenv(
	key, fallback string,
) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func ErrWithCause(
	msg string,
	err error,
) error {
	return fmt.Errorf("%s | cause: `%w`", msg, err)
}