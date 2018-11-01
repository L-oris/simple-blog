package bucketrepository

import (
	"context"
	"fmt"

	"github.com/L-oris/yabb/logger"
	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/gcsblob"
	"github.com/google/go-cloud/gcp"
)

// setupGCP sets up connection to Google Cloud Bucket
func setupGCP(ctx context.Context, bucket string) (*blob.Bucket, error) {
	credentials, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		err = fmt.Errorf("cannot get default credentials: %s", err.Error())
		logger.Log.Error(err.Error())
		return nil, err
	}

	client, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(credentials))
	if err != nil {
		err = fmt.Errorf("cannot create HTTPClient: %s", err.Error())
		logger.Log.Error(err.Error())
		return nil, err
	}

	return gcsblob.OpenBucket(ctx, bucket, client)
}
