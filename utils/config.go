package utils

type KlappConfig struct {
	Retry              int    `default:"3" split_words:"true"`
	RawLogEnable       bool   `default:"true" split_words:"true"`
	JSONLogEnable      bool   `default:"false" split_words:"true"`
	FlipStoreType      string `split_words:"true"`

	/* Consul Store */
	ConsulDatacenter   string `split_words:"true"`
	ConsulAllowStale   bool   `default:"true" split_words:"true"`
	ConsulEndpoint     string `default:"127.0.0.1:8500" split_words:"true"`
	ConsulScheme       string `default:"http" split_words:"true"`
	/* End Consul Store */

	FlipPrefix         string `default:"klapp/flips/" split_words:"true"`
	CacheEnable		   bool   `default:"false" split_words:"true"`
}
