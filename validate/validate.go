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
	RuleUnique       = "unique"
)

var (
	MessageNotEmpty     = "cannot be empty"
	MessageLength       = "length must be %d"
	MessageLengthIn     = "length must be between %d and %d"
	MessageMin          = "must be %d or more"
	MessageMax          = "must be %d or less"
	MessageMatch        = "must match %s"
	MessageAlphanumeric = "must only contain letters a to z and digits"
	MessageUnique       = "must be unique"

	ReAlphanumeric = regexp.MustCompile("^[a-zA-Z0-9]+$")
)

func Validate(errs ...*Error) error {
	return NewErrors(errs...).ErrorOrNil()
}

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
			params: []any{n},
			msg:    fmt.Sprintf(MessageLength, n),
		}
	}
	return nil
}

func LengthIn[T ~string | ~[]any | ~map[any]any](key string, val T, min, max int) *Error {
	if len(val) < min || len(val) > max {
		return &Error{
			key:    key,
			rule:   RuleLengthIn,
			params: []any{min, max},
			msg:    fmt.Sprintf(MessageLengthIn, min, max),
		}
	}
	return nil
}

func Min[T int | int64 | float64](key string, val T, min T) *Error {
	if val < min {
		return &Error{
			key:    key,
			rule:   RuleMin,
			params: []any{min},
			msg:    fmt.Sprintf(MessageMin, min),
		}
	}
	return nil
}

func Max[T int | int64 | float64](key string, val T, max T) *Error {
	if val > max {
		return &Error{
			key:    key,
			rule:   RuleMax,
			params: []any{max},
			msg:    fmt.Sprintf(MessageMax, max),
		}
	}
	return nil
}

func Match(key, val string, r *regexp.Regexp) *Error {
	if !r.MatchString(val) {
		return &Error{
			key:    key,
			rule:   RuleNotEmpty,
			params: []any{r},
			msg:    fmt.Sprintf(MessageMatch, r),
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

func ErrNotUnique(key string) *Error {
	return &Error{
		key:  key,
		rule: RuleUnique,
		msg:  MessageUnique,
	}
}
