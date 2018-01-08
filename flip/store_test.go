package flip

import (
	"testing"
	"github.com/dudesons/klapp/config"
	"github.com/stretchr/testify/assert"
	"github.com/kelseyhightower/envconfig"
	"fmt"
	"go.uber.org/zap"
)

func TestFlipStoreClient_Get(t *testing.T) {
	var conf config.KlappConfig
	err := envconfig.Process("klapp", &conf)
	if err != nil {
		t.Fatal(err)
	}
	conf.FlipStoreType = "mock"

	store, err := NewFlipStore(&conf)
	assert.Nil(t, err, "get a flip store should not fail")

	for _, fixture := range flipsFixture {
		resp, err := store.Get(&fixture.Flip)
		assert.Nil(t, err, fmt.Sprintf("Flip request on %s should not fail", fixture.Flip))

		assert.IsType(t, &Flip{}, resp, "Flip response type should be pointer on Flip")
		assert.Equal(t, &fixture.Flip, resp.Name)
	}
}

func TestFlipStoreClient_Get_FailureUnknownTag(t *testing.T) {
	featureTag := "404 tag not found"
	var conf config.KlappConfig
	err := envconfig.Process("klapp", &conf)
	if err != nil {
		t.Fatal(err)
	}
	conf.FlipStoreType = "mock"

	store, err := NewFlipStore(&conf)
	assert.Nil(t, err, "get a flip store should not fail")

	resp, err := store.Get(&featureTag)
	assert.NotNil(t, err, "The tag doesn't exist")

	assert.Nil(t, resp, "Flip response should be nil")
}

func TestFlipStoreClient_Watch(t *testing.T) {
	chanCache := make(chan *map[string]*Flip)
	var conf config.KlappConfig
	err := envconfig.Process("klapp", &conf)
	if err != nil {
		t.Fatal(err)
	}
	conf.FlipStoreType = "mock"

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	stringPointer := func(s string) *string{
		return &s
	}

	expectedresponse := map[string]*Flip{
		"Foo": {
			Activated: flipResponse(true),
			Name: stringPointer("Foo"),
			Description: stringPointer(""),
			Content: float64(42),
			},
		"Bar": {
			Activated: flipResponse(false),
			Name: stringPointer("Bar"),
			Description: stringPointer(""),
			Content: "toto",
			},
	}

	store, err := NewFlipStore(&conf)
	assert.Nil(t, err, "get a flip store should not fail")

	go store.Watch(chanCache, logger)
	cacheFlips := <-chanCache

	assert.IsType(t, &expectedresponse, cacheFlips)
	assert.Equal(t, &expectedresponse, cacheFlips)
}
