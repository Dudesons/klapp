package flip

import (
	"strings"
	"github.com/teh-cmc/cast"
	"errors"
	"strconv"
)

func flipResponse(b bool) *bool {
	return &b
}

func flipTag(s string) *string {
	return &s
}

func rollOutOfString(flipTarget *string, flipValues interface{}) (*bool, error){
	flips, err := cast.ToStringSliceE(flipValues)
	if err != nil {
		return flipResponse(false), err
	}

	for _, i := range flips {
		if strings.Compare(*flipTarget, i) == 0 {
			return flipResponse(true), nil
		}
	}

	return flipResponse(false), nil
}

func rollOutOfInt(flipTarget *int64, flipValues interface{}) (*bool, error){
	flips, err := cast.ToIntSliceE(flipValues)
	if err != nil {
		return flipResponse(false), err
	}

	for _, i := range flips {
		if *flipTarget == int64(i) {
			return flipResponse(true), nil
		}
	}

	return flipResponse(false), nil
}

func isFlipString(flipTarget *string, flipValue interface{}) (result *bool, err error){
	defer func() {
		if r := recover(); r != nil {
			result = flipResponse(false)
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case int:
				err = errors.New(strconv.Itoa(x))
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	return flipResponse(strings.Compare(*flipTarget, flipValue.(string)) == 0), nil
}

func isFlipInt(flipTarget *int64, flipValue interface{}) (result *bool, err error){
	defer func() {
		if r := recover(); r != nil {
			result = flipResponse(false)
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case int:
				err = errors.New(strconv.Itoa(x))
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	val, err := cast.ToIntE(flipValue)
	if err != nil {
		panic(err)
	}
	return flipResponse(*flipTarget == int64(val)), nil
}
