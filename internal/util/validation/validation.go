package validation

import (
	"errors"
	"time"

	"adeia-api/internal/util"
	"adeia-api/internal/util/crypto"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Errors represents a map of validation errors containing, fieldName:error pairs.
type Errors map[string]error

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

// Validation represents a map of validation rules.
type Validation struct {
	Errors
}

// Validate calls the underlying library's Filter() and returns an
// ErrValidationFailed when validation fails.
func (v *Validation) Validate() error {
	err := validation.Errors(v.Errors).Filter()
	if err == nil {
		return nil
	}

	e := make(map[string]string)
	for k, v := range err.(validation.Errors) {
		e[k] = v.Error()
	}
	return util.ErrValidationFailed.ValidationErr(e)
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

// ==========
// Policies for common fields
// ==========

// NamePolicy is the validation policy for user's name.
var NamePolicy = []validation.Rule{
	validation.Required,
	validation.RuneLength(1, 255),
	// we should not have any other checks here, because names can contain anything
	// See: https://www.kalzumeus.com/2010/06/17/falsehoods-programmers-believe-about-names/
	// and https://shinesolutions.com/2018/01/08/falsehoods-programmers-believe-about-names-with-examples/
}

// DesignationPolicy is the validation policy for user's designation.
var DesignationPolicy = []validation.Rule{
	validation.Required,
	validation.RuneLength(2, 255),
	WithWhiteAndSpecialChars(".+-/_&()[]{}", is.Alphanumeric),
}

// EmailPolicy is the validation policy for user's email.
var EmailPolicy = []validation.Rule{
	validation.Required,
	validation.RuneLength(3, 120),
	is.EmailFormat,
}

var empIDPolicyCore = []validation.Rule{
	validation.RuneLength(5, 10),
	is.Alphanumeric,
}

// EmpIDPolicy is the validation policy for user's employee ID.
var EmpIDPolicy = append([]validation.Rule{validation.Required}, empIDPolicyCore...)

// EmpIDPolicyOptional is a special EmpIDPolicy that validates employee ID only if
// it's not empty (meaning, employee ID is optional). This is used in CreateUser controller
// handle.
func EmpIDPolicyOptional(value interface{}) []validation.Rule {
	return append([]validation.Rule{validation.Required.When(value != "")}, empIDPolicyCore...)
}

// PwdPolicy is the validation policy for user's password.
var PwdPolicy = []validation.Rule{
	validation.Required,
	validation.RuneLength(12, 128),
}

// PwdStrengthPolicy is the validation policy for the strength user's password. This
// in combination with PwdPolicy, is used during user-onboarding.
var PwdStrengthPolicy = []validation.Rule{
	validation.By(func(value interface{}) error {
		s, _ := value.(string)
		if crypto.PasswordStrength(s) < 3 {
			return errors.New("password is weak")
		}
		return nil
	}),
}

// ConfirmPwdPolicy is the validation policy for the confirm-password field. This policy
// just checks whether the value is equal to the password. DO NOT use this separately, because
// it does not validate other aspects of a password. This must be used in combination
// with the PwdPolicy.
func ConfirmPwdPolicy(password string) []validation.Rule {
	return []validation.Rule{
		validation.By(func(value interface{}) error {
			s, _ := value.(string)
			if s != password {
				return errors.New("passwords do not match")
			}
			return nil
		}),
	}
}

// ResourceNamePolicy is the validation policy for names of resources, like name
// of holiday, etc.
var ResourceNamePolicy = []validation.Rule{
	validation.Required,
	validation.RuneLength(2, 255),
	WithWhiteAndSpecialChars(".+-/_&()[]{}", is.Alphanumeric),
}

// ResourceIDPolicy is the validation policy for ids of resources. This is used
// to validate db's auto-increment ids.
var ResourceIDPolicy = []validation.Rule{
	validation.Required,
	is.Digit,
	validation.Min(1),
}

// DatePolicy is the validation policy for a date string.
var DatePolicy = []validation.Rule{
	validation.Required,
	validation.Date(time.RFC3339),
}

// YearPolicy is the validation policy for a year value. A year must contain exactly
// four digits.
var YearPolicy = []validation.Rule{
	validation.Required,
	is.Digit,
	validation.Min(4),
	validation.Max(4),
}

// MonthPolicy is the validation policy for a month value. It must only be between 1-12.
var MonthPolicy = []validation.Rule{
	validation.Required,
	is.Digit,
	validation.Min(1),
	validation.Max(12),
}

// DayPolicy is the validation policy for a day of month value. It must only be between 1-31.
var DayPolicy = []validation.Rule{
	validation.Required,
	is.Digit,
	validation.Min(1),
	validation.Max(31),
}

// TokenPolicy is the validation policy for secure random tokens.
var TokenPolicy = []validation.Rule{
	validation.Required,
}

// ==========
// Validation funcs
// ==========

// ValidateName validates user names.
func ValidateName(value string) error {
	return validation.Validate(value, NamePolicy...)
}

// ValidateDesignation validates user designations.
func ValidateDesignation(value string) error {
	return validation.Validate(value, DesignationPolicy...)
}

// ValidateEmail validates user email IDs.
func ValidateEmail(value string) error {
	return validation.Validate(value, EmailPolicy...)
}

// ValidateLoginPwd validates the password entered during login (without strength
// estimation).
func ValidateLoginPwd(value string) error {
	return validation.Validate(value, PwdPolicy...)
}

// ValidateNewPwd validates the password entered during user-onboarding (with strength
// estimation).
func ValidateNewPwd(value string) error {
	policy := append(PwdPolicy, PwdStrengthPolicy...)
	return validation.Validate(value, policy...)
}

// ValidateConfirmPwd validates the confirm-password entered during user-onboarding. The value
// must be equal to the password. No strength estimation is done here.
func ValidateConfirmPwd(value string, password string) error {
	policy := append(PwdPolicy, ConfirmPwdPolicy(password)...)
	return validation.Validate(value, policy...)
}

// ValidateEmpID validates user employee IDs.
func ValidateEmpID(value string) error {
	return validation.Validate(value, EmpIDPolicy...)
}

// ValidateEmpIDOptional validates the optional employee ID entered during user-onboarding.
func ValidateEmpIDOptional(value interface{}) error {
	return validation.Validate(value, EmpIDPolicyOptional(value)...)
}

// ValidateCustom is used to write custom validation rules for one-off use-cases.
func ValidateCustom(value interface{}, rules ...validation.Rule) error {
	return validation.Validate(value, rules...)
}

// ValidateResourceName validates resource names.
func ValidateResourceName(value string) error {
	return validation.Validate(value, ResourceNamePolicy...)
}

// ValidateResourceID validates resource ids.
func ValidateResourceID(value string) error {
	return validation.Validate(value, ResourceIDPolicy...)
}

// ValidateDate validates a date string.
func ValidateDate(value string) error {
	return validation.Validate(value, DatePolicy...)
}

// ValidateYear validates a year value (eg., 2012).
func ValidateYear(value string) error {
	return validation.Validate(value, YearPolicy...)
}

// ValidateMonth validates a month value (eg., 4).
func ValidateMonth(value string) error {
	return validation.Validate(value, MonthPolicy...)
}

// ValidateDay validates a day of month value (eg., 24).
func ValidateDay(value string) error {
	return validation.Validate(value, DayPolicy...)
}

// ValidateToken validates secure random tokens.
func ValidateToken(value string) error {
	return validation.Validate(value, TokenPolicy...)
}
