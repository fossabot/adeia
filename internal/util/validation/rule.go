package validation

import (
	"strconv"

	"adeia-api/internal/util"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// WrappedRule is a custom validation rule that wraps existing rules. Currently, WrappedRule
// is used to allow certain characters in strings, like whitespaces and special characters,
// before actually executing the wrapped rule. For example, with WrappedRule, a custom validator
// can be created that allows only "a space-separated UTFLetter" string.
type WrappedRule struct {
	rules               []validation.Rule
	allowWhitespace     bool
	allowSpecialChars   bool
	allowedSpecialChars string
}

// Validate performs the actual validation. This is basically the Validate method of
// validation.StringRule, with plugged-in logic to allow certain additional characters that
// existing rules won't allow.
func (w *WrappedRule) Validate(value interface{}) error {
	// copied over from StringRule
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}
	str, err := validation.EnsureString(value)
	if err != nil {
		return err
	}

	switch {
	case w.allowWhitespace:
		str = util.StripWhitespace(str)
		fallthrough

	case w.allowSpecialChars:
		str = util.StripChars(str, w.allowedSpecialChars)
		fallthrough

	default:
		break
	}

	// run the actual rule
	return validation.Validate(str, w.rules...)
}

// StringThresholdRule is a validation rule that checks if a numeric-string value
// satisfies the specified threshold requirement. This is adapted from validation.ThresholdRule.
type StringThresholdRule struct {
	threshold int64
	operator  int
	err       validation.Error
}

const (
	greaterThan = iota
	greaterEqualThan
	lessThan
	lessEqualThan
)

// Validate checks if the given value is valid or not.
func (r StringThresholdRule) Validate(value interface{}) error {
	// copied over from StringRule
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}
	str, err := validation.EnsureString(value)
	if err != nil {
		return err
	}

	// convert to number
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return is.ErrDigit
	}

	// compare value
	if r.compareInt(r.threshold, num) {
		return nil
	}
	return r.err.SetParams(map[string]interface{}{"threshold": r.threshold})
}

func (r StringThresholdRule) compareInt(threshold, value int64) bool {
	switch r.operator {
	case greaterThan:
		return value > threshold
	case greaterEqualThan:
		return value >= threshold
	case lessThan:
		return value < threshold
	default:
		return value <= threshold
	}
}

// ==========
// Custom rules
// ==========

// Min returns a validation rule that checks if a numeric-string value is greater or
// equal than the specified value. By calling Exclusive,the rule will check if the value
// is strictly greater than the specified value. An empty value is considered valid.
// Please use the Required rule to make sure a value is not empty.
func Min(min int64) StringThresholdRule {
	return StringThresholdRule{
		threshold: min,
		operator:  greaterEqualThan,
		err:       validation.ErrMinGreaterEqualThanRequired,
	}
}

// Max returns a validation rule that checks if a numeric-string value is less or
// equal than the specified value. By calling Exclusive,the rule will check if the value
// is strictly less than the specified value. An empty value is considered valid.
// Please use the Required rule to make sure a value is not empty.
func Max(max int64) StringThresholdRule {
	return StringThresholdRule{
		threshold: max,
		operator:  lessEqualThan,
		err:       validation.ErrMaxLessEqualThanRequired,
	}
}

// Exclusive sets the comparison to exclude the boundary value.
func (r StringThresholdRule) Exclusive() StringThresholdRule {
	if r.operator == greaterEqualThan {
		r.operator = greaterThan
		r.err = validation.ErrMinGreaterThanRequired
	} else if r.operator == lessEqualThan {
		r.operator = lessThan
		r.err = validation.ErrMaxLessThanRequired
	}
	return r
}

// WithWhitespace is a wrapped rule that allows the use of whitespaces along with
// the wrapped rule. Basically, it strips of any whitespace before passing it to
// the wrapped rule for validation.
func WithWhitespace(rules ...validation.Rule) validation.Rule {
	return &WrappedRule{
		rules:           rules,
		allowWhitespace: true,
	}
}

// WithSpecialChars is a wrapped rule that allows the use of specified chars along with
// the wrapped rule. Basically, it strips of any of the specified chars before passing it to
// the wrapped rule for validation.
func WithSpecialChars(chars string, rules ...validation.Rule) validation.Rule {
	return &WrappedRule{
		rules:               rules,
		allowSpecialChars:   true,
		allowedSpecialChars: chars,
	}
}

// WithWhiteAndSpecialChars is a wrapped rule that allows the use of whitespaces and specified
// chars along with the wrapped rule. Basically, it strips of any of the specified chars
// before passing it to the wrapped rule for validation.
func WithWhiteAndSpecialChars(chars string, rules ...validation.Rule) validation.Rule {
	return &WrappedRule{
		rules:               rules,
		allowWhitespace:     true,
		allowSpecialChars:   true,
		allowedSpecialChars: chars,
	}
}
