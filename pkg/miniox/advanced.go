package miniox

import (
	"context"
	"log/slog"
	"time"

	"github.com/aeternitas-infinita/rmlog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/tags"
)

// GetObjectTagging gets the tags of an object with automatic path prefix handling
func (c *Client) GetObjectTagging(ctx context.Context, objectPath string, opts minio.GetObjectTaggingOptions) (*tags.Tags, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, err
	}

	fullPath := c.buildPath(objectPath)

	rmlog.DebugCtxMin(ctx, "[MinIO] Getting object tags",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	return c.minio.GetObjectTagging(ctx, c.bucketName, fullPath, opts)
}

// PutObjectTagging sets the tags of an object with automatic path prefix handling
func (c *Client) PutObjectTagging(ctx context.Context, objectPath string, objectTags *tags.Tags, opts minio.PutObjectTaggingOptions) error {
	if err := c.ValidatePath(objectPath); err != nil {
		return err
	}

	fullPath := c.buildPath(objectPath)

	rmlog.DebugCtxMin(ctx, "[MinIO] Setting object tags",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	return c.minio.PutObjectTagging(ctx, c.bucketName, fullPath, objectTags, opts)
}

// RemoveObjectTagging removes all tags from an object with automatic path prefix handling
func (c *Client) RemoveObjectTagging(ctx context.Context, objectPath string, opts minio.RemoveObjectTaggingOptions) error {
	if err := c.ValidatePath(objectPath); err != nil {
		return err
	}

	fullPath := c.buildPath(objectPath)

	rmlog.DebugCtxMin(ctx, "[MinIO] Removing object tags",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	return c.minio.RemoveObjectTagging(ctx, c.bucketName, fullPath, opts)
}

// GetObjectRetention gets the retention settings of an object with automatic path prefix handling
func (c *Client) GetObjectRetention(ctx context.Context, objectPath string, versionID string) (*minio.RetentionMode, *time.Time, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, nil, err
	}

	fullPath := c.buildPath(objectPath)

	rmlog.DebugCtxMin(ctx, "[MinIO] Getting object retention",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath),
		slog.String("versionID", versionID))

	return c.minio.GetObjectRetention(ctx, c.bucketName, fullPath, versionID)
}

// PutObjectRetention sets the retention settings of an object with automatic path prefix handling
func (c *Client) PutObjectRetention(ctx context.Context, objectPath string, opts minio.PutObjectRetentionOptions) error {
	if err := c.ValidatePath(objectPath); err != nil {
		return err
	}

	fullPath := c.buildPath(objectPath)

	rmlog.DebugCtxMin(ctx, "[MinIO] Setting object retention",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	return c.minio.PutObjectRetention(ctx, c.bucketName, fullPath, opts)
}

// GetObjectLegalHold gets the legal hold status of an object with automatic path prefix handling
func (c *Client) GetObjectLegalHold(ctx context.Context, objectPath string, opts minio.GetObjectLegalHoldOptions) (*minio.LegalHoldStatus, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, err
	}

	fullPath := c.buildPath(objectPath)

	rmlog.DebugCtxMin(ctx, "[MinIO] Getting object legal hold",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	return c.minio.GetObjectLegalHold(ctx, c.bucketName, fullPath, opts)
}

// PutObjectLegalHold sets the legal hold status of an object with automatic path prefix handling
func (c *Client) PutObjectLegalHold(ctx context.Context, objectPath string, opts minio.PutObjectLegalHoldOptions) error {
	if err := c.ValidatePath(objectPath); err != nil {
		return err
	}

	fullPath := c.buildPath(objectPath)

	rmlog.DebugCtxMin(ctx, "[MinIO] Setting object legal hold",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	return c.minio.PutObjectLegalHold(ctx, c.bucketName, fullPath, opts)
}

// SelectObjectContent performs SQL select on object content with automatic path prefix handling
func (c *Client) SelectObjectContent(ctx context.Context, objectPath string, opts minio.SelectObjectOptions) (*minio.SelectResults, error) {
	// Note: We don't validate path here as SelectObjectContent might work with special paths
	// and we want to maintain compatibility with the underlying MinIO client

	fullPath := c.buildPath(objectPath)

	rmlog.DebugCtxMin(ctx, "[MinIO] Selecting object content",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	return c.minio.SelectObjectContent(ctx, c.bucketName, fullPath, opts)
}
