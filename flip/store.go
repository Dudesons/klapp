package flip

import (
	"time"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/docker/libkv/store/boltdb"
	"github.com/docker/libkv"
	"github.com/dudesons/klapp/config"
	"github.com/docker/libkv/store/consul"
	"fmt"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
)

type Flip struct {
	Activated 	*bool 		`json:"activated"`
	Name 		*string 	`json:"name"`
	Description *string 	`json:"description"`
	Content     interface{} `json:"content"`
}


type StoreWriter interface {
	// Get a value given its key
	Put(key string, value []byte, options *store.WriteOptions) error
}

type StoreReader interface {
	// Get a value given its key
	Get(key string) (*store.KVPair, error)
}

type StoreWatcher interface {
	WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*store.KVPair, error)
}

type StoreCleanerWriterReaderWatcher interface {
	StoreWriter
	StoreReader
	StoreWatcher
	StoreCleaner
}

type StoreCleaner interface {
	Delete(key string) error
	DeleteTree(directory string) error
}

type FlipStore interface {
	Get(featureTag *string) (*Flip, error)
	Watch(flipUpdate chan *map[string]*Flip, logger *zap.Logger)
	Put(key string, flip *Flip) error
	Delete(target string, isDir bool) error
}

// FakeStore is used only for running unit test
type FakeStore struct {}

// Get a value given its key
func (fs *FakeStore) Get(key string) (*store.KVPair, error) {
	var flips map[string]Flip
	
	json.Unmarshal([]byte(`
{
  "klapp/flips//boolon": {
    "activated": true,
    "name": "boolon",
    "description": ""
  },
  "klapp/flips//booloff": {
    "activated": false,
    "name": "booloff",
    "description": ""
  },
  "klapp/flips//string": {
    "activated": true,
    "name": "string",
    "description": "",
    "content": "mystring"
  },
  "klapp/flips//stringslice": {
    "activated": true,
    "name": "stringslice",
    "description": "",
    "content": ["q", "w", "e", "r", "t", "y"]
  },
  "klapp/flips//int": {
    "activated": true,
    "name": "int",
    "description": "",
    "content": 42
  },
  "klapp/flips//intslice": {
    "activated": true,
    "name": "intslice",
    "description": "",
    "content": [0, 1, 2, 3, 4, 5]
  }
}
`), &flips)
	_, ok := flips[key]
	if ok {
		data, err := json.Marshal(flips[key])
		if err != nil {
			panic(err)
		}
		return &store.KVPair{Value:data}, nil
	}
	return nil, errors.New("Unknown key")
}

func (fs *FakeStore) WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*store.KVPair, error) {
	chanKV := make(chan []*store.KVPair)
	var lastIndex uint64
	go func() {
		chanKV <- []*store.KVPair{
			{Key: "Foo", Value: []byte(`{"activated": true, "name": "Foo", "description": "", "content": 42}`), LastIndex: lastIndex},
			{Key: "Bar", Value: []byte(`{"activated": false, "name": "Bar", "description": "", "content": "toto"}`), LastIndex: lastIndex},
		}
	}()
	return chanKV, nil
}

func (fs *FakeStore) Put(key string, value []byte, options *store.WriteOptions) error {
	var err error
	return err
}

func (fs *FakeStore) Delete(key string) error  {
	var err error
	return err
}

func (fs *FakeStore) DeleteTree(directory string) error  {
	var err error
	return err
}


func initStore() {
	consul.Register()
	etcd.Register()
	zookeeper.Register()
	boltdb.Register()
}

func NewFlipStore (config *config.KlappConfig) (FlipStore, error) {
	var kvClient StoreCleanerWriterReaderWatcher
	var err error

	if config.FlipStoreType == "mock" {
		kvClient = &FakeStore{}
	} else {
		initStore()
		kvClient, err = libkv.NewStore(
			config.FlipStoreType,
			[]string{config.FlipStoreEndpoint},
			&store.Config{
				ConnectionTimeout: time.Duration(config.FlipStoreConnexionTimeout) * time.Second,
			},
		)
	}

	return &FlipStoreClient{kvClient: kvClient, Config: config}, err
}

type FlipStoreClient struct {
	kvClient StoreCleanerWriterReaderWatcher
	Config *config.KlappConfig
}

func (f *FlipStoreClient) Get(featureTag *string) (*Flip, error) {
	flipPayload := Flip{}
	response, err := f.kvClient.Get(fmt.Sprintf("%s/%s", f.Config.FlipPrefix, *featureTag))

	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New(fmt.Sprintf("flip: %s unknown", *featureTag))
	}

	err = json.Unmarshal(response.Value, &flipPayload)

	if err != nil {
		return nil, err
	}

	return &flipPayload, nil
}

func (f *FlipStoreClient) Watch(flipUpdate chan *map[string]*Flip, logger *zap.Logger) {
	stopCh := make(<-chan struct{})
	events, err := f.kvClient.WatchTree(f.Config.FlipPrefix, stopCh)

	if err != nil {
		logger.Fatal(
			"cache_watch_tree_failure",
			zap.String("caused_by", err.Error()),
			zap.String("flip_store_prefix", f.Config.FlipPrefix))
	}

	for {
		select {
		case pairs := <-events:
			flips := map[string]*Flip{}
			logger.Info("cache_watch_tree_event_trigger")
			// If there is a need to optimize we can try to send by chunk, this proto is all cache
			for _, pair := range pairs {
				flipPayload := Flip{}
				err = json.Unmarshal(pair.Value, &flipPayload)

				if err != nil {
					logger.Error(
						"cache_watch_tree_decode_failure",
						zap.String("caused_by", err.Error()),
						zap.String("flip_payload", string(pair.Value)))
				}

				if flipPayload.Name != nil {
					flips[*flipPayload.Name] = &flipPayload
				}
			}
			flipUpdate <- &flips
		}
	}
}

func (f *FlipStoreClient) Put(key string, flip *Flip) error {
	val, err := json.Marshal(flip)
	if err != nil {
		return err
	}

	return f.kvClient.Put(key, val, nil)
}

func (f *FlipStoreClient) Delete(target string, isDir bool) error  {
	if isDir {
		return f.kvClient.DeleteTree(target)
	}

	return f.kvClient.Delete(target)
}
