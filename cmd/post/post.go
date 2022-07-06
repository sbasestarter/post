package main

import (
	"context"
	"time"

	"github.com/gorilla/mux"
	"github.com/sbasestarter/post/internal/config"
	"github.com/sbasestarter/post/internal/post/server"
	"github.com/sbasestarter/proto-repo/gen/protorepo-post-go"
	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libconfig"
	"github.com/sgostarter/libeasygo/stg"
	"github.com/sgostarter/liblogrus"
	"github.com/sgostarter/librediscovery"
	"github.com/sgostarter/libservicetoolset/servicetoolset"
	"google.golang.org/grpc"
)

func main() {
	logger := l.NewWrapper(liblogrus.NewLogrus())
	logger.GetLogger().SetLevel(l.LevelDebug)

	var cfg config.Config

	_, err := libconfig.Load("config.yaml", &cfg)
	if err != nil {
		logger.Fatalf("load config failed: %v", err)

		return
	}

	redisCli, err := stg.InitRedis(cfg.RedisDSN)
	if err != nil {
		panic(err)
	}

	cfg.GRpcServerConfig.DiscoveryExConfig.Setter, err = librediscovery.NewSetter(context.Background(), logger, redisCli,
		"", time.Minute)
	if err != nil {
		logger.Fatalf("create rediscovery setter failed: %v", err)

		return
	}

	postServer := server.NewPostServer(context.Background(), &cfg, logger)
	serviceToolset := servicetoolset.NewServerToolset(context.Background(), logger)
	_ = serviceToolset.CreateGRpcServer(&cfg.GRpcServerConfig, nil, func(s *grpc.Server) error {
		postpb.RegisterPostServiceServer(s, postServer)

		return nil
	})

	r := mux.NewRouter()
	postServer.HTTPRegister(r)
	cfg.HttpServerConfig.Handler = r
	_ = serviceToolset.CreateHTTPServer(&cfg.HttpServerConfig)
	serviceToolset.Wait()
}
