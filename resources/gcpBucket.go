package resources

import (
	"context"

	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/gcsblob"
	"github.com/google/go-cloud/gcp"
)

// SetupGCP sets up Google Cloud Bucket
func SetupGCP(ctx context.Context, bucket string) (*blob.Bucket, error) {
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