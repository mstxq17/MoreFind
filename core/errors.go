package core

import "github.com/pkg/errors"

func NewError(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, message)
}
