package flip

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestStringFlipOperatorOnString(t *testing.T) {
	var expectedType *bool
	featureTag := "feature42"
	var flipValue interface{}
	flipValue = "feature42"

	isActivated, err := StringFlipOperator(&featureTag, flipValue)

	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(true), isActivated,"Expected pointer boolean")
	assert.Nil(t, err)
}

func TestStringFlipOperatorOnStringRollOut(t *testing.T) {
	var expectedType *bool
	featureTag := "feature42"
	var flipValues interface{}
	flipValues = []string{"feature1", "feature2", "feature42", "feature99"}

	isActivated, err := StringFlipOperator(&featureTag, flipValues)

	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(true), isActivated,"Expected pointer boolean")
	assert.Nil(t, err)
}

func TestStringFlipOperatorFailure(t *testing.T) {
	var expectedType *bool
	featureTag := "feature42"
	var flipValue interface{}
	flipValue = 852

	isActivated, err := StringFlipOperator(&featureTag, flipValue)

	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated,"Expected pointer boolean")
	assert.NotNil(t, err)
}

func TestIntFlipOperatorOnInt(t *testing.T) {
	var expectedType *bool
	var featureTag int64
	featureTag = 42
	var flipValue interface{}
	flipValue = 42

	isActivated, err := IntFlipOperator(&featureTag, flipValue)

	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(true), isActivated,"Expected pointer boolean")
	assert.Nil(t, err)
}

func TestIntFlipOperatorOnIntRollOut(t *testing.T) {
	var expectedType *bool
	var featureTag int64
	featureTag = 42
	var flipValues interface{}
	flipValues = []int{10, 42, 56, 254}

	isActivated, err := IntFlipOperator(&featureTag, flipValues)

	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(true), isActivated,"Expected pointer boolean")
	assert.Nil(t, err)

}
func TestIntFlipOperatorFailure(t *testing.T) {
	var expectedType *bool
	var featureTag int64
	featureTag = 42
	var flipValue interface{}
	flipValue = "feature562"

	isActivated, err := IntFlipOperator(&featureTag, flipValue)

	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated,"Expected pointer boolean")
	assert.NotNil(t, err)
}
