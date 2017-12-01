package utils

import (
    apiconsul "github.com/hashicorp/consul/api"
    "github.com/dudesons/klapp/pb"
    "golang.org/x/net/context"
)

type FlipServer interface {
    IsFlip(context.Context, *pb.FlipRequest) (*pb.FlipResponse, error)
    Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error)
    UpdateCache()
}

type ConsulFlipStore interface {
    Get(key string, q *apiconsul.QueryOptions) (*apiconsul.KVPair, *apiconsul.QueryMeta, error)
}

type Flip struct {
    Activated 	*bool 		`json:"activated"`
    Name 		*string 	`json:"name"`
    Description *string 	`json:"description"`
    Content     interface{} `json:"content"`
    Type        int 		`json:"type"`
}

type FlipResponse struct {

}

type FlipStore interface {
    Get(request *pb.FlipRequest) (*bool, error)
}