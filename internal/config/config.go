package config

import (
	"github.com/jiuzhou-zhao/go-fundamental/servicetoolset"
)

type Config struct {
	GRpcServerConfig  servicetoolset.GRpcServerConfig
	HttpServerConfig  servicetoolset.HttpServerConfig
	ProtocolProviders map[string]map[string][]EndPoint
}

type EndPoint struct {
	Name     string
	Argument string
}
