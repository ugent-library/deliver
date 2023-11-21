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
		return NewError(key, RuleNotEmpty).WithMessage(MessageNotEmpty)
	}
	return nil
}

func Length[T ~string | ~[]any | ~map[any]any](key string, val T, n int) *Error {
	if len(val) != n {
		return NewError(key, RuleLength, n).WithMessage(fmt.Sprintf(MessageLength, n))
	}
	return nil
}

func LengthIn[T ~string | ~[]any | ~map[any]any](key string, val T, min, max int) *Error {
	if len(val) < min || len(val) > max {
		return NewError(key, RuleLengthIn, min, max).WithMessage(fmt.Sprintf(MessageLengthIn, min, max))
	}
	return nil
}

func Min[T int | int64 | float64](key string, val T, min T) *Error {
	if val < min {
		return NewError(key, RuleMin, min).WithMessage(fmt.Sprintf(MessageMin, min))
	}
	return nil
}

func Max[T int | int64 | float64](key string, val T, max T) *Error {
	if val > max {
		return NewError(key, RuleMax, max).WithMessage(fmt.Sprintf(MessageMax, max))
	}
	return nil
}

func Match(key, val string, r *regexp.Regexp) *Error {
	if !r.MatchString(val) {
		return NewError(key, RuleMatch, r).WithMessage(fmt.Sprintf(MessageMatch, r))
	}
	return nil
}

func Alphanumeric(key, val string) *Error {
	if !ReAlphanumeric.MatchString(val) {
		return NewError(key, RuleAlphanumeric).WithMessage(MessageAlphanumeric)
	}
	return nil
}

func ErrNotUnique(key string) *Error {
	return NewError(key, RuleUnique).WithMessage(MessageUnique)
}
