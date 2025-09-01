package miniox

import (
	"context"
	"io"
	"log/slog"

	"github.com/aeternitas-infinita/rmlog"
	"github.com/minio/minio-go/v7"
)

// StatObject performs StatObject with automatic bucket name and path prefix handling
func (c *Client) StatObject(ctx context.Context, objectPath string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return minio.ObjectInfo{}, err
	}

	fullPath := c.buildPath(objectPath)
	rmlog.DebugCtxMin(ctx, "[MinIO] Getting object info",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	info, err := c.minio.StatObject(ctx, c.bucketName, fullPath, opts)
	if err != nil {
		return info, err
	}

	// Strip base path from returned object info to maintain relative paths for external usage
	info.Key = c.stripBasePath(info.Key)
	return info, nil
}

// GetObject performs GetObject with automatic bucket name and path prefix handling
func (c *Client) GetObject(ctx context.Context, objectPath string, opts minio.GetObjectOptions) (*minio.Object, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return nil, err
	}

	fullPath := c.buildPath(objectPath)
	rmlog.DebugCtxMin(ctx, "[MinIO] Getting object",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	return c.minio.GetObject(ctx, c.bucketName, fullPath, opts)
}

// PutObject performs PutObject with automatic bucket name and path prefix handling
func (c *Client) PutObject(ctx context.Context, objectPath string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	if err := c.ValidatePath(objectPath); err != nil {
		return minio.UploadInfo{}, err
	}

	fullPath := c.buildPath(objectPath)
	rmlog.DebugCtxMin(ctx, "[MinIO] Putting object",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath),
		slog.Int64("size", objectSize))

	uploadInfo, err := c.minio.PutObject(ctx, c.bucketName, fullPath, reader, objectSize, opts)
	if err != nil {
		return uploadInfo, err
	}

	// Strip base path from returned upload info
	uploadInfo.Key = c.stripBasePath(uploadInfo.Key)
	return uploadInfo, nil
}

// RemoveObject performs RemoveObject with automatic bucket name and path prefix handling
func (c *Client) RemoveObject(ctx context.Context, objectPath string, opts minio.RemoveObjectOptions) error {
	if err := c.ValidatePath(objectPath); err != nil {
		return err
	}

	fullPath := c.buildPath(objectPath)
	rmlog.DebugCtxMin(ctx, "[MinIO] Removing object",
		slog.String("bucket", c.bucketName),
		slog.String("object", fullPath))

	return c.minio.RemoveObject(ctx, c.bucketName, fullPath, opts)
}

// ListObjects lists objects with automatic bucket name and path prefix handling
func (c *Client) ListObjects(ctx context.Context, prefix string, recursive bool) <-chan minio.ObjectInfo {
	if prefix != "" {
		if err := c.ValidatePath(prefix); err != nil {
			// Return a channel with the error
			errorCh := make(chan minio.ObjectInfo, 1)
			errorCh <- minio.ObjectInfo{Err: err}
			close(errorCh)
			return errorCh
		}
	}

	fullPrefix := c.buildPath(prefix)
	rmlog.DebugCtxMin(ctx, "[MinIO] Listing objects",
		slog.String("bucket", c.bucketName),
		slog.String("prefix", fullPrefix),
		slog.Bool("recursive", recursive))

	opts := minio.ListObjectsOptions{
		Prefix:    fullPrefix,
		Recursive: recursive,
	}

	objectCh := c.minio.ListObjects(ctx, c.bucketName, opts)

	// Create a new channel to strip base paths from returned objects
	strippedCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(strippedCh)
		for objectInfo := range objectCh {
			if objectInfo.Err == nil {
				objectInfo.Key = c.stripBasePath(objectInfo.Key)
			}
			strippedCh <- objectInfo
		}
	}()

	return strippedCh
}

// CopyObject copies an object from source to destination with automatic path handling
func (c *Client) CopyObject(ctx context.Context, destObjectPath string, srcObjectPath string, opts minio.CopyDestOptions) (minio.UploadInfo, error) {
	if err := c.ValidatePath(destObjectPath); err != nil {
		return minio.UploadInfo{}, err
	}
	if err := c.ValidatePath(srcObjectPath); err != nil {
		return minio.UploadInfo{}, err
	}

	fullDestPath := c.buildPath(destObjectPath)
	fullSrcPath := c.buildPath(srcObjectPath)

	rmlog.DebugCtxMin(ctx, "[MinIO] Copying object",
		slog.String("bucket", c.bucketName),
		slog.String("src", fullSrcPath),
		slog.String("dest", fullDestPath))

	// Create source object options
	srcOpts := minio.CopySrcOptions{
		Bucket: c.bucketName,
		Object: fullSrcPath,
	}

	uploadInfo, err := c.minio.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: c.bucketName,
		Object: fullDestPath,
	}, srcOpts)

	if err != nil {
		return uploadInfo, err
	}

	// Strip base path from returned upload info
	uploadInfo.Key = c.stripBasePath(uploadInfo.Key)
	return uploadInfo, nil
}
