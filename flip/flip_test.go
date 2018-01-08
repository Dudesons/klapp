package flip

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFlipResponse(t *testing.T) {
	var expected *bool
	r0 := flipResponse(true)
	r1 := flipResponse(false)

	assert.IsType(t, expected, r0, "Expected pointer boolean")
	assert.IsType(t, expected, r1, "Expected pointer boolean")
}

func TestRollOutOfString(t *testing.T) {
	featureTag := "feature42"
	var flipValues interface{}
	var expectedType *bool
	flipValues = []string{"feature1", "feature2", "feature42", "feature99"}

	isActivated, err := rollOutOfString(&featureTag, flipValues)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(true), isActivated,"Expected pointer boolean with value true")
	assert.Nil(t, err)
}

func TestRollOutOfStringFailureBadFlipValue(t *testing.T) {
	featureTag := "feature42"
	var flipValues interface{}
	var expectedType *bool

	flipValues = []int{0, 1, 2}

	isActivated, err := rollOutOfString(&featureTag, flipValues)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated,"Expected pointer boolean with value false")
	assert.NotNil(t, err)
}

func TestRollOutOfStringFailureTagNotfound(t *testing.T) {
	featureTag := "feature42"
	var flipValues interface{}
	var expectedType *bool
	flipValues = []string{"feature1", "feature2", "feature51", "feature99"}

	isActivated, err := rollOutOfString(&featureTag, flipValues)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated,"Expected pointer boolean with value false")
	assert.Nil(t, err)
}

func TestRollOutOfInt(t *testing.T) {
	var featureTag int64
	featureTag = 42
	var flipValues interface{}
	var expectedType *bool
	flipValues = []int{10, 42, 56, 254}

	isActivated, err := rollOutOfInt(&featureTag, flipValues)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(true), isActivated,"Expected pointer boolean with value true")
	assert.Nil(t, err)
}

func TestRollOutOfIntFailureBadFlipValue(t *testing.T) {
	var featureTag int64
	featureTag = 42
	var flipValues interface{}
	var expectedType *bool
	flipValues = []string{"feature42", "feature1", "feature2"}

	isActivated, err := rollOutOfInt(&featureTag, flipValues)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated,"Expected pointer boolean with value false")
	assert.NotNil(t, err)
}

func TestRollOutOfIntFailureTagNotfound(t *testing.T) {
	var featureTag int64
	featureTag = 42
	var flipValues interface{}
	var expectedType *bool
	flipValues = []int{0, 1, 2}

	isActivated, err := rollOutOfInt(&featureTag, flipValues)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated,"Expected pointer boolean with value false")
	assert.Nil(t, err)
}

func TestIsFlipString(t *testing.T) {
	featureTag := "feature42"
	var flipValue interface{}
	var expectedType *bool
	flipValue = "feature42"

	isActivated, err := isFlipString(&featureTag, flipValue)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(true), isActivated,"Expected pointer boolean with value true")
	assert.Nil(t, err)
}

func TestIsFlipStringFailureBadFlipValue(t *testing.T) {
	featureTag := "feature42"
	var flipValue interface{}
	var expectedType *bool

	flipValue = 42

	isActivated, err := isFlipString(&featureTag, flipValue)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated,"Expected pointer boolean with value false")
	assert.NotNil(t, err)
}

func TestIsFlipStringFailureTagNotEqual(t *testing.T) {
	featureTag := "feature42"
	var flipValue interface{}
	var expectedType *bool
	flipValue = "feature1"

	isActivated, err := isFlipString(&featureTag, flipValue)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated, "Expected pointer boolean with value false")
	assert.Nil(t, err)
}

func TestIsFlipInt(t *testing.T) {
	var featureTag int64
	featureTag = 42
	var flipValue interface{}
	var expectedType *bool
	flipValue = 42

	isActivated, err := isFlipInt(&featureTag, flipValue)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(true), isActivated, "Expected pointer boolean with value false")
	assert.Nil(t, err)
}

func TestIsFlipIntFailureBadFlipValue(t *testing.T) {
	var featureTag int64
	featureTag = 42
	var flipValue interface{}
	var expectedType *bool

	flipValue = "feature78"

	isActivated, err := isFlipInt(&featureTag, flipValue)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated, "Expected pointer boolean with value ")
	assert.NotNil(t, err)
}

func TestIsFlipIntFailureTagNotEqual(t *testing.T) {
	var featureTag int64
	featureTag = 42
	var flipValue interface{}
	var expectedType *bool
	flipValue = 51

	isActivated, err := isFlipInt(&featureTag, flipValue)
	assert.IsType(t, expectedType, isActivated,"Expected pointer boolean")
	assert.Equal(t, flipResponse(false), isActivated, "Expected pointer boolean with value ")
	assert.Nil(t, err)
}
