package profitbricks

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

type S3Buckets struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner   S3BucketOwner
	Buckets []S3BucketsItem `xml:"Buckets>Bucket"`
}

type S3BucketOwner struct {
	ID string
	DisplayName string
}

type S3BucketsItem struct {
	Name string
	CreationDate string
}

type S3Bucket struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	Name    string
	Contents []S3BucketContent
}

type S3BucketContent struct {
	Key          string
	LastModified string
	StorageClass string
	Size         int64
	ETag         string
	Owner        S3BucketOwner
}

func s3BucketsPath() string {
	return "/"
}

func (s3 S3Client) ListS3Buckets() (*S3Buckets, error) {
	url := s3BucketsPath()
	ret := &S3Buckets{}
	err := s3.Get(url, ret, http.StatusOK)
	return ret, err

}

func (s3 S3Client) GetBucketHost(name string) (string, error) {
	parts := strings.Split(s3.HostURL, "//")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid host: %s", s3.HostURL)
	}

	return fmt.Sprintf("%s//%s.%s", parts[0], name, parts[1]), nil
}

func (s3 S3Client) GetBucket(name string) (*S3Bucket, error) {

	bucketHost, err := s3.GetBucketHost(name)
	if err != nil {
		return nil, err
	}
	s3.HostURL = bucketHost
	ret := &S3Bucket{}
	err = s3.Get("/", ret, http.StatusOK)

	return ret, err
}
