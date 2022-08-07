package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
      "Sid": "GetListObjectsBuckets",
      "Effect" : "Allow",
      "Action" : [
        "s3:GetObject"
      ],
        "Resource" : [
          "arn:aws:s3:::photostudio/%s"
        ]
    }
  ]
}`

func main() {
	host := viper.GetString("s3.host")
	accessKey := viper.GetString("s3.accessKey")
	secretKey := viper.GetString("s3.secretKey")
	log.Info("host: ", host)

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String(host),
		Region:      aws.String("minio"),
	})

	if err != nil {
		fmt.Println("NewSession Error", err)
		return
	}
	//participantPolicy, err := getPolicy("policy/participant.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	id := "Sibay"
	// Create a STS client
	svc := sts.New(sess)
	role, err := svc.AssumeRole(&sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(3600),
		Policy:          aws.String(fmt.Sprintf(participantPolicy, id)),
		RoleArn:         aws.String("arn:x:ignored:by:minio:"), // arn:partition:service:region:account-id:resource-id
		RoleSessionName: aws.String("ignored-by-minio"),
		SourceIdentity:  aws.String("photostudio"),
	})
	if err != nil {
		log.Fatal("AssumeRole: ", err)
	}
	fmt.Println("Role: ", *role)
	//
	//mainS3svc := s3.New(sess)
	//if _, err := mainS3svc.PutObject(&s3.PutObjectInput{
	//	ACL:    aws.String("public-read"),
	//	Body:   nil,
	//	Bucket: aws.String("photostudio"),
	//	Key:    aws.String("photostudio/" + id + "/"),
	//}); err != nil {
	//	log.Fatal("PutObject: s", err)
	//}
	userSess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(*role.Credentials.AccessKeyId, *role.Credentials.SecretAccessKey, *role.Credentials.SessionToken),
		Endpoint:    aws.String(host),
		Region:      aws.String("minio"),
	})
	if err != nil {
		log.Fatal("NewSession: ", err)
	}
	userS3svc := s3.New(userSess)
	getObject, err := userS3svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("photostudio"),
	})
	//buckets, err := userS3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Error("GetObject: ", err)
	}
	fmt.Println("getObject: ", getObject)

	accessKeyInfo, err := svc.GetAccessKeyInfo(&sts.GetAccessKeyInfoInput{
		AccessKeyId: role.Credentials.AccessKeyId,
	})
	if err != nil {
		log.Error("GetAccessKeyInfo: ", err)
	}
	fmt.Println(accessKeyInfo)
}
