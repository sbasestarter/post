package config

import (
	"github.com/jiuzhou-zhao/go-fundamental/clienttoolset"
	"github.com/jiuzhou-zhao/go-fundamental/dbtoolset"
	"github.com/jiuzhou-zhao/go-fundamental/servicetoolset"
)

type Config struct {
	GRpcServerConfig    servicetoolset.GRpcServerConfig
	GRpcClientConfigTpl clienttoolset.GRpcClientConfig
	HttpServerConfig    servicetoolset.HttpServerConfig
	DbConfig            dbtoolset.DBConfig
	ProtocolProviders   map[string]map[string][]EndPoint
}

type EndPoint struct {
	Name     string
	Argument string
}
