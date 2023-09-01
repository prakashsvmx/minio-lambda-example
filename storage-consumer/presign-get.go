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

	//bucketName := "test-bucket"
	//objName := "1.txt"

	// Connect to the MinIO deployment
	s3Client, err := minio.New("localhost:22000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio123", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Set the Lambda function target using its ARN
	reqParams := make(url.Values)
	reqParams.Set("lambdaArn", "arn:minio:s3-object-lambda::myfunction:webhook")
	reqParams.Set("my-range", "10-20")
	// Generate a presigned url to access the original object

	presignedURL, err := s3Client.PresignedGetObject(context.Background(), "test-bucket", "1.txt", time.Duration(1000)*time.Second, reqParams)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest(http.MethodGet, presignedURL.String(), nil)
	if err != nil {
		fmt.Println("HTTP request to Presigned URL failed", err)
		return
	}

	//req.Header.Add("Range", "bytes=10-20")

	transport, err := minio.DefaultTransport(true)
	if err != nil {
		fmt.Println("HTTP request to Presigned URL failed", err)
		return
	}

	httpClient := &http.Client{
		// Setting a sensible time out of 30secs to wait for response
		// headers. Request is pro-actively canceled after 30secs
		// with no response.
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	/*fmt.Println("============= Req Headers")
	for name, values := range req.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
	fmt.Println("============= Req Headers")
	*/
	// rangeValue := "bytes=10-30"
	// req.Header.Add("Range", rangeValue)

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

	print("Done...")
	print(resp.ContentLength)

}
