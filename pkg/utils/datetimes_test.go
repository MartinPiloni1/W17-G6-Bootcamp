package utils_test

import (
	"testing"
	"time"

	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Date time.Time `validate:"not_future"`
}

func TestNotFutureDatetimeValidator(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("not_future", utils.NotFutureDatetime)
	require.NoError(t, err)

	tests := []struct {
		name    string
		date    time.Time
		wantErr bool
	}{
		{
			"CurrentTime",
			time.Now(),
			false,
		},
		{
			"PastTime",
			time.Now().Add(-1 * time.Hour),
			false,
		},
		{
			"FutureTime",
			time.Now().Add(1 * time.Hour),
			true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// arrange
			s := TestStruct{Date: testCase.date}
			err := v.Struct(s)

			// act
			gotErr := err != nil

			// assert
			assert.Equal(t, testCase.wantErr, gotErr)
		})
	}
}
