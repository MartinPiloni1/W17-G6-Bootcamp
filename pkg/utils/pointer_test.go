package utils_test

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPtr(t *testing.T) {
	n := 7
	p := utils.Ptr(n)
	require.Equal(t, &n, p)

	s := "hola"
	ps := utils.Ptr(s)
	require.Equal(t, &s, ps)

	// Nil value
	var x *int
	require.Nil(t, x)
}
