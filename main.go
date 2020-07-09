package main

import (
	"flag"
	"fmt"
	"time"
	"strconv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var storage_session *session.Session

func main() {

	server 	:= flag.String("server", "http://", "s3 storage endpoint")
	key 		:= flag.String("key", "key", "s3 access key")
	secret 	:= flag.String("secret", "secret", "s3 access secret")
	bucket 	:= flag.String("bucket", "library", "bucket for book storage")
	days 		:= flag.String("days", "0", "bucket for book storage")

	flag.Parse()

	storage_session = ConnectStorage(*server, *key, *secret)

	objects, err := ListObjects(*bucket, "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	days_int, err := strconv.Atoi(*days)
	if err != nil {
		fmt.Println("Days given not in acceptable format")
		return
	}

	LifecycleObjects(*bucket, days_int, objects)

}

// ConnectStorage
// Create a connection to the s3 server
func ConnectStorage(url, key, secret string) *session.Session {
	// Configure s3 remote
	storage_config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(url),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	// Return new s3 session for starting connections
	return session.New(storage_config)
}

// ListObjects
// Get the objects in a bucket
func ListObjects(bucket, prefix string) ([]*s3.Object, error) {

	// Connection to s3 server
	storage_connection := s3.New(storage_session)

	// Upload a new object
	objects, err := storage_connection.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	return objects.Contents, nil
}

// LifecycleObjects
// Delete objects older than given date
func LifecycleObjects(bucket string, days int, objects []*s3.Object) error {

	// Connection to s3 server
	storage_connection := s3.New(storage_session)

	cutoff := time.Now().AddDate(0, 0, -days)

	fmt.Printf("Deleting objects lat modified before %s\n\n", cutoff)

	for _, item := range objects {

		if item.LastModified.Before(cutoff) {
			fmt.Println("Object older than cutoff, deleting...")
			fmt.Println("  Name:         ", *item.Key)
			fmt.Println("  Last modified:", *item.LastModified)
			fmt.Println("  Size:         ", *item.Size)
			fmt.Println("  Storage class:", *item.StorageClass)
			fmt.Println("")

			_, err := storage_connection.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key: aws.String(*item.Key),
			})
			if err != nil {
				fmt.Printf("Unable to delete object %q from bucket %q, %v\n\n", *item.Key, bucket, err)
				continue
			}

			err = storage_connection.WaitUntilObjectNotExists(&s3.HeadObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(*item.Key),
			})
			if err != nil {
				fmt.Printf("Unable to delete object %s from bucket %s, %v\n\n", *item.Key, bucket, err)
				continue
			}

			fmt.Printf("Object %s successfully deleted\n\n\n", *item.Key)
		}
	}

	// Return no error
	return nil
}
