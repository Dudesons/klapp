package flip

import (
	"reflect"
	"errors"
	"fmt"
	"github.com/dudesons/klapp/utils"
)

func Factory(flipTarget *string, flip *utils.Flip) (*bool, error) {
	if flip.Content == nil {
		return flip.Activated, nil
	}

	switch reflect.TypeOf(flip.Content).Kind() {
	case reflect.Slice:
		if flip.Type == 0 {
			return rolloutOfString(flipTarget, reflect.ValueOf(flip.Content))
		} else if flip.Type == 1 {
			return rolloutOfInt(flipTarget, reflect.ValueOf(flip.Content))
		} else {
			return flipResponse(false), errors.New("data type is too complex only accept for rollout slice of string or int")
		}
	case reflect.String:
		return isFlipString(flipTarget, flip.Content.(string))
	case reflect.Int:
		return isFlipInt(flipTarget, flip.Content.(int))
	case reflect.Float64:
		return isFlipInt(flipTarget, int(flip.Content.(float64)))
	default:
		return flipResponse(false), errors.New(fmt.Sprintf("unknown type, data_type=%s", reflect.TypeOf(flip.Content).Kind().String()))

	}
}

func FlipNewStoreFactory(config *utils.KlappConfig) (utils.FlipStore, error) {
	switch config.FlipStoreType {
	case "consul":
		return newConsulStore(config), nil
	default:
		return nil, errors.New(fmt.Sprintf("Don't know flip store type %s", config.FlipStoreType))
	}
}

func flipResponse(b bool) *bool {
	return &b
}