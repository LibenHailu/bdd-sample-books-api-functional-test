// packege errors contains errors types mthod and funtions to hide implementation details from the user
// when errors occur, rather than sending an error occured we send a custom error to not expose impleentaion details
package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// type Err struct {
// 	Code    string `json:"code"`
// 	Message string `json:"message"`
// }

// //ErrorEnvelope to give more context for an error when error occurs
// //	to send out put json like this
// // {
// // 	error:{
// // 		message:"invalid input"
// // 	}
// // }
// // The consumer knows if he is accessing the correct data or not in the above case
// //
// // unenveloped error will look like this
// // {
// //	"message":"invalid input"
// // }

// type ErrEnvelope map[string]Err

// // NewEnvelope creates a new ErrEnvelope
// func NewEnvelope(envStr string, err Err) ErrEnvelope {
// 	return ErrEnvelope{envStr: err}
// }

// // type validationError map[string]map[string]string

// // // // NewValidationErrorEnvelope creates a new validationErrorEnvelope
// // // func NewValidationErrorEnvelope(envStr string, err map[string]map[string]string) validationError {
// // // 	return map[string]map[string]map[string]string{envStr: err}
// // // }

// NewUserValidationError creates a new map of string to string of errors
func ValidationError(valError error) map[string]map[string]string {

	errs := make(map[string]string)

	var verr validator.ValidationErrors

	if errors.As(valError, &verr) {
		for _, f := range verr {
			err := f.ActualTag()

			if err == "required" {
				err = fmt.Sprintf("%s is required %s", f.Field(), err)
			}
			errs[strings.ToLower(f.Field())] = err
		}
	}
	return map[string]map[string]string{"errors": errs}
}

func ErrorMap(msg string) map[string]map[string]string {
	return map[string]map[string]string{
		"errors": map[string]string{"error": msg},
	}
}
