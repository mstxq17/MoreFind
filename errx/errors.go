package errx

import (
	"github.com/pkg/errors"
)

// NewWrapError 自定义错误处理函数
func NewWrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, message)

}

func NewWithMsgf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return errors.WithMessagef(err, format, args...)
}

func NewMsgf(format string, args ...any) error {
	return errors.Errorf(format, args...)
}

func NewMsg(message string) error {
	return errors.New(message)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
