package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	cr "github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
)

// @project photo-studio
// @created 10.08.2022

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
		  "arn:aws:s3:::%[1]s"
		],
		"Condition":{"StringLike":{"s3:prefix":"%[2]s"}}
    },
    {
      "Sid": "GetListObjectsBuckets",
      "Effect" : "Allow",
      "Action" : [
        "s3:GetObject"
      ],
	  "Resource": ["arn:aws:s3:::%[1]s/%[2]s/*"]
    }
  ]
}`

type UserMinioClient struct {
	UserID      string
	Bucket      string
	Endpoint    *url.URL
	Credentials *cr.Credentials
}

func CreateNewUserClient(userID string) (*UserMinioClient, error) {
	host := viper.GetString("minio.host")
	accessKey := viper.GetString("minio.accessKey")
	secretKey := viper.GetString("minio.secretKey")
	bucket := viper.GetString("minio.bucket")

	endpoint, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("CreateNewUserCredentia: Error parsing sts endpoint: %v ", err)
	}

	var stsOpts cr.STSAssumeRoleOptions
	stsOpts.AccessKey = accessKey
	stsOpts.SecretKey = secretKey

	policy := fmt.Sprintf(participantPolicy, bucket, userID)
	stsOpts.Policy = policy
	stsOpts.DurationSeconds = 300

	li, err := cr.NewSTSAssumeRole(host, stsOpts)
	if err != nil {
		return nil, fmt.Errorf("CreateNewUserClient: Error initializing STS Identity: %v ", err)
	}

	if _, err := li.Get(); err != nil {
		return nil, fmt.Errorf("CreateNewUserClient: Error retrieving STS credentials: %v ", err)
	}

	return &UserMinioClient{
		UserID:      userID,
		Bucket:      bucket,
		Endpoint:    endpoint,
		Credentials: li,
	}, nil
}

func (umc *UserMinioClient) GetUserFiles() ([]minio.ObjectInfo, error) {
	// todo: renew credentials if that expired
	opts := &minio.Options{
		Creds:  umc.Credentials,
		Secure: umc.Endpoint.Scheme == "https",
	}
	// Use generated credentials to authenticate with MinIO server
	minioClient, err := minio.New(umc.Endpoint.Host, opts)
	if err != nil {
		return nil, fmt.Errorf("GetUserFiles: Error initializing client: %w ", err)
	}

	log.Printf("GetUserFiles: Calling list objects on bucket named `%s` with temp creds", umc.Bucket)

	objCh := minioClient.ListObjects(context.Background(), umc.Bucket, minio.ListObjectsOptions{
		Prefix:    umc.UserID,
		Recursive: true,
	})
	result := make([]minio.ObjectInfo, 0, 0)
	for obj := range objCh {
		if obj.Err != nil {
			log.Errorf("GetUserFiles: Listing error: %v ", obj.Err)
			continue
		}
		result = append(result, obj)
	}
	return result, nil
}
