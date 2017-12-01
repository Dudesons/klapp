package flip

import (
	"strconv"
	"strings"
	"reflect"
)

func rolloutOfString(flipTarget *string, flipValues reflect.Value) (*bool, error){
	for i := 0; i < flipValues.Len(); i++ {
		if strings.Compare(*flipTarget, flipValues.Index(i).Interface().(string)) == 0 {
			return flipResponse(true), nil
		}
	}

	return flipResponse(false), nil
}

func rolloutOfInt(flipTarget *string, flipValues reflect.Value) (*bool, error){
	val, err := strconv.Atoi(*flipTarget)

	if err != nil {
		return flipResponse(false), err
	}

	for i := 0; i < flipValues.Len(); i++ {
		// Stange case, reflect return me float64 for int
		if val == int(flipValues.Index(i).Interface().(float64)) {
			return flipResponse(true), nil
		}
	}

	return flipResponse(false), nil
}

func isFlipString(flipTarget *string, flipValue string) (*bool, error){
	return flipResponse(strings.Compare(*flipTarget, flipValue) == 0), nil
}

func isFlipInt(flipTarget *string, flipValue int) (*bool, error){
	val, err := strconv.Atoi(*flipTarget)

	if err != nil {
		return flipResponse(false), err
	}
	return flipResponse(val == flipValue), nil
}