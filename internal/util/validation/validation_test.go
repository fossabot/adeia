package validation

import (
	"errors"
	"strings"
	"testing"

	"adeia-api/internal/util"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/stretchr/testify/assert"
)

func TestValidation_Validate(t *testing.T) {
	v := Validation{
		Errors{
			"B": errors.New("b1"),
			"C": nil,
			"A": errors.New("a1"),
		},
	}
	err := v.Validate()
	assert.Equal(t, 2, len(err.(util.ResponseError).ValidationErrors))
	if assert.NotNil(t, err) {
		assert.Equal(t,
			map[string]string{
				"A": "a1",
				"B": "b1",
			},
			err.(util.ResponseError).ValidationErrors,
		)
	}

	v = Validation{}
	assert.Nil(t, v.Validate())

	v = Validation{
		Errors{
			"B": nil,
			"C": nil,
		},
	}

	assert.Nil(t, v.Validate())
}

type testCase struct {
	in   string
	want error
	msg  string
}

func TestValidateName(t *testing.T) {
	// success cases
	testCases := []testCase{
		{"foo", nil, "should not return err on normal names"},
		{"田中太郎", nil, "should accept all kinds of unicode characters"},
		{"Guðmundsdóttir", nil, "should accept all kinds of unicode characters"},
		{"Ельцин", nil, "should accept all kinds of unicode characters"},
		{"Björk Guðmundsdóttir", nil, "should accept names with spaces"},
		{" Björk Guðmundsdóttir ", nil, "should accept names with spaces"},
		{"Jennifer 8 Lee", nil, "should accept names with numbers"},
		{"María-Jose", nil, "should accept names with special characters"},
	}
	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			got := ValidateName(tc.in)
			assert.Nil(t, got, tc.msg)
		})
	}

	// err cases
	errTestCases := []testCase{
		{
			strings.Repeat("f", 256),
			validation.ErrLengthOutOfRange.SetParams(map[string]interface{}{"min": 1, "max": 255}),
			"name too long",
		},
		{
			"",
			validation.ErrRequired,
			"name cannot be empty",
		},
	}
	for _, tc := range errTestCases {
		t.Run(tc.in, func(t *testing.T) {
			got := ValidateName(tc.in)
			if assert.NotNil(t, got) {
				assert.Equal(t, tc.want, got, tc.msg)
			}
		})
	}
}

func TestValidateDesignation(t *testing.T) {
	// success cases
	testCases := []testCase{
		{"foo", nil, "should not return err on normal designations"},
		{"foo bar", nil, "designations can contain spaces"},
		{"Administrator Level 2", nil, "designations can contain numbers"},
		{"Administrator-Lvl 2 (junior & intermediate)", nil, "designations can contain special characters .+-/_&()[]{}"},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			got := ValidateDesignation(tc.in)
			assert.Nil(t, got, tc.msg)
		})
	}

	// err cases
	errTestCases := []testCase{
		{
			strings.Repeat("f", 256),
			validation.ErrLengthOutOfRange.SetParams(map[string]interface{}{"min": 2, "max": 255}),
			"designation too long",
		},
		{
			"a",
			validation.ErrLengthOutOfRange.SetParams(map[string]interface{}{"min": 2, "max": 255}),
			"designation too short",
		},
		{"", validation.ErrRequired, "designation cannot be empty"},
		{"Administrator * Level 2", is.ErrAlphanumeric, "designation cannot contain special characters other than .+-/_&()[]{}"},
		{"田中太郎", is.ErrAlphanumeric, "designation cannot contain other language characters"},
	}
	for _, tc := range errTestCases {
		t.Run(tc.in, func(t *testing.T) {
			got := ValidateDesignation(tc.in)
			if assert.NotNil(t, got) {
				assert.Equal(t, tc.want, got, tc.msg)
			}
		})
	}
}
