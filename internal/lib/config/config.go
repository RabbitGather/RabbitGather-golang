package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/meowalien/RabbitGather-golang.git/internal/lib/errs"
)

var ConfigFile = "config.yaml"

type Config struct {
	GRPCServer GRPCConfig `yaml:"grpc_server"`
}

type GRPCConfig struct {
	Port uint32 `yaml:"port"`
}

func ParseConfig() (cf Config) {
	file, err := os.ReadFile(ConfigFile)
	if err != nil {
		panic(errs.New(err))
	}

	err = yaml.Unmarshal(file, &cf)
	if err != nil {
		panic(errs.New(err))
	}
	return
}
