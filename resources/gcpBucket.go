package resources

import (
	"context"

	"github.com/L-oris/yabb/logger"
	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/gcsblob"
	"github.com/google/go-cloud/gcp"
)

var yabbBucket *blob.Bucket

// GetYabbBucket gets pointer to 'yabb' bucket
func GetYabbBucket() (*blob.Bucket, error) {
	if yabbBucket != nil {
		logger.Log.Debug("found existing bucket")
		return yabbBucket, nil
	}

	logger.Log.Debug("setting up bucket connection")
	var err error // avoid variable shadowing
	yabbBucket, err = setupGCP(CTX, "yabb")
	return yabbBucket, err
}

// setupGCP sets up connection to Google Cloud Bucket
func setupGCP(ctx context.Context, bucket string) (*blob.Bucket, error) {
	credentials, err := gcp.DefaultCredentials(CTX)
	if err != nil {
		return nil, err
	}

	client, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(credentials))
	if err != nil {
		return nil, err
	}

	return gcsblob.OpenBucket(CTX, bucket, client)
}
