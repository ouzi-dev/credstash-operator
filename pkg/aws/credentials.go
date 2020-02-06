package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const awsSdkMaxRetries = 5

// Gets the aws session to use for looking up credstash secrets
func GetAwsSession(region string, awsAccessKeyID string, awsSecretAccessKey string) (*session.Session, error) {
	if awsAccessKeyID == "" || awsSecretAccessKey == "" {
		config := aws.Config{
			Region:     aws.String(region),
			MaxRetries: aws.Int(awsSdkMaxRetries),
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
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}

// Gets the aws session to use for looking up credstash secrets falling back to the environment config
func GetAwsSessionFromEnv() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}
