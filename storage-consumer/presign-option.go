package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {

	// Connect to the MinIO deployment
	s3Client, err := minio.New("localhost:22000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio123", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	rangeValue := "bytes=10-30"

	// Set the Lambda function target using its ARN
	reqParams := make(url.Values)
	//**************** Does not work ***************************
	// reqParams.Set("lambdaArn", "arn:minio:s3-object-lambda::myfunction:webhook")

	bucketName := "test-bucket"
	objName := "1.txt"

	presignedURL, err := s3Client.Presign(context.Background(), "GET", bucketName, objName, 3600*time.Second, reqParams)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, presignedURL.String(), nil)
	if err != nil {
		fmt.Println("HTTP request to Presigned URL failed", err)
		return
	}

	transport, err := minio.DefaultTransport(true)
	if err != nil {
		fmt.Println("HTTP request to Presigned URL failed", err)
		return
	}

	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
	//req.Header.Add("lambdaArn", "arn:minio:s3-object-lambda::myfunction:webhook")
	req.Header.Add("Range", rangeValue)

	/*q := req.URL.Query()                                                 // Get a copy of the query values.
	q.Add("lambdaArn", "arn:minio:s3-object-lambda::myfunction:webhook") // Add a new value to the set.
	req.URL.RawQuery = q.Encode()
	*/
	fmt.Println("============= Req Headers")
	for name, values := range req.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
	fmt.Println("============= Req Headers")

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("HTTP request to Presigned URL failed", err)

		return
	}

	fmt.Println("============= Resp Headers")

	for name, values := range resp.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
	fmt.Println("============= Resp Headers")

	fmt.Println("====================Response=============")

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))

	fmt.Println("Done...")
	fmt.Println(resp.ContentLength)
	// Print the URL to stdout
	//fmt.Println(presignedURL)
}
