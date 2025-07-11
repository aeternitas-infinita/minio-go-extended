// Package miniox provides an extended MinIO Go client library
package miniox

import "github.com/aeternitas-infinita/minio-go-extended/pkg/miniox"

// Config is a type alias for the main configuration
type Config = miniox.Config

// Client is a type alias for the main client
type Client = miniox.Client

// New creates a new MinIO extended client
func New(config *Config) (*Client, error) {
	return miniox.New(config)
}
