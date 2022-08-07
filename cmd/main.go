package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	cr "github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"os"
)

// @project photostudio
// @created 27.07.2022

func init() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

}

func getPolicy(fn string) (string, error) {
	open, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer open.Close()
	buf := make([]byte, 1024)
	n, err := open.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

var participantPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "GetObjectsInFolder",
      "Effect" : "Allow",
      "Action" : [
        "s3:ListBucket"
      ],
		"Resource" : [
		  "arn:aws:s3:::photostudio"
		],
		"Condition":{"StringLike":{"s3:prefix":"%[1]s"}}
    },
    {
      "Sid": "GetListObjectsBuckets",
      "Effect" : "Allow",
      "Action" : [
        "s3:GetObject"
      ],
	  "Resource": ["arn:aws:s3:::photostudio/%[1]s/*"]
    }
  ]
}`

func main() {
	host := viper.GetString("s3.host")
	accessKey := viper.GetString("s3.accessKey")
	secretKey := viper.GetString("s3.secretKey")
	log.Info("host: ", host)

	var stsOpts cr.STSAssumeRoleOptions
	stsOpts.AccessKey = accessKey
	stsOpts.SecretKey = secretKey

	id := "Sibay"
	policy := fmt.Sprintf(participantPolicy, id)
	stsOpts.Policy = policy
	stsOpts.DurationSeconds = 300

	li, err := cr.NewSTSAssumeRole(host, stsOpts)
	if err != nil {
		log.Fatalf("Error initializing STS Identity: %v", err)
	}

	stsEndpointURL, err := url.Parse(host)
	if err != nil {
		log.Fatalf("Error parsing sts endpoint: %v", err)
	}

	opts := &minio.Options{
		Creds:  li,
		Secure: stsEndpointURL.Scheme == "https",
	}

	if _, err := li.Get(); err != nil {
		log.Fatalf("Error retrieving STS credentials: %v", err)
	}

	//fmt.Println("Only displaying credentials:")
	//fmt.Println("AccessKeyID:", v.AccessKeyID)
	//fmt.Println("SecretAccessKey:", v.SecretAccessKey)
	//fmt.Println("SessionToken:", v.SessionToken)

	// Use generated credentials to authenticate with MinIO server
	minioClient, err := minio.New(stsEndpointURL.Host, opts)
	if err != nil {
		log.Fatalf("Error initializing client: %v", err)
	}

	bucket := "photostudio"
	fmt.Printf("\nCalling list objects on bucket named `%s` with temp creds:\n===\n", bucket)
	objCh := minioClient.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
		Prefix:    id,
		Recursive: true,
	})
	for obj := range objCh {
		if obj.Err != nil {
			log.Fatalf("Listing error: %v", obj.Err)
		}
		fmt.Printf("Key: %s\nSize: %d\nLast Modified: %s\n===\n", obj.Key, obj.Size, obj.LastModified)
	}
}
