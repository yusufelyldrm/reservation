package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

// Form creates a custom from struct which embeds an url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters", length))
		return false
	}
	return true

}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if form field is in post is not empty
func (f *Form) Has(field string) bool {

	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// IsEmail checks if email is valid
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
