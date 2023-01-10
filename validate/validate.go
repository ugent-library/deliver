package validate

import (
	"fmt"
)

const (
	RuleNotEmpty = "not_empty"
	RuleLengthIs = "length_is"
	RuleLengthIn = "length_in"
)

var (
	MessageNotEmpty = "cannot be empty"
	MessageLengthIs = "length must be %d"
	MessageLengthIn = "length must be between %d and %d"
)

func NotEmpty[T ~string | ~[]any | ~map[any]any](key string, val T) *Error {
	if len(val) == 0 {
		return &Error{key: key, rule: RuleNotEmpty, msg: MessageNotEmpty}
	}
	return nil
}

func LengthIs[T ~string | ~[]any | ~map[any]any](key string, val T, n int) *Error {
	if len(val) != n {
		return &Error{
			key:    key,
			rule:   RuleLengthIs,
			msg:    fmt.Sprintf(MessageLengthIs, n),
			params: []any{n},
		}
	}
	return nil
}

func LengthIn[T ~string | ~[]any | ~map[any]any](key string, val T, min, max int) *Error {
	if len(val) < min || len(val) > max {
		return &Error{
			key:    key,
			rule:   RuleLengthIn,
			msg:    fmt.Sprintf(MessageLengthIn, min, max),
			params: []any{min, max},
		}
	}
	return nil
}
