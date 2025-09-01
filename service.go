package miniox

import "github.com/aeternitas-infinita/minio-go-extended/pkg/miniox"

type Config = miniox.Config

type Client = miniox.Client

func New(config *Config) (*Client, error) {
	return miniox.New(config)
}
