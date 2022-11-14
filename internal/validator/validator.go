// Filename: internal/validator/validator.go
package validator

import (
	"regexp"
)

var (
	EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// we create a type that wraps our validation errors map
type Validator struct {
	Errors map[string]string
}

//new creates a new validator instance
func New() *Validator{
	return &Validator{
		Errors: make(map[string]string),
	}
}

//valid() checks the erros map for entires
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

//IN() check if an element can be found in a provided list of elements
func In(element string, list ...string) bool {
	for i := range list{
		if element == list[i]{
			return true
		}
	}
	return false
}

//matches() returns true if a string value matches a specific regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

//adderror() aadds an error entry to Error map
func (v *Validator) AddError (key, message string) {
	if _, exists := v.Errors[key]; !exists{
		v.Errors[key] = message
	}
}

//check() performs the validation checks and calls the adderrors
//method in turn if an error entry needs to be added
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
	
}

//unique() check that there are no repeating values in the slice
func Unique(values []string) bool {
	UniqueValues := make(map[string]bool)
	for _,value := range values{
		UniqueValues[value] = true
	}
	return len(values) == len(UniqueValues)
}
