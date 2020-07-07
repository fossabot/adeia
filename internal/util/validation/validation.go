package validation

import (
	"errors"

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

func WithWhitespace(rules ...validation.Rule) validation.Rule {
	return &WrappedRule{
		rules:           rules,
		allowWhitespace: true,
	}
}

func WithSpecialChars(chars string, rules ...validation.Rule) validation.Rule {
	return &WrappedRule{
		rules:               rules,
		allowSpecialChars:   true,
		allowedSpecialChars: chars,
	}
}

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

var NamePolicy = []validation.Rule{
	validation.Required,
	validation.RuneLength(1, 255),
	// we should not have any other checks here, because names can contain anything
	// See: https://www.kalzumeus.com/2010/06/17/falsehoods-programmers-believe-about-names/
	// and https://shinesolutions.com/2018/01/08/falsehoods-programmers-believe-about-names-with-examples/
}

var DesignationPolicy = []validation.Rule{
	validation.Required,
	validation.RuneLength(2, 255),
	WithWhiteAndSpecialChars(".+-/_&()[]{}", is.Alphanumeric),
}

var EmailPolicy = []validation.Rule{
	validation.Required,
	validation.RuneLength(3, 120),
	is.EmailFormat,
}

var empIDPolicyCore = []validation.Rule{
	validation.RuneLength(5, 10),
	is.Alphanumeric,
}

var EmpIDPolicy = append([]validation.Rule{validation.Required}, empIDPolicyCore...)

func EmpIDPolicyOptional(value interface{}) []validation.Rule {
	return append([]validation.Rule{validation.Required.When(value != "")}, empIDPolicyCore...)
}

var PwdPolicy = []validation.Rule{
	validation.Required,
	validation.RuneLength(12, 128),
}

var PwdStrengthPolicy = []validation.Rule{
	validation.By(func(value interface{}) error {
		s, _ := value.(string)
		if crypto.PasswordStrength(s) < 3 {
			return errors.New("password is weak")
		}
		return nil
	}),
}

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

// ==========
// Validation funcs
// ==========

func ValidateName(value string) error {
	return validation.Validate(value, NamePolicy...)
}

func ValidateDesignation(value string) error {
	return validation.Validate(value, DesignationPolicy...)
}

func ValidateEmail(value string) error {
	return validation.Validate(value, EmailPolicy...)
}

func ValidateLoginPwd(value string) error {
	return validation.Validate(value, PwdPolicy...)
}

func ValidateNewPwd(value string) error {
	policy := append(PwdPolicy, PwdStrengthPolicy...)
	return validation.Validate(value, policy...)
}

func ValidateConfirmPwd(value string, password string) error {
	policy := append(PwdPolicy, ConfirmPwdPolicy(password)...)
	return validation.Validate(value, policy...)
}

func ValidateEmpID(value string) error {
	return validation.Validate(value, EmpIDPolicy...)
}

func ValidateEmpIDOptional(value interface{}) error {
	return validation.Validate(value, EmpIDPolicyOptional(value)...)
}

func ValidateCustom(value interface{}, rules ...validation.Rule) error {
	return validation.Validate(value, rules...)
}
