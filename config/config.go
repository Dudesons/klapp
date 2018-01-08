package config

import (
	"github.com/docker/libkv/store"
)

type KlappConfig struct {
	Retry                     int           `default:"3" split_words:"true"`
	RawLogEnable              bool          `default:"true" split_words:"true"`
	JSONLogEnable             bool          `default:"false" split_words:"true"`
	FlipPrefix                string        `default:"klapp/flips/" split_words:"true"`
	CacheEnable		          bool          `default:"false" split_words:"true"`
	FlipStoreType             store.Backend `default:"consul" split_words:"true"`
	FlipStoreEndpoint         string        `split_words:"true"`
	FlipStoreConnexionTimeout int           `default:"10" split_words:"true"`
	RPCListen                 int           `default:"50051" split_words:"true"`
	ProxyEnable               bool          `default:"false" split_words:"true"`
	ProxyToRPCEndpoint        string        `default:"127.0.0.1:50051" split_words:"true"`
	ProxyListen               int           `default:"8080" split_words:"true"`
}
