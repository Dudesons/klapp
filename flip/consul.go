package flip

import (
	"github.com/dudesons/klapp/utils"
	consulapi "github.com/hashicorp/consul/api"
	"errors"
	"fmt"
	"github.com/dudesons/klapp/pb"
	"encoding/json"
)

func newConsulStore(config *utils.KlappConfig) utils.FlipStore {
	consulConfig := &consulapi.Config{}

	consulConfig.Address = config.ConsulEndpoint
	consulConfig.Scheme = config.ConsulScheme
	consulConfig.Datacenter = config.ConsulDatacenter

	client, err := consulapi.NewClient(consulConfig)

	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	return &ConsulStore{
		kv,
		config.ConsulAllowStale,
		config,
	}
}

type ConsulStore struct {
	kvClient utils.ConsulFlipStore
	// AllowStale allows any Consul server (non-leader) to service
	// a read. This allows for lower latency and higher throughput
	allowStale bool
	config *utils.KlappConfig
}

//func (c *ConsulStore) Get(flipName *string) (*pb.FlipResponse, error) {
func (c *ConsulStore) Get(request *pb.FlipRequest) (*bool, error) {
	flipPayload := utils.Flip{}

		response, _, err := c.kvClient.Get(fmt.Sprintf("%s/%s", c.config.FlipPrefix, request.Flip), &consulapi.QueryOptions{
		AllowStale: c.allowStale,
	})

	if err != nil {
		return flipResponse(false), err
	}

	if response == nil {
		return flipResponse(false), errors.New(fmt.Sprintf("flip: %s unknown", request.Flip))
	}

	err = json.Unmarshal(response.Value, &flipPayload)

	if err != nil {
		// TODO(Remove panic)
		panic(err)
	}

	return Factory(&request.Target, &flipPayload)
}