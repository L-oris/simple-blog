package bucketrepository

import (
	"fmt"
	"io/ioutil"

	"github.com/L-oris/yabb/logger"
	"github.com/google/go-cloud/blob"
)

// Config to create Repository
type Config struct {
	BucketName string
}

// Repository for GCP buckets; exported for DI declarations
type Repository struct {
	bucket *blob.Bucket
}

// New creates a new Repository
func New(config Config) (*Repository, error) {
	bucket, err := setupGCP(ctx, config.BucketName)
	if err != nil {
		return &Repository{}, err
	}

	return &Repository{bucket}, nil
}

// Write writes a file to bucket
func (c Repository) Write(fileName string, fileBytes []byte) error {
	defaultError := fmt.Errorf("cannot write %s to bucket", fileName)

	bucketWriter, err := c.bucket.NewWriter(ctx, fileName, nil)
	if err != nil {
		logger.Log.Errorf("create bucketWriter error: %s", err.Error())
		return defaultError
	}

	if _, err := bucketWriter.Write(fileBytes); err != nil {
		logger.Log.Errorf("write to bucket error: %s", err.Error())
		return defaultError
	}

	if err := bucketWriter.Close(); err != nil {
		logger.Log.Errorf("close bucket error: %s", err.Error())
		return defaultError
	}

	return nil
}

// Read reads a file from bucket
func (c Repository) Read(fileName string) ([]byte, error) {
	defaultError := fmt.Errorf("cannot read %s from bucket", fileName)

	bucketReader, err := c.bucket.NewReader(ctx, fileName)
	if err != nil {
		logger.Log.Errorf("cannot find file %s: %s", fileName, err.Error())
		return nil, defaultError
	}
	defer bucketReader.Close()

	newFile, err := ioutil.ReadAll(bucketReader)
	if err != nil {
		logger.Log.Fatalf("cannot create new file: %s", err.Error())
		return nil, defaultError
	}

	return newFile, nil
}

// Delete deletes a file from bucket
func (c Repository) Delete(fileName string) error {
	if err := c.bucket.Delete(ctx, fileName); err != nil {
		logger.Log.Errorf("delete file %s error: %s", fileName, err.Error())
		return fmt.Errorf("cannot delete %s from bucket", fileName)
	}

	return nil
}
