package utils

import (
	"errors"
	"reflect"
)

func ApplyNonZero[T any](dst *T, patch T) error {
	if dst == nil {
		return errors.New("the original struct cannot be nil")
	}

	dstVal := reflect.ValueOf(dst).Elem() // *T → T
	patchVal := reflect.ValueOf(patch)

	for patchVal.Kind() == reflect.Pointer {
		if patchVal.IsNil() {
			return errors.New("the patch struct is a nil pointer")
		}
		patchVal = patchVal.Elem()
	}

	if dstVal.Type() != patchVal.Type() || dstVal.Kind() != reflect.Struct {
		return errors.New("the original struct and the patch struct must be structs of the same type")
	}

	for i := 0; i < dstVal.NumField(); i++ {
		fDst := dstVal.Field(i)
		fSrc := patchVal.Field(i)

		if !fDst.CanSet() {
			continue
		}

		if fSrc.IsZero() {
			continue
		}

		fDst.Set(fSrc)
	}

	return nil
}
