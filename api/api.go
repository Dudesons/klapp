package api

import (
	"github.com/dudesons/klapp/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"github.com/dudesons/klapp/utils"
	"github.com/dudesons/klapp/flip"
	"github.com/dudesons/klapp/cache"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type flipServer struct{
	Cache     	map[string]*utils.Flip
	ChanCache	chan *map[string]*utils.Flip
	FlipStore 	utils.FlipStore
	Config    	*utils.KlappConfig
	Logger    	*zap.Logger
}

func newFlipServer(chanCache chan *map[string]*utils.Flip, config *utils.KlappConfig, flipStore utils.FlipStore, logger *zap.Logger) utils.FlipServer {
	return &flipServer{
		ChanCache: chanCache,
		Config: config,
		FlipStore: flipStore,
		Logger: logger,
	}
}

func (s *flipServer) UpdateCache() {
	go cache.WatchTree(&s.Config.FlipPrefix, &s.Config.ConsulEndpoint, s.ChanCache, s.Logger)
	for {
		select {
		case cacheFlips := <-s.ChanCache:
			s.Logger.Info("flip_store_sync")
			s.Cache = *cacheFlips
		default:
		}
	}
}

func (s *flipServer) IsFlip(ctx context.Context, flipRequest *pb.FlipRequest) (*pb.FlipResponse, error) {
	var isFlip *bool
	var err error

	if s.Config.CacheEnable {
		_, ok := s.Cache[flipRequest.Flip]
		if ok {
			s.Logger.Info("flip_request_in_cache", zap.String("key", flipRequest.Flip), zap.String("value", flipRequest.Target))
			isFlip, err = flip.Factory(&flipRequest.Target, s.Cache[flipRequest.Flip])
		} else {
			s.Logger.Warn("flip_not_found")
			return &pb.FlipResponse{Activated: false}, grpc.Errorf(codes.NotFound, "flip_not_found")
		}

	} else {
		s.Logger.Info("flip_request_in_flipstore", zap.String("key", flipRequest.Flip), zap.String("value", flipRequest.Target))
		isFlip, err = s.FlipStore.Get(flipRequest)
	}

	if err != nil {
		s.Logger.Error("flip_request_error", zap.String("flip_requested", flipRequest.Flip), zap.String("caused_by", err.Error()))
		// TODO(Create a real error management)
		return &pb.FlipResponse{Activated: *isFlip}, grpc.Errorf(codes.Internal, err.Error())
	}

	return &pb.FlipResponse{Activated: *isFlip}, nil

}

func (s *flipServer) Health(ctx context.Context, flipRequest *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{}, nil
}

func Run() error {
	chanCache := make(chan *map[string]*utils.Flip)
	var config utils.KlappConfig
	err := envconfig.Process("klapp", &config)
	kv, err := flip.FlipNewStoreFactory(&config)
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	flipServer := newFlipServer(chanCache, &config, kv, logger)
	go flipServer.UpdateCache()

	server := grpc.NewServer()
	pb.RegisterFlipServer(server, flipServer)
	server.Serve(listen)
	return nil
}
