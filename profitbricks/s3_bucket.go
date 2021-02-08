package profitbricks

import "net/http"

type S3Buckets struct {

}

func s3BucketsPath() string {
	return "/"
}

func (s3 *S3Client) ListBuckets() (*S3Buckets, error) {
	url := s3BucketsPath()
	ret := &S3Buckets{}
	err := s3.Get(url, ret, http.StatusOK)
	return ret, err

}
