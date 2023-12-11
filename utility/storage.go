package utility

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func UploadFileToGCS(file multipart.File, filename string, bucketName string, credentials string) (string, error) {
	ctx := context.Background()

	// Create a new Google Cloud Storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentials))
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Get the bucket handle
	bucket := client.Bucket(bucketName)

	// Create an object handle
	object := bucket.Object(filename)

	// Create a new Writer
	wc := object.NewWriter(ctx)

	// Copy the file data to the object in GCS
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}

	// Close the Writer to flush the data to GCS
	if err := wc.Close(); err != nil {
		return "", err
	}

	// Generate the public URL for the uploaded file
	// Set proper ACL rules to make the object publicly accessible
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	// Construct the public URL
	gcsURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, attrs.Name)

	return gcsURL, nil
}
