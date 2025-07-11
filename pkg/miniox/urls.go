package miniox

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/aeternitas-infinita/sloglog"
	"github.com/minio/minio-go/v7"
)

// GetPresignedURL generates a presigned URL for GET operation with automatic path prefix handling
func (c *Client) GetPresignedURL(ctx context.Context, objectPath string, expiry time.Duration) (*url.URL, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, err
	}

	fullPath := c.buildPath(objectPath)

	sloglog.DebugCtx(ctx, "[MinIO] Generating presigned GET URL",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath),
		slog.Duration("expiry", expiry))

	return c.minio.PresignedGetObject(ctx, c.bucketName, fullPath, expiry, nil)
}

// GetPresignedURLWithParams generates a presigned URL for GET operation with custom parameters
func (c *Client) GetPresignedURLWithParams(ctx context.Context, objectPath string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, err
	}

	fullPath := c.buildPath(objectPath)

	sloglog.DebugCtx(ctx, "[MinIO] Generating presigned GET URL with params",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath),
		slog.Duration("expiry", expiry))

	return c.minio.PresignedGetObject(ctx, c.bucketName, fullPath, expiry, reqParams)
}

// GetPresignedPutURL generates a presigned URL for PUT operation with automatic path prefix handling
func (c *Client) GetPresignedPutURL(ctx context.Context, objectPath string, expiry time.Duration) (*url.URL, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, err
	}

	fullPath := c.buildPath(objectPath)

	sloglog.DebugCtx(ctx, "[MinIO] Generating presigned PUT URL",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath),
		slog.Duration("expiry", expiry))

	return c.minio.PresignedPutObject(ctx, c.bucketName, fullPath, expiry)
}

// GetPresignedPostPolicy generates a presigned POST policy with automatic path prefix handling
func (c *Client) GetPresignedPostPolicy(ctx context.Context, policy *minio.PostPolicy) (*url.URL, map[string]string, error) {
	sloglog.DebugCtx(ctx, "[MinIO] Generating presigned POST policy",
		slog.String("bucket", c.bucketName))

	// Note: PostPolicy object key should be set with prefix applied before calling this method
	return c.minio.PresignedPostPolicy(ctx, policy)
}

// GetPublicURL generates a public URL for an object (requires public bucket or appropriate policy)
func (c *Client) GetPublicURL(objectPath string) (*url.URL, error) {
	if c.publicBaseURL == "" {
		return nil, fmt.Errorf("public base URL not configured")
	}

	if err := c.ValidatePath(objectPath); err != nil {
		return nil, err
	}

	fullPath := c.buildPath(objectPath)

	// Clean up the public base URL
	baseURL := strings.TrimSuffix(c.publicBaseURL, "/")

	// Construct the full public URL
	publicURLString := fmt.Sprintf("%s/%s/%s", baseURL, c.bucketName, fullPath)

	return url.Parse(publicURLString)
}

// ComposeObject composes an object from existing objects with automatic path prefix handling
func (c *Client) ComposeObject(ctx context.Context, destObjectPath string, srcObjects []minio.CopySrcOptions, opts minio.CopyDestOptions) (minio.UploadInfo, error) {
	if err := c.ValidatePath(destObjectPath); err != nil {
		return minio.UploadInfo{}, err
	}

	fullDestPath := c.buildPath(destObjectPath)

	// Apply prefix to source objects if they are from the same bucket
	for i := range srcObjects {
		if srcObjects[i].Bucket == c.bucketName {
			// Validate source object path
			if err := c.ValidatePath(srcObjects[i].Object); err != nil {
				return minio.UploadInfo{}, fmt.Errorf("invalid source object path %s: %w", srcObjects[i].Object, err)
			}
			srcObjects[i].Object = c.buildPath(srcObjects[i].Object)
		}
	}

	sloglog.DebugCtx(ctx, "[MinIO] Composing object",
		slog.String("bucket", c.bucketName),
		slog.String("dest", fullDestPath),
		slog.Int("sources", len(srcObjects)))

	// Set the destination in the opts
	opts.Bucket = c.bucketName
	opts.Object = fullDestPath

	uploadInfo, err := c.minio.ComposeObject(ctx, opts, srcObjects...)
	if err != nil {
		return uploadInfo, err
	}

	// Strip base path from returned upload info
	uploadInfo.Key = c.stripBasePath(uploadInfo.Key)
	return uploadInfo, nil
}

// PresignedGetObject generates a presigned URL for GET operation with automatic path prefix handling
// This is an alias for GetPresignedURL for consistency with MinIO naming
func (c *Client) PresignedGetObject(ctx context.Context, objectPath string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
	return c.GetPresignedURLWithParams(ctx, objectPath, expiry, reqParams)
}

// PresignedPutObject generates a presigned URL for PUT operation with automatic path prefix handling
// This is an alias for GetPresignedPutURL for consistency with MinIO naming
func (c *Client) PresignedPutObject(ctx context.Context, objectPath string, expiry time.Duration) (*url.URL, error) {
	return c.GetPresignedPutURL(ctx, objectPath, expiry)
}

// PresignedHeadObject generates a presigned URL for HEAD operation with automatic path prefix handling
func (c *Client) PresignedHeadObject(ctx context.Context, objectPath string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, err
	}

	fullPath := c.buildPath(objectPath)

	sloglog.DebugCtx(ctx, "[MinIO] Generating presigned HEAD URL",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath),
		slog.Duration("expiry", expiry))

	return c.minio.PresignedHeadObject(ctx, c.bucketName, fullPath, expiry, reqParams)
}

// PresignedPostPolicyForUpload creates a presigned POST policy for browser-based uploads
func (c *Client) PresignedPostPolicyForUpload(ctx context.Context, objectPath string, expiry time.Duration, maxSize int64) (*url.URL, map[string]string, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, nil, err
	}

	fullPath := c.buildPath(objectPath)

	sloglog.DebugCtx(ctx, "[MinIO] Generating presigned POST policy for upload",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath),
		slog.Duration("expiry", expiry),
		slog.Int64("maxSize", maxSize))

	policy := minio.NewPostPolicy()
	policy.SetBucket(c.bucketName)
	policy.SetKey(fullPath)
	policy.SetExpires(time.Now().Add(expiry))

	// Set content length range if specified
	if maxSize > 0 {
		policy.SetContentLengthRange(1, maxSize)
	}

	return c.minio.PresignedPostPolicy(ctx, policy)
}

// PresignedPostPolicyWithConditions creates a presigned POST policy with custom conditions
func (c *Client) PresignedPostPolicyWithConditions(ctx context.Context, objectPath string, expiry time.Duration, contentType string, maxSize int64) (*url.URL, map[string]string, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, nil, err
	}

	fullPath := c.buildPath(objectPath)

	sloglog.DebugCtx(ctx, "[MinIO] Generating presigned POST policy with conditions",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath),
		slog.Duration("expiry", expiry),
		slog.String("contentType", contentType),
		slog.Int64("maxSize", maxSize))

	policy := minio.NewPostPolicy()
	policy.SetBucket(c.bucketName)
	policy.SetKey(fullPath)
	policy.SetExpires(time.Now().Add(expiry))

	// Set content type if specified
	if contentType != "" {
		policy.SetContentType(contentType)
	}

	// Set content length range if specified
	if maxSize > 0 {
		policy.SetContentLengthRange(1, maxSize)
	}

	return c.minio.PresignedPostPolicy(ctx, policy)
}
