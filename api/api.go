package api

import (
	"github.com/dudesons/klapp/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"github.com/dudesons/klapp/flip"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/dudesons/klapp/config"
	"fmt"
	"github.com/dudesons/klapp/proxy"
)

type stringFlipOperator func(flipValue *string, flipValueRegistered interface{}) (*bool, error)
type intFlipOperator func(flipValue *int64, flipValueRegistered interface{}) (*bool, error)

type flipOperator struct {
	string stringFlipOperator
	int    intFlipOperator
}

type flipServer struct{
	cache     	map[string] *flip.Flip
	chanCache	chan *map[string]*flip.Flip
	flipStore 	flip.FlipStore
	config    	*config.KlappConfig
	logger    	*zap.Logger
	operator    *flipOperator
}

func (s *flipServer) UpdateCache() {
	go s.flipStore.Watch(s.chanCache, s.logger)
	for {
		select {
		case cacheFlips := <-s.chanCache:
			s.logger.Info("flip_store_sync")
			s.cache = *cacheFlips
		default:
		}
	}
}

func (s *flipServer) IsFlip(ctx context.Context, flipRequest *pb.FlipRequest) (*pb.FlipResponse, error) {
	s.logger.Info("is_cache_enable", zap.Bool("cache enable", s.config.CacheEnable))
	if s.config.CacheEnable {
		_, ok := s.cache[flipRequest.FeatureTag]
		if ok {
			s.logger.Info("flip_request_in_cache", zap.String("key", flipRequest.FeatureTag))
			return &pb.FlipResponse{Activated: *s.cache[flipRequest.FeatureTag].Activated}, nil
		}
		s.logger.Warn("flip_not_found")
		return nil, status.Errorf(codes.NotFound, "flip_not_found")

	}
	s.logger.Info("flip_request_in_flipstore", zap.String("key", flipRequest.FeatureTag))
	isFlip, err := s.flipStore.Get(&flipRequest.FeatureTag)
	if err != nil {
		s.logger.Error("flip_request_error", zap.String("flip_requested", flipRequest.FeatureTag), zap.String("caused_by", err.Error()))
		// TODO(Create a real error management)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.FlipResponse{Activated: *isFlip.Activated}, nil
}

func (s *flipServer) IsFlipString(ctx context.Context, flipRequest *pb.FlipStringRequest) (*pb.FlipResponse, error) {
	if s.config.CacheEnable {
		_, ok := s.cache[flipRequest.FeatureTag]
		if ok {
			s.logger.Info("flip_request_in_cache", zap.String("key", flipRequest.FeatureTag), zap.String("value", flipRequest.FeatureValue))
			isFlip, err := s.operator.string(&flipRequest.FeatureValue, s.cache[flipRequest.FeatureTag].Content)
			if err != nil {
				s.logger.Error("flip_request_error", zap.String("flip_requested", flipRequest.FeatureTag), zap.String("caused_by", err.Error()))
				// TODO(Create a real error management)
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			return &pb.FlipResponse{Activated: *s.cache[flipRequest.FeatureTag].Activated && *isFlip}, nil
		}

		s.logger.Warn("flip_not_found")
		return nil, status.Errorf(codes.NotFound, "flip_not_found")
	}

	s.logger.Info("flip_request_in_flipstore", zap.String("key", flipRequest.FeatureTag))
	flipResponse, err := s.flipStore.Get(&flipRequest.FeatureTag)
	isFlip, err := s.operator.string(&flipRequest.FeatureValue, flipResponse.Content)

	if err != nil {
		s.logger.Error("flip_request_error", zap.String("flip_requested", flipRequest.FeatureTag), zap.String("caused_by", err.Error()))
		// TODO(Create a real error management)
		return &pb.FlipResponse{Activated: *flipResponse.Activated && *isFlip}, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.FlipResponse{Activated: *flipResponse.Activated && *isFlip}, nil
}

func (s *flipServer) IsFlipInteger(ctx context.Context, flipRequest *pb.FlipIntegerRequest) (*pb.FlipResponse, error) {
	if s.config.CacheEnable {
		_, ok := s.cache[flipRequest.FeatureTag]
		if ok {
			s.logger.Info("flip_request_in_cache", zap.String("key", flipRequest.FeatureTag), zap.Int64("value", flipRequest.FeatureValue))
			isFlip, err := s.operator.int(&flipRequest.FeatureValue, s.cache[flipRequest.FeatureTag].Content)
			if err != nil {
				s.logger.Error("flip_request_error", zap.String("flip_requested", flipRequest.FeatureTag), zap.String("caused_by", err.Error()))
				// TODO(Create a real error management)
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			return &pb.FlipResponse{Activated: *s.cache[flipRequest.FeatureTag].Activated && *isFlip}, nil
		}

		s.logger.Warn("flip_not_found")
		return nil, status.Errorf(codes.NotFound, "flip_not_found")
	}

	s.logger.Info("flip_request_in_flipstore", zap.String("key", flipRequest.FeatureTag))
	flipResponse, err := s.flipStore.Get(&flipRequest.FeatureTag)
	isFlip, err := s.operator.int(&flipRequest.FeatureValue, flipResponse.Content)

	if err != nil {
		s.logger.Error("flip_request_error", zap.String("flip_requested", flipRequest.FeatureTag), zap.String("caused_by", err.Error()))
		// TODO(Create a real error management)
		return &pb.FlipResponse{Activated: *flipResponse.Activated && *isFlip}, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.FlipResponse{Activated: *flipResponse.Activated && *isFlip}, nil
}

func (s *flipServer) Health(ctx context.Context, flipRequest *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{}, nil
}

func Run() error {
	chanCache := make(chan *map[string]*flip.Flip)
	var conf config.KlappConfig

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger.Info("Starting klapp server")
	logger.Info("Retrieving configuration")
	err = envconfig.Process("klapp", &conf)
	if err != nil {
		panic(err)
	}
	logger.Info("klapp_config", zap.Any("configuration", conf))
	logger.Info("Retrieving flip store client")
	kv, err := flip.NewFlipStore(&conf)
	if err != nil {
		panic(err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.RPCListen))
	if err != nil {
		return err
	}
	logger.Info("Listen port", zap.Int("port", conf.RPCListen))

	flipServer := &flipServer{
		chanCache: chanCache,
		config: &conf,
		flipStore: kv,
		logger: logger,
		operator: &flipOperator{
			string: flip.StringFlipOperator,
			int: flip.IntFlipOperator,
		},
	}

	// TODO(Maybe add a recovery strategy)
	if conf.CacheEnable {
		logger.Info("Starting cache")
		go flipServer.UpdateCache()
	}

	// TODO(Maybe add a recovery strategy)
	if conf.ProxyEnable {
		logger.Info("Starting http proxy server")
		go proxy.HttpProxy(&conf.ProxyToRPCEndpoint, &conf.ProxyListen)
	}

	server := grpc.NewServer()
	pb.RegisterFlipServer(server, flipServer)
	logger.Info("Klapp server started")
	return server.Serve(listen)
}
