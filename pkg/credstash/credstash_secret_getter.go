package credstash

//nolint
//go:generate mockgen -destination=../mocks/mock_credential_helper.go -package=mocks github.com/singularityconsulting/hoolio/internal/credentials EnvironmentCredentialHelper

import (
	"github.com/ouzi-dev/credstash-operator/pkg/apis/credstash/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/versent/unicreds"
)

type CredstashSecretGetter struct {
}

func NewHelper(awsSession *session.Session) *CredstashSecretGetter {
	unicreds.SetKMSConfig(awsSession.Config)
	unicreds.SetDynamoDBConfig(awsSession.Config)
	return &CredstashSecretGetter{}
}

func (h *CredstashSecretGetter) GetCredstashSecretsForCredstashSecretDefs(credstashDefs []v1alpha1.CredstashSecretDef) map[string]string {
	ecryptionContext := unicreds.NewEncryptionContextValue()
	secretsMap := map[string]string{}
	for _,v := range credstashDefs {
		//TODO validate table is correct and/or use defaults
		if v.Version == "" {
			creds, err := unicreds.GetHighestVersionSecret(aws.String(v.Table), v.Key, ecryptionContext)
			if err != nil {
				//TODO log error here
				continue
			}
			secretsMap[v.Key] = creds.Secret
		} else {
			creds, err := unicreds.GetSecret(aws.String(v.Table), v.Key, v.Version, ecryptionContext)
			if err != nil {
				//TODO log error here
				continue
			}
			secretsMap[v.Key] = creds.Secret
		}
	}

	return secretsMap
}
