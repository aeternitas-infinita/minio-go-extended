# MinIO Go Extended

A powerful and user-friendly extension library for MinIO Go client that provides enhanced functionality and simplified API for object storage operations.

## Features

- 🚀 **Simplified API**: Easy-to-use wrapper around the official MinIO Go client
- 📁 **Automatic Path Management**: Built-in support for base directory prefixes
- 🔒 **Path Validation**: Prevents path traversal attacks and ensures secure operations
- 📂 **Folder Operations**: Native support for folder creation and management
- 🔗 **Presigned URLs**: Easy generation of presigned URLs for secure file access
- 🏷️ **Object Tagging**: Built-in support for object tagging and metadata
- 📊 **Comprehensive Logging**: Structured logging with configurable levels
- 🛡️ **Error Handling**: Robust error handling with detailed error messages

## Installation

```bash
go get github.com/aeternitas-infinita/minio-go-extended
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "strings"
    
    "github.com/aeternitas-infinita/minio-go-extended/pkg/miniox"
    "github.com/minio/minio-go/v7"
)

func main() {
    // Initialize the client
    config := &miniox.Config{
        Endpoint:      "localhost:9000",
        AccessKey:     "your-access-key",
        SecretKey:     "your-secret-key",
        UseSSL:        false,
        BucketName:    "your-bucket",
        BaseDirPrefix: "app-data", // Optional: all operations will be relative to this path
        PublicURL:     "http://localhost:9000", // Optional: for generating public URLs
    }
    
    client, err := miniox.New(config)
    if err != nil {
        log.Fatal("Failed to initialize MinIO client:", err)
    }
    
    ctx := context.Background()
    
    // Upload a file
    uploadInfo, err := client.PutObject(ctx, "documents/hello.txt", 
        strings.NewReader("Hello, World!"), 13, minio.PutObjectOptions{})
    if err != nil {
        log.Fatal("Failed to upload file:", err)
    }
    
    fmt.Printf("File uploaded successfully: %s\n", uploadInfo.ETag)
    
    // Download a file
    object, err := client.GetObject(ctx, "documents/hello.txt", minio.GetObjectOptions{})
    if err != nil {
        log.Fatal("Failed to download file:", err)
    }
    defer object.Close()
    
    // Create a folder
    err = client.CreateFolder(ctx, "images")
    if err != nil {
        log.Fatal("Failed to create folder:", err)
    }
    
    // Generate a presigned URL (valid for 1 hour)
    presignedURL, err := client.GetPresignedURL(ctx, "documents/hello.txt", time.Hour)
    if err != nil {
        log.Fatal("Failed to generate presigned URL:", err)
    }
    
    fmt.Printf("Presigned URL: %s\n", presignedURL.String())
}
```

## API Reference

### Client Configuration

```go
type Config struct {
    Endpoint      string // MinIO server endpoint (e.g., "localhost:9000")
    AccessKey     string // Access key for authentication
    SecretKey     string // Secret key for authentication
    UseSSL        bool   // Whether to use HTTPS
    BucketName    string // Default bucket name for operations
    BaseDirPrefix string // Optional: Base directory prefix for all operations
    PublicURL     string // Optional: Public URL for generating accessible links
}
```

### Core Operations

#### Object Operations
- `PutObject(ctx, path, reader, size, opts)` - Upload an object
- `GetObject(ctx, path, opts)` - Download an object
- `StatObject(ctx, path, opts)` - Get object metadata
- `RemoveObject(ctx, path, opts)` - Delete an object
- `CopyObject(ctx, destPath, srcPath, opts)` - Copy an object

#### Folder Operations
- `CreateFolder(ctx, path)` - Create a folder
- `FolderExists(ctx, path)` - Check if a folder exists
- `RemoveFolder(ctx, path)` - Remove a folder and all its contents
- `ListObjects(ctx, path, recursive)` - List objects in a folder

#### URL Operations
- `GetPresignedURL(ctx, path, expiry)` - Generate a presigned GET URL
- `GetPresignedPutURL(ctx, path, expiry)` - Generate a presigned PUT URL
- `GetPublicURL(path)` - Generate a public URL (if configured)

#### Bucket Operations
- `BucketExists(ctx)` - Check if the configured bucket exists
- `ListBuckets(ctx)` - List all buckets
- `GetBucketLocation(ctx)` - Get bucket location

#### Raw Client Access
- `GetRawClient()` - Access the underlying MinIO client for advanced operations

### Advanced Features

#### Path Management
All paths are automatically managed with the configured `BaseDirPrefix`. This means:
- Input paths are relative to your base directory
- No need to manually prepend the base path
- Automatic path validation prevents traversal attacks

#### Error Handling
The library provides detailed error messages and proper error wrapping:

```go
object, err := client.GetObject(ctx, "nonexistent/file.txt", minio.GetObjectOptions{})
if err != nil {
    var errResponse minio.ErrorResponse
    if errors.As(err, &errResponse) {
        if errResponse.Code == "NoSuchKey" {
            fmt.Println("File does not exist")
        }
    }
}
```

#### Logging
The library uses structured logging. You can configure the log level through the `sloglog` package:

