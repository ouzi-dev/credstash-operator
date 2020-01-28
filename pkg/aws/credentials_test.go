package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAwsSession(t *testing.T) {
	sess, err := GetAwsSession("us-west-2", "AKIAVWXUA4X6HCCHYC77", "usdmqO85cKpL7TzkOX4kZ7i6z5niMuntLU6a1ACC")

//	sess.Config.Credentials.Get()

	assert.NoError(t, err)

	s3svc := s3.New(sess)

	listBucketsInput := &s3.ListBucketsInput{}

	out, err := s3svc.ListBuckets(listBucketsInput)

	assert.NoError(t, err)

	fmt.Print(out)
}