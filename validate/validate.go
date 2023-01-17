package validate

import (
	"fmt"
	"regexp"
)

const (
	RuleNotEmpty     = "not_empty"
	RuleLength       = "length"
	RuleLengthIn     = "length_in"
	RuleMin          = "min"
	RuleMax          = "max"
	RuleMatch        = "match"
	RuleAlphanumeric = "alphanumeric"
)

var (
	MessageNotEmpty     = "cannot be empty"
	MessageLength       = "length must be %d"
	MessageLengthIn     = "length must be between %d and %d"
	MessageMin          = "must be %d or more"
	MessageMax          = "must be %d or less"
	MessageMatch        = "must match %s"
	MessageAlphanumeric = "must only contain letters a to z and digits"

	ReAlphanumeric = regexp.MustCompile("^[a-zA-Z0-9]+$")
)

func Validate(errs ...*Error) error {
	return NewErrors(errs...).ErrorOrNil()
}

// TODO return *Errors object from rules
func NotEmpty[T ~string | ~[]any | ~map[any]any](key string, val T) *Error {
	if len(val) == 0 {
		return &Error{
			key:  key,
			rule: RuleNotEmpty,
			msg:  MessageNotEmpty,
		}
	}
	return nil
}

func Length[T ~string | ~[]any | ~map[any]any](key string, val T, n int) *Error {
	if len(val) != n {
		return &Error{
			key:    key,
			rule:   RuleLength,
			msg:    fmt.Sprintf(MessageLength, n),
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

func Min[T int | int64 | float64](key string, val T, min T) *Error {
	if val < min {
		return &Error{
			key:    key,
			rule:   RuleMin,
			msg:    fmt.Sprintf(MessageMin, min),
			params: []any{min},
		}
	}
	return nil
}

func Max[T int | int64 | float64](key string, val T, max T) *Error {
	if val > max {
		return &Error{
			key:    key,
			rule:   RuleMax,
			msg:    fmt.Sprintf(MessageMax, max),
			params: []any{max},
		}
	}
	return nil
}

func Match(key, val string, r *regexp.Regexp) *Error {
	if !r.MatchString(val) {
		return &Error{
			key:    key,
			rule:   RuleNotEmpty,
			msg:    fmt.Sprintf(MessageMatch, r),
			params: []any{r},
		}
	}
	return nil
}

func Alphanumeric(key, val string) *Error {
	if !ReAlphanumeric.MatchString(val) {
		return &Error{
			key:  key,
			rule: RuleAlphanumeric,
			msg:  MessageAlphanumeric,
		}
	}
	return nil
}