```go
import "github.com/aeternitas-infinita/sloglog"

// Set log level to debug for detailed operation logs
sloglog.SetLevel(slog.LevelDebug)
```

## Examples

### Uploading Files with Metadata

```go
opts := minio.PutObjectOptions{
    ContentType: "image/jpeg",
    UserMetadata: map[string]string{
        "photographer": "John Doe",
        "location":     "New York",
    },
}

uploadInfo, err := client.PutObject(ctx, "photos/sunset.jpg", file, fileSize, opts)
```

### Listing Objects with Pagination

```go
// List objects in the "documents" folder
for objectInfo := range client.ListObjects(ctx, "documents", false) {
    if objectInfo.Err != nil {
        log.Printf("Error listing objects: %v", objectInfo.Err)
        continue
    }
    fmt.Printf("Object: %s, Size: %d, Modified: %s\n", 
        objectInfo.Key, objectInfo.Size, objectInfo.LastModified)
}
```

### Working with Presigned URLs

```go
// Generate a presigned URL for uploading
putURL, err := client.GetPresignedPutURL(ctx, "uploads/document.pdf", time.Hour)
if err != nil {
    log.Fatal(err)
}

// Use the URL for direct upload from frontend
fmt.Printf("Upload your file to: %s\n", putURL.String())

// Generate a presigned URL for downloading
getURL, err := client.GetPresignedURL(ctx, "uploads/document.pdf", time.Hour*24)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Download link: %s\n", getURL.String())
```

### Using Raw MinIO Client

For advanced operations not covered by the extended client, you can access the underlying MinIO client:

```go
// Get the raw MinIO client
rawClient := client.GetRawClient()

// Use any native MinIO functionality directly
// Note: Operations on raw client don't benefit from automatic path prefix handling
buckets, err := rawClient.ListBuckets(ctx)
if err != nil {
    log.Fatal(err)
}

for _, bucket := range buckets {
    fmt.Printf("Bucket: %s, Created: %s\n", bucket.Name, bucket.CreationDate)
}

// Example: Use raw client for operations not yet wrapped
notification, err := rawClient.GetBucketNotification(ctx, "my-bucket")
if err != nil {
    log.Fatal(err)
}
```

### Working with Multiple Buckets

Each client instance is configured for a specific bucket. To work with multiple buckets, simply create multiple clients:

```go
// Client for user uploads
uploadsConfig := &miniox.Config{
    Endpoint:      "localhost:9000",
    AccessKey:     "your-access-key",
    SecretKey:     "your-secret-key",
    UseSSL:        false,
    BucketName:    "user-uploads",
    BaseDirPrefix: "uploads",
}
uploadsClient, err := miniox.New(uploadsConfig)

// Client for application data
dataConfig := &miniox.Config{
    Endpoint:      "localhost:9000",
    AccessKey:     "your-access-key",
    SecretKey:     "your-secret-key",
    UseSSL:        false,
    BucketName:    "app-data",
    BaseDirPrefix: "data",
}
dataClient, err := miniox.New(dataConfig)

// Use each client for its specific bucket
uploadInfo, err := uploadsClient.PutObject(ctx, "documents/file.pdf", reader, size, opts)
appData, err := dataClient.GetObject(ctx, "config/settings.json", minio.GetObjectOptions{})
```

This approach provides several benefits:
- **Isolation**: Each client handles its own bucket and path prefix
- **Configuration**: Different settings per bucket (SSL, endpoints, etc.)
- **Security**: Different credentials per bucket if needed
- **Performance**: Each client maintains its own connection pool
- **Clarity**: Code clearly shows which bucket is being used

## Best Practices

1. **Always use context**: Pass context for proper cancellation and timeouts
2. **Handle errors gracefully**: Check for specific MinIO error codes when needed
3. **Use appropriate content types**: Set ContentType in PutObjectOptions for better file handling
4. **Validate inputs**: The library validates paths, but validate business logic in your application
5. **Use presigned URLs for client uploads**: More secure than exposing credentials
6. **Set reasonable expiry times**: Don't make presigned URLs valid longer than necessary
7. **Multiple buckets**: Create separate client instances for different buckets - each client is lightweight and manages its own configuration
8. **Connection reuse**: Clients maintain their own connection pools, so reuse client instances when possible

## Error Handling

Common error scenarios and how to handle them:

```go
// Check if object exists before operating on it
_, err := client.StatObject(ctx, "path/to/file.txt", minio.StatObjectOptions{})
if err != nil {
    errResponse := minio.ToErrorResponse(err)
    if errResponse.Code == "NoSuchKey" {
        fmt.Println("File does not exist")
    } else {
        log.Printf("Error checking file: %v", err)
    }
}

// Handle network errors
object, err := client.GetObject(ctx, "large-file.bin", minio.GetObjectOptions{})
if err != nil {
    if strings.Contains(err.Error(), "connection") {
        fmt.Println("Network connection issue, please retry")
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Changelog

### v1.0.0
- Initial release with core functionality
- Support for object and folder operations
- Presigned URL generation
- Automatic path management
- Comprehensive error handling

## Support

If you encounter any issues or have questions, please open an issue on GitHub.
