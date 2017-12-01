package cache

import (
	"github.com/hashicorp/consul/watch"
	consulapi "github.com/hashicorp/consul/api"
	"encoding/json"
	"github.com/dudesons/klapp/utils"
	"go.uber.org/zap"
)


// TODO(Add interface to support not only consul)
func WatchTree (keyPrefix *string, flipStorePeer *string, flipUpdate chan *map[string]*utils.Flip, logger *zap.Logger) {
	params := make(map[string]interface{})
	params["prefix"] = *keyPrefix
	params["type"] = "keyprefix"

	wp, err := watch.Parse(params)
	defer wp.Stop()
	if err != nil {
		logger.Fatal(
			"cache_watch_tree_failure",
			zap.String("caused_by", err.Error()),
			zap.String("flip_store_prefix", *keyPrefix),
			zap.String("flip_store_peer", *flipStorePeer))
	}

	wp.Handler = func(idx uint64, data interface{}) {
		flips := map[string]*utils.Flip{}
		logger.Info("cache_watch_tree_event_trigger")
		// I have to test if I send all the cache or by chunk, this proto is all cache
		for _, i := range data.(consulapi.KVPairs) {
			flipPayload := utils.Flip{}
			err = json.Unmarshal(i.Value, &flipPayload)

			if err != nil {
				logger.Error(
					"cache_watch_tree_decode_failure",
					zap.String("caused_by", err.Error()),
					zap.String("flip_payload", string(i.Value)))
			}

			if flipPayload.Name != nil {
				flips[*flipPayload.Name] = &flipPayload
			}
		}
		flipUpdate <- &flips
	}

	if err := wp.Run(*flipStorePeer); err != nil {
		logger.Fatal(
			"cache_watch_tree_init_failure",
			zap.String("caused_by", err.Error()),
			zap.String("flip_store_prefix", *keyPrefix),
			zap.String("flip_store_peer", *flipStorePeer))
	}
}
