package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Gets the aws session to use for looking up credstash secrets
func GetAwsSession(region string, awsAccessKeyId string, awsSecretAccessKey string) (*session.Session, error){

	if awsAccessKeyId == "" || awsSecretAccessKey == "" {
		config := aws.Config{
			Region:     aws.String(region),
			MaxRetries: aws.Int(5),
		}

		sess, err := session.NewSessionWithOptions(session.Options{
			Config:            config,
			SharedConfigState: session.SharedConfigEnable,
		})

		if err != nil {
			return nil, err
		}

		return sess, nil

	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyId, awsSecretAccessKey, ""),
	})

	if err != nil {
		return nil, err
	}

	return sess, nil

}