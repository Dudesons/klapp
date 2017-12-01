package main

import (
	"github.com/dudesons/klapp/api"
	"flag"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := api.Run(); err != nil {
		glog.Fatal(err)
	}
}
