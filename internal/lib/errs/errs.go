package errs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/meowalien/RabbitGather-golang.git/internal/lib/caller"
)

type WithLineError interface {
	error
}

/*
New Usage:

	New(any...) => make a withLineError of errors.New(fmt.Sprint(any...))
	New(string , any...) => make a withLineError of fmt.Errorf(string, obj...)
	New(error1(could be nil) , error2/string , error3/string ...) => wrap error1(error2(error3 ...)))
	New(error) make a withLineError of error
*/
func New(err any, obj ...any) WithLineError {
	ee := newWithLineErrorFromAny(true, caller.Caller(1, caller.CALLER_FORMAT_SHORT), err, obj...)
	// to make sure that the returned error is nil type and nil value
	if ee == nil {
		return nil
	}
	return ee
}

func newWithLineErrorFromAny(deliverMode bool, caller string, err any, obj ...any) *withLineError {
	if err == nil || err == (*withLineError)(nil) {
		if len(obj) == 0 {
			return nil
		} else if len(obj) == 1 {
			return newWithLineErrorFromAny(deliverMode, caller, obj[0])
		} else {
			return newWithLineErrorFromAny(deliverMode, caller, obj[0], obj[1:]...)
		}
	}
	switch errTp := err.(type) {
	case string:
		if strings.Contains(errTp, "%") {
			return newWithLineErrorFromError(fmt.Errorf(errTp, obj...), caller)
		} else {
			return newWithLineErrorFromError(errors.New(fmt.Sprint(append([]any{errTp}, obj...)...)), caller)
		}
	case error:
		var parentErr *withLineError
		parentErr, ok := errTp.(*withLineError) //nolint:errorlint
		if ok {
			if deliverMode {
				parentErr = parentErr.deliver(caller)
			}
		} else {
			parentErr = newWithLineErrorFromError(errTp, caller)
		}

		if len(obj) == 0 {
			return parentErr
		}

		var toWrapErr error

		switch obj0 := obj[0].(type) {
		case error:
			toWrapErr = errors.New(fmt.Sprint(obj...))
			return parentErr.wrap(toWrapErr, "")

		case string:
			if strings.Contains(obj0, "%") {
				toWrapErr = fmt.Errorf(obj0, obj[1:]...)
			} else {
				toWrapErr = errors.New(fmt.Sprint(obj...))
			}

		}

		return parentErr.wrap(toWrapErr, caller)
	default:
		return newWithLineErrorFromError(errors.New(fmt.Sprint(append([]any{errTp}, obj...)...)), caller)
	}
}

func newWithLineErrorFromError(err error, caller string) *withLineError {
	return &withLineError{caller: caller, error: err}
}

type withLineError struct {
	error
	parent *withLineError
	caller string
	layer  int
}

func (w *withLineError) Unwrap() error {
	if w.parent == nil {
		return w.error
	} else {
		return w.parent
	}
}

func (w *withLineError) wrap(a error, caller string) (res *withLineError) {
	ne := newWithLineErrorFromAny(false, caller, a)
	ne.layer = w.layer + 1
	ne.parent = w
	return ne
}

func (w withLineError) deliver(caller string) *withLineError {
	w.caller = fmt.Sprintf("%s <= %s", w.caller, caller)
	x := w
	return &x
}
