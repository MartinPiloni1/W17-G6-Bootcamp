package httperrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetErrorData(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Case: Bad Request error",
			err:          BadRequestError{Message: "bad request"},
			expectedCode: 400,
			expectedMsg:  "bad request",
		},
		{
			name:         "Case: Not Found error",
			err:          NotFoundError{Message: "not found"},
			expectedCode: 404,
			expectedMsg:  "not found",
		},
		{
			name:         "Case: Conflict error",
			err:          ConflictError{Message: "conflict"},
			expectedCode: 409,
			expectedMsg:  "conflict",
		},
		{
			name:         "Case: Unprocessable Entity error",
			err:          UnprocessableEntityError{Message: "unprocessable"},
			expectedCode: 422,
			expectedMsg:  "unprocessable",
		},
		{
			name:         "Case: Internal server error",
			err:          InternalServerError{Message: "Internal Server Error"},
			expectedCode: 500,
			expectedMsg:  "Internal Server Error",
		},
		{
			name:         "Case: Handle of external errors",
			err:          errors.New("some random error"),
			expectedCode: 500,
			expectedMsg:  "Internal Server Error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualCode, actualMsg := GetErrorData(tc.err)
			require.Equal(t, tc.expectedCode, actualCode)
			require.Equal(t, tc.expectedMsg, actualMsg)
		})
	}
}
