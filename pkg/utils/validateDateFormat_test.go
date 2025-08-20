package utils

import (
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/require"
)

func TestValidateDateFormat(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.RegisterValidation("dateformat", ValidateDateFormat))
	type S struct {
		Date string `validate:"dateformat"`
	}
	tests := []struct {
		Name    string
		Date    string
		WantErr bool
	}{
		{"valida", "2024-01-20", false},
		{"bad-day", "2024-02-31", true},
		{"empty", "", true},
		{"slash", "2000/10/10", true},
		{"text", "hola mundo", true},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			s := S{Date: tc.Date}
			err := v.Struct(s)
			if tc.WantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
