package main

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/jiuzhou-zhao/go-fundamental/servicetoolset"
	"github.com/sbasestarter/post/internal/config"
	"github.com/sbasestarter/post/internal/post/server"
	"github.com/sbasestarter/proto-repo/gen/protorepo-post-go"
	"github.com/sgostarter/libconfig"
	"google.golang.org/grpc"
)

func main() {
	loge.SetGlobalLogger(loge.NewLogger(nil))

	var cfg config.Config
	_, err := libconfig.Load("config", &cfg)
	if err != nil {
		loge.Fatalf(context.Background(), "load config failed: %v", err)
		return
	}

	postServer := server.NewPostServer(context.Background(), &cfg)
	serviceToolset := servicetoolset.NewServerToolset(context.Background(), loge.GetGlobalLogger().GetLogger())
	_ = serviceToolset.CreateGRpcServer(&cfg.GRpcServerConfig, nil, func(s *grpc.Server) {
		postpb.RegisterPostServiceServer(s, postServer)
	})

	r := mux.NewRouter()
	postServer.HTTPRegister(r)
	cfg.HttpServerConfig.Handler = r
	_ = serviceToolset.CreateHttpServer(&cfg.HttpServerConfig)
	serviceToolset.Wait()
}
