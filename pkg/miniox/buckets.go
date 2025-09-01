package miniox

import (
	"context"
	"log/slog"

	"github.com/aeternitas-infinita/rmlog"
	"github.com/minio/minio-go/v7"
)

// BucketExists checks if the configured bucket exists
func (c *Client) BucketExists(ctx context.Context) (bool, error) {
	rmlog.DebugCtxMin(ctx, "[MinIO] Checking bucket existence",
		slog.String("bucket", c.bucketName))

	return c.minio.BucketExists(ctx, c.bucketName)
}

// ListBuckets lists all buckets (no prefix applied here as it's bucket-level operation)
func (c *Client) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	rmlog.DebugCtxMin(ctx, "[MinIO] Listing all buckets")

	return c.minio.ListBuckets(ctx)
}

// GetBucketLocation gets the location of the configured bucket
func (c *Client) GetBucketLocation(ctx context.Context) (string, error) {
	rmlog.DebugCtxMin(ctx, "[MinIO] Getting bucket location",
		slog.String("bucket", c.bucketName))

	return c.minio.GetBucketLocation(ctx, c.bucketName)
}

// GetBucketPolicy gets the bucket policy for the configured bucket
func (c *Client) GetBucketPolicy(ctx context.Context) (string, error) {
	rmlog.DebugCtxMin(ctx, "[MinIO] Getting bucket policy",
		slog.String("bucket", c.bucketName))

	return c.minio.GetBucketPolicy(ctx, c.bucketName)
}

// SetBucketPolicy sets the bucket policy for the configured bucket
func (c *Client) SetBucketPolicy(ctx context.Context, policy string) error {
	rmlog.DebugCtxMin(ctx, "[MinIO] Setting bucket policy",
		slog.String("bucket", c.bucketName))

	return c.minio.SetBucketPolicy(ctx, c.bucketName, policy)
}
