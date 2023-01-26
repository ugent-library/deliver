package validate

import (
	"fmt"
	"strings"
)

type Errors struct {
	errors []*Error
}

func NewErrors(errs ...*Error) *Errors {
	e := &Errors{}
	return e.Add(errs...)
}

func (e *Errors) Error() string {
	msg := ""
	for i, err := range e.errors {
		msg += err.Error()
		if i < len(e.errors)-1 {
			msg += ", "
		}
	}
	return msg
}

func (e *Errors) Add(errs ...*Error) *Errors {
	for _, err := range errs {
		if err != nil {
			e.errors = append(e.errors, err)
		}
	}
	return e
}

func (e *Errors) AddWithPrefix(prefix string, errs ...*Error) *Errors {
	for _, err := range errs {
		if err != nil {
			err.key = prefix + err.key
			e.errors = append(e.errors, err)
		}
	}
	return e
}

func (e *Errors) Get(key string) *Error {
	for _, e := range e.errors {
		if e.key == key {
			return e
		}
	}
	return nil
}

func (e *Errors) WithPrefix(prefix string) *Errors {
	ee := &Errors{}

	for _, err := range e.errors {
		if strings.HasPrefix(err.key, prefix) {
			ee.errors = append(ee.errors, err)
		}
	}

	return ee
}

func (e *Errors) ErrorOrNil() error {
	if len(e.errors) > 0 {
		return e
	}
	return nil
}

type Error struct {
	key    string
	rule   string
	params []any
	msg    string
}

func NewError(key, rule string, params ...any) *Error {
	return &Error{
		key:    key,
		rule:   rule,
		params: params,
	}
}

func (e *Error) Key() string {
	return e.key
}

func (e *Error) Rule() string {
	return e.rule
}

func (e *Error) Message() string {
	return e.msg
}

func (e *Error) WithMessage(msg string) *Error {
	if e != nil {
		e.msg = msg
	}
	return e
}

func (e *Error) Error() string {
	msg := e.key
	if msg != "" {
		msg += " "
	}
	if e.msg != "" {
		msg += e.msg
	} else if e.rule != "" {
		msg += e.rule
		if len(e.params) > 0 {
			msg += "["
			for i, p := range e.params {
				msg += fmt.Sprintf("%v", p)
				if i < len(e.params)-1 {
					msg += ", "
				}
			}
			msg += "]"
		}
	}
	return msg
}
