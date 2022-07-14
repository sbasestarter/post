package config

import (
	"github.com/sgostarter/libservicetoolset/clienttoolset"
	"github.com/sgostarter/libservicetoolset/dbtoolset"
	"github.com/sgostarter/libservicetoolset/servicetoolset"
)

type Config struct {
	GRpcServerConfig    servicetoolset.GRPCServerConfig  `yaml:"grpc_server_config"`
	GRpcClientConfigTpl clienttoolset.GRPCClientConfig   `yaml:"grpc_client_config_tpl"`
	HTTPServerConfig    servicetoolset.HTTPServerConfig  `yaml:"http_server_config"`
	DbConfig            dbtoolset.Config                 `yaml:"db_config"`
	ProtocolProviders   map[string]map[string][]EndPoint `yaml:"protocol_providers"`

	RedisDSN string `yaml:"redis_dsn"`
}

type EndPoint struct {
	Name     string `yaml:"name"`
	Argument string `yaml:"argument"`
}
