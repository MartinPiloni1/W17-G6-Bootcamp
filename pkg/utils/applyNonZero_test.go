package utils_test

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

type patchStruct struct {
	A int
	B string
	C float64
	d bool
}

func TestApplyNonZero_Basic(t *testing.T) {
	dst := &patchStruct{A: 1, B: "hello", C: 4.2}
	patch := patchStruct{A: 0, B: "bye", C: 0}
	err := utils.ApplyNonZero(dst, patch)
	require.NoError(t, err)
	require.Equal(t, &patchStruct{A: 1, B: "bye", C: 4.2}, dst)
}

func TestApplyNonZero_PrivateField_NoChange(t *testing.T) {
	type S struct {
		X int
		y int // no se debe tocar
	}
	dst := &S{X: 1, y: 9}
	patch := S{X: 2, y: 99}
	err := utils.ApplyNonZero(dst, patch)
	require.NoError(t, err)
	require.Equal(t, 2, dst.X)
	require.Equal(t, 9, dst.y)
}

func TestApplyNonZero_NilDst(t *testing.T) {
	var dst *patchStruct
	patch := patchStruct{A: 3}
	err := utils.ApplyNonZero(dst, patch)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot be nil")
}

func TestApplyNonZero_DifferentTypes(t *testing.T) {
	dst := &patchStruct{}
	err := utils.ApplyNonZero(dst, patchStruct{}) // correct usage
	require.NoError(t, err)
}

func TestApplyNonZero_ZeroPatch(t *testing.T) {
	dst := &patchStruct{A: 7, B: "q"}
	patch := patchStruct{}
	err := utils.ApplyNonZero(dst, patch)
	require.NoError(t, err)
	require.Equal(t, &patchStruct{A: 7, B: "q"}, dst)
}

func TestApplyNonZero_DifferentType_NonStruct(t *testing.T) {
	// dst = int, patch = int
	i := 10
	err := utils.ApplyNonZero(&i, 15)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be structs of the same type")
}
