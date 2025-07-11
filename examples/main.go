// Example demonstrating basic usage of the MinIO Go Extended library
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
		BaseDirPrefix: "app-data",              // Optional: all operations will be relative to this path
		PublicURL:     "http://localhost:9000", // Optional: for generating public URLs
	}

	client, err := miniox.New(config)
	if err != nil {
		log.Fatal("Failed to initialize MinIO client:", err)
	}

	ctx := context.Background()

	// Example 1: Upload a file
	fmt.Println("=== Example 1: Upload a file ===")
	content := "Hello, World! This is a test file."
	uploadInfo, err := client.PutObject(ctx, "documents/hello.txt",
		strings.NewReader(content), int64(len(content)), minio.PutObjectOptions{
			ContentType: "text/plain",
		})
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
	} else {
		fmt.Printf("File uploaded successfully: %s (ETag: %s)\n", uploadInfo.Key, uploadInfo.ETag)
	}

	// Example 2: Download a file
	fmt.Println("\n=== Example 2: Download a file ===")
	object, err := client.GetObject(ctx, "documents/hello.txt", minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Failed to download file: %v", err)
	} else {
		defer object.Close()

		// Read the content
		buf := make([]byte, 1024)
		n, err := object.Read(buf)
		if err != nil && err.Error() != "EOF" {
			log.Printf("Error reading object: %v", err)
		} else {
			fmt.Printf("Downloaded content: %s\n", string(buf[:n]))
		}
	}

	// Example 3: Create a folder
	fmt.Println("\n=== Example 3: Create a folder ===")
	err = client.CreateFolder(ctx, "images")
	if err != nil {
		log.Printf("Failed to create folder: %v", err)
	} else {
		fmt.Println("Folder 'images' created successfully")
	}

	// Example 4: List objects
	fmt.Println("\n=== Example 4: List objects ===")
	for objectInfo := range client.ListObjects(ctx, "", true) {
		if objectInfo.Err != nil {
			log.Printf("Error listing objects: %v", objectInfo.Err)
			continue
		}
		fmt.Printf("Object: %s, Size: %d, Modified: %s\n",
			objectInfo.Key, objectInfo.Size, objectInfo.LastModified.Format(time.RFC3339))
	}

	// Example 5: Generate a presigned URL
	fmt.Println("\n=== Example 5: Generate presigned URL ===")
	presignedURL, err := client.GetPresignedURL(ctx, "documents/hello.txt", time.Hour)
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
	} else {
		fmt.Printf("Presigned URL (valid for 1 hour): %s\n", presignedURL.String())
	}

	// Example 6: Access raw MinIO client for advanced operations
	fmt.Println("\n=== Example 6: Raw MinIO client access ===")
	rawClient := client.GetRawClient()

	// Use raw client to list all buckets (not affected by BaseDirPrefix)
	buckets, err := rawClient.ListBuckets(ctx)
	if err != nil {
		log.Printf("Failed to list buckets: %v", err)
	} else {
		fmt.Printf("Available buckets:\n")
		for _, bucket := range buckets {
			fmt.Printf("  - %s (created: %s)\n", bucket.Name, bucket.CreationDate.Format(time.RFC3339))
		}
	}

	// Example 7: Use raw client for bucket-level operations
	fmt.Println("\n=== Example 7: Raw client bucket operations ===")
	bucketLocation, err := rawClient.GetBucketLocation(ctx, client.GetBucketName())
	if err != nil {
		log.Printf("Failed to get bucket location: %v", err)
	} else {
		fmt.Printf("Bucket '%s' is located in: %s\n", client.GetBucketName(), bucketLocation)
	}

	// Example 8: Multiple bucket usage (demonstration concept)
	fmt.Println("\n=== Example 8: Multiple bucket concept ===")
	fmt.Println("To work with multiple buckets, create multiple clients:")
	fmt.Println("  uploadsClient, _ := miniox.New(&miniox.Config{BucketName: \"uploads\", ...})")
	fmt.Println("  dataClient, _ := miniox.New(&miniox.Config{BucketName: \"app-data\", ...})")
	fmt.Println("Each client manages its own bucket and path prefix independently.")

	fmt.Println("\n=== All examples completed ===")
}
