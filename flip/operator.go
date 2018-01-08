package flip

import (
	"reflect"
	"errors"
	"fmt"
)

func StringFlipOperator(flipValue *string, flipValueRegistered interface{}) (*bool, error) {
	switch reflect.TypeOf(flipValueRegistered).Kind() {
	case reflect.Slice:
			return rollOutOfString(flipValue, flipValueRegistered)
	case reflect.String:
		return isFlipString(flipValue, flipValueRegistered)
	default:
		return flipResponse(false), errors.New(fmt.Sprintf("unknown type, data_type=%s", reflect.TypeOf(flipValueRegistered).Kind().String()))

	}
}

func IntFlipOperator(flipValue *int64, flipValueRegistered interface{}) (*bool, error) {
	switch reflect.TypeOf(flipValueRegistered).Kind() {
	case reflect.Slice:
		return rollOutOfInt(flipValue, flipValueRegistered)
	case reflect.Int:
		return isFlipInt(flipValue, flipValueRegistered)
	default:
		return flipResponse(false), errors.New(fmt.Sprintf("unknown type, data_type=%s", reflect.TypeOf(flipValueRegistered).Kind().String()))

	}
}
