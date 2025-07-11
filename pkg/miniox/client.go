// Package miniox provides an extended MinIO Go client with enhanced functionality
// and simplified API for object storage operations.
package miniox

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/aeternitas-infinita/sloglog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config represents the configuration for MinIO client initialization
type Config struct {
	Endpoint      string // MinIO server endpoint (e.g., "localhost:9000")
	AccessKey     string // Access key for authentication
	SecretKey     string // Secret key for authentication
	UseSSL        bool   // Whether to use HTTPS
	BucketName    string // Default bucket name for operations
	BaseDirPrefix string // Optional: Base directory prefix for all operations
	PublicURL     string // Optional: Public URL for generating accessible links
}

// Client represents an extended MinIO client with additional functionality
type Client struct {
	minio         *minio.Client
	bucketName    string
	baseDirPrefix string
	publicBaseURL string
}

// New creates and initializes a new MinIO extended client
func New(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}

	if config.AccessKey == "" {
		return nil, fmt.Errorf("access key is required")
	}

	if config.SecretKey == "" {
		return nil, fmt.Errorf("secret key is required")
	}

	if config.BucketName == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		sloglog.Error("error while creating minio client", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Check if bucket exists
	exists, err := client.BucketExists(context.Background(), config.BucketName)
	if err != nil {
		sloglog.Error("error checking bucket existence", slog.Any("error", err))
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", config.BucketName)
	}

	extendedClient := &Client{
		minio:         client,
		bucketName:    config.BucketName,
		baseDirPrefix: config.BaseDirPrefix,
		publicBaseURL: config.PublicURL,
	}

	sloglog.Info("[MinIO] successfully connected to MinIO",
		slog.String("endpoint", config.Endpoint),
		slog.String("bucket", config.BucketName))

	return extendedClient, nil
}

// GetBucketName returns the configured bucket name
func (c *Client) GetBucketName() string {
	return c.bucketName
}

// GetBaseDirPrefix returns the configured base directory prefix
func (c *Client) GetBaseDirPrefix() string {
	return c.baseDirPrefix
}

// GetPublicBaseURL returns the configured public base URL
func (c *Client) GetPublicBaseURL() string {
	return c.publicBaseURL
}

// buildPath constructs the full path with base directory prefix
// Ensures proper forward slash formatting for MinIO compatibility
func (c *Client) buildPath(path string) string {
	// Clean the input path: remove leading/trailing slashes and convert to forward slashes
	cleanPath := strings.Trim(filepath.ToSlash(path), "/")

	// If no base directory prefix is set, return the clean path
	if c.baseDirPrefix == "" {
		if cleanPath == "" {
			return ""
		}
		return cleanPath
	}

	// Clean the base directory prefix: remove leading/trailing slashes and convert to forward slashes
	cleanPrefix := strings.Trim(filepath.ToSlash(c.baseDirPrefix), "/")

	// If the clean path is empty, return just the prefix
	if cleanPath == "" {
		return cleanPrefix
	}

	// Combine prefix and path with forward slash
	fullPath := cleanPrefix + "/" + cleanPath

	return fullPath
}

// stripBasePath removes the base directory prefix from a full path
// This is useful when returning paths to external callers who expect relative paths
func (c *Client) stripBasePath(fullPath string) string {
	if c.baseDirPrefix == "" {
		return fullPath
	}

	cleanPrefix := strings.Trim(filepath.ToSlash(c.baseDirPrefix), "/")
	cleanFullPath := filepath.ToSlash(fullPath)

	// Check if the full path starts with the prefix
	if strings.HasPrefix(cleanFullPath, cleanPrefix+"/") {
		return strings.TrimPrefix(cleanFullPath, cleanPrefix+"/")
	} else if cleanFullPath == cleanPrefix {
		return ""
	}

	// If path doesn't start with prefix, return as-is (shouldn't happen in normal usage)
	return cleanFullPath
}

// ValidatePath ensures the path is safe and doesn't try to escape the base directory
func (c *Client) ValidatePath(path string) error {
	cleanPath := filepath.ToSlash(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected: %s", path)
	}

	// Check for absolute paths (should be relative to base directory)
	if strings.HasPrefix(cleanPath, "/") {
		return fmt.Errorf("absolute paths are not allowed: %s", path)
	}

	return nil
}
