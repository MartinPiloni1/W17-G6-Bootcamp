package utils

import (
	"time"

	"github.com/go-playground/validator"
)

// ValidateDateFormat validates the date format YYYY-MM-DD
func ValidateDateFormat(fl validator.FieldLevel) bool {
	const layout = "2006-01-02"
	dateStr := fl.Field().String()
	_, err := time.Parse(layout, dateStr)
	return err == nil
}