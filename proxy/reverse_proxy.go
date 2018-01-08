package proxy

import (
	"net/http"

	"golang.org/x/net/context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	gw "github.com/dudesons/klapp/pb"
	"fmt"
)

func HttpProxy(rpcEndpoint *string, proxyPort *int) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterFlipHandlerFromEndpoint(ctx, mux, *rpcEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", *proxyPort), mux)
}
