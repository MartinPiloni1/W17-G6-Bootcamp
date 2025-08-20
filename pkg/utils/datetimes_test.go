package utils_test

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/require"
)

func TestNotFutureDatetime(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.RegisterValidation("not_future", utils.NotFutureDatetime))
	type S struct {
		Date time.Time `validate:"not_future"`
	}

	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		Name    string
		Date    time.Time
		WantErr bool
	}{
		{"now", now, false},
		{"past", past, false},
		{"future", future, true},
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

func TestNotFutureDatetime_NotTimeType(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.RegisterValidation("not_future", utils.NotFutureDatetime))
	type S struct {
		Fake int `validate:"not_future"`
	}
	s := S{Fake: 9}
	err := v.Struct(s)
	require.Error(t, err)
}
