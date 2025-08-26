package miniox

import (
	"context"
	"log/slog"
	"strings"

	"github.com/aeternitas-infinita/rmlog"
	"github.com/minio/minio-go/v7"
)

// FolderExists checks if a folder exists with automatic path prefix handling
func (c *Client) FolderExists(ctx context.Context, folderPath string) (bool, error) {
	if err := c.ValidatePath(folderPath); err != nil {
		return false, err
	}

	fullPath := c.buildPath(folderPath)
	filePath := fullPath + "/.empty"

	rmlog.DebugCtx(ctx, "[MinIO] Checking folder existence",
		slog.String("bucket", c.bucketName),
		slog.String("folder", fullPath))

	_, err := c.minio.StatObject(ctx, c.bucketName, filePath, minio.StatObjectOptions{})
	if err == nil {
		return true, nil
	}
	if minio.ToErrorResponse(err).Code == "NoSuchKey" {
		return false, nil
	}
	return false, err
}

// CreateFolder creates an empty folder with automatic path prefix handling
func (c *Client) CreateFolder(ctx context.Context, folderPath string) error {
	if err := c.ValidatePath(folderPath); err != nil {
		return err
	}

	exists, err := c.FolderExists(ctx, folderPath)
	if exists {
		return nil // Folder already exists
	} else if err != nil {
		return err
	}

	fullPath := c.buildPath(folderPath)
	filePath := fullPath + "/.empty"

	rmlog.DebugCtx(ctx, "[MinIO] Creating folder",
		slog.String("bucket", c.bucketName),
		slog.String("folder", fullPath))

	_, err = c.minio.PutObject(ctx, c.bucketName, filePath, nil, 0, minio.PutObjectOptions{})
	return err
}

// RemoveFolder removes all objects with a given prefix (folder) with automatic path prefix handling
func (c *Client) RemoveFolder(ctx context.Context, folderPath string) error {
	if err := c.ValidatePath(folderPath); err != nil {
		return err
	}

	fullPath := c.buildPath(folderPath)
	if !strings.HasSuffix(fullPath, "/") {
		fullPath += "/"
	}

	rmlog.DebugCtx(ctx, "[MinIO] Removing folder",
		slog.String("bucket", c.bucketName),
		slog.String("folder", fullPath))

	opts := minio.ListObjectsOptions{
		Prefix:    fullPath,
		Recursive: true,
	}

	objectCh := c.minio.ListObjects(ctx, c.bucketName, opts)
	errorCh := c.minio.RemoveObjects(ctx, c.bucketName, objectCh, minio.RemoveObjectsOptions{})

	// Check for errors during removal
	for removeErr := range errorCh {
		if removeErr.Err != nil {
			rmlog.ErrorCtx(ctx, "[MinIO] Error removing object during folder deletion",
				slog.String("object", removeErr.ObjectName),
				slog.Any("error", removeErr.Err))
			return removeErr.Err
		}
	}

	return nil
}

// ListFolders lists folders (common prefixes) in the given path
func (c *Client) ListFolders(ctx context.Context, prefix string) ([]string, error) {
	if prefix != "" {
		if err := c.ValidatePath(prefix); err != nil {
			return nil, err
		}
	}

	fullPrefix := c.buildPath(prefix)
	if fullPrefix != "" && !strings.HasSuffix(fullPrefix, "/") {
		fullPrefix += "/"
	}

	rmlog.DebugCtx(ctx, "[MinIO] Listing folders",
		slog.String("bucket", c.bucketName),
		slog.String("prefix", fullPrefix))

	opts := minio.ListObjectsOptions{
		Prefix:    fullPrefix,
		Recursive: false,
	}

	objectCh := c.minio.ListObjects(ctx, c.bucketName, opts)
	var folders []string
	seenFolders := make(map[string]bool)

	for objectInfo := range objectCh {
		if objectInfo.Err != nil {
			return nil, objectInfo.Err
		}

		// Extract folder name from object key
		relativePath := c.stripBasePath(objectInfo.Key)
		if prefix != "" {
			// Remove the prefix from the relative path
			prefixToRemove := strings.Trim(prefix, "/") + "/"
			relativePath = strings.TrimPrefix(relativePath, prefixToRemove)
		}

		// Get the first directory component
		parts := strings.Split(relativePath, "/")
		if len(parts) > 1 && parts[0] != "" {
			folderName := parts[0]
			if !seenFolders[folderName] {
				folders = append(folders, folderName)
				seenFolders[folderName] = true
			}
		}
	}

	return folders, nil
}
