package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"strings"
)

func main() {

	bucketName := "test-bucket"
	objName := "1.txt"

	// Connect to the MinIO deployment
	s3Client, err := minio.New("localhost:22000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio123", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	//Just try to see if it works. Though it is not fine.

	// Set the Lambda function target using its ARN
	// reqParams := make(url.Values)
	// reqParams.Set("lambdaArn", "arn:minio:s3-object-lambda::myfunction:webhook")

	// Generate a presigned url to access the original object

	options := &minio.GetObjectOptions{}
	// options.SetReqParam("lambdaArn", "arn:minio:s3-object-lambda::myfunction:webhook") // this does not work.
	err = options.SetRange(10, 20) // this works.
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	reader, err := s3Client.GetObject(context.Background(), bucketName, objName, *options)
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	buf := new(strings.Builder)
	n, err := io.Copy(buf, reader)
	// check errors
	fmt.Println(buf.String())
	fmt.Println("Done...", n)

}
