package gcpstorage

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type clientWithContext struct {
	client *storage.Client
	ctx    context.Context
}

type StorageAsset struct {
	BucketName      string
	ObjectName      string
	CredentialsFile string
}

func clientFromCredentials(credentialsFile string) (clientWithContext, error) {
	ctx := context.Background()
	credentialsFile = os.ExpandEnv(strings.Replace(credentialsFile, "~", "$HOME", 1))
	opt := option.WithCredentialsFile(credentialsFile)
	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		return clientWithContext{}, fmt.Errorf("error initializing storage client: %s", err)
	}
	return clientWithContext{
		client: client,
		ctx:    ctx,
	}, nil
}

func getClient(credentialsFile string) (clientWithContext, error) {
	if credentialsFile == "" {
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			return clientWithContext{}, fmt.Errorf("error initializing storage client: %s", err)
		}
		return clientWithContext{
			client: client,
			ctx:    ctx,
		}, nil
	}
	return clientFromCredentials(credentialsFile)
}

func Upload(storageAsset StorageAsset, content []byte) (err error) {
	clientWithContext, err := getClient(storageAsset.CredentialsFile)
	if err != nil {
		return fmt.Errorf("error getting storage client: %s", err)
	}
	client := clientWithContext.client
	ctx := clientWithContext.ctx

	bucketName := storageAsset.BucketName
	objectName := storageAsset.ObjectName

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	w := object.NewWriter(ctx)
	buffer := bytes.NewBuffer(content)

	_, err = io.Copy(w, buffer)
	if err != nil {
		return fmt.Errorf("error writing to object %s in bucket %s: %s", objectName, bucketName, err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("error closing handle for object %s in bucket %s: %s", objectName, bucketName, err)
	}

	return
}

func Download(storageAsset StorageAsset) (content []byte, err error) {
	clientWithContext, err := getClient(storageAsset.CredentialsFile)
	if err != nil {
		return content, fmt.Errorf("error getting storage client: %s", err)
	}
	client := clientWithContext.client
	ctx := clientWithContext.ctx

	bucketName := storageAsset.BucketName
	objectName := storageAsset.ObjectName

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	reader, err := object.NewReader(ctx)
	if err != nil {
		return content, fmt.Errorf("error getting a reader for object %s in bucket %s: %v", objectName, bucketName, err)
	}

	defer func(r io.ReadCloser) {
		err := r.Close()
		if err != nil {
			log.Printf("error closing reader for object %s in bucket %s", objectName, bucketName)
		}
	}(reader)

	content, err = ioutil.ReadAll(reader)
	if err != nil {
		return content, fmt.Errorf("error reading object %s in bucket %s: %s", objectName, bucketName, err)
	}
	return
}
