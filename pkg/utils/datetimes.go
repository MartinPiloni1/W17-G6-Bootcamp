package utils

import (
	"time"

	"github.com/go-playground/validator"
)

// intended use, add it as a new parameter in playground/validator
// validates that the field Time is not seted into the future
func NotFutureDatetime(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	now := time.Now()
	return !date.After(now)
}
