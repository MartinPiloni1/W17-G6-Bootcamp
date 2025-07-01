package utils

import (
	"errors"
	"reflect"
)

func ApplyNonZero[T any](dst *T, patch T) error {
	if dst == nil {
		return errors.New("dst no puede ser nil")
	}

	dstVal := reflect.ValueOf(dst).Elem() // *T → T
	patchVal := reflect.ValueOf(patch)

	for patchVal.Kind() == reflect.Pointer {
		if patchVal.IsNil() {
			return errors.New("patch es puntero nil")
		}
		patchVal = patchVal.Elem()
	}

	if dstVal.Type() != patchVal.Type() || dstVal.Kind() != reflect.Struct {
		return errors.New("dst y patch deben ser structs del mismo tipo")
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
