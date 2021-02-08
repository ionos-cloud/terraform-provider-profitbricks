package profitbricks

import (
	"encoding/xml"
	"net/http"
)

type S3Buckets struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner S3BucketsOwner
	Buckets []S3Bucket `xml:"Buckets>Bucket"`
}

type S3BucketsOwner struct {
	ID string
	DisplayName string
}

type S3Bucket struct {
	Name string
	CreationDate string
}

func s3BucketsPath() string {
	return "/"
}

func (s3 *S3Client) ListS3Buckets() (*S3Buckets, error) {
	url := s3BucketsPath()
	ret := &S3Buckets{}
	err := s3.Get(url, ret, http.StatusOK)
	return ret, err

}
