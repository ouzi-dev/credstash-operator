package credstash

//nolint
//go:generate mockgen -destination=../mocks/mock_credential_helper.go -package=mocks github.com/singularityconsulting/hoolio/internal/credentials EnvironmentCredentialHelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ouzi-dev/credstash-operator/pkg/apis/credstash/v1alpha1"
	"github.com/versent/unicreds"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type secretGetter struct {
}

const (
	defaultCredstashTable = "credential-store"
	credstashVersionLength = 19
)

var log = logf.Log.WithName("credstashsecret_getter")

func NewSecretGetter(awsSession *session.Session) SecretGetter {
	unicreds.SetKMSConfig(awsSession.Config)
	unicreds.SetDynamoDBConfig(awsSession.Config)
	return &secretGetter{}
}

func (h *secretGetter) GetCredstashSecretsForCredstashSecretDefs(credstashDefs []v1alpha1.CredstashSecretDef) (map[string][]byte, error) {
	ecryptionContext := unicreds.NewEncryptionContextValue()
	secretsMap := map[string][]byte{}
	for _,v := range credstashDefs {
		table := v.Table
		if table == "" {
			table = defaultCredstashTable
		}

		if v.Version == "" {
			creds, err := unicreds.GetHighestVersionSecret(aws.String(table), v.Key, ecryptionContext)
			if err != nil {
				log.Error(err, "Failed fetching secret from credstash",
					"Secret.Key", v.Key, "Secret.Version", "latest", "Secret.Table", table)
				return nil, err
			}
			secretsMap[v.Key] = []byte(creds.Secret)
		} else {
			formattedVersion, err := formatCredstashVersion(v.Version)
			if err != nil {
				log.Error(err, "Failed formatting secret version",
					"Secret.Key", v.Key, "Secret.Version", v.Version, "Secret.Table", table)
				return nil, err
			}

			creds, err := unicreds.GetSecret(aws.String(table), v.Key, formattedVersion, ecryptionContext)
			if err != nil {
				log.Error(err, "Failed fetching secret from credstash",
					"Secret.Key", v.Key, "Secret.Version", formattedVersion, "Secret.Table", table)
				return nil, err
			}
			secretsMap[v.Key] = []byte(creds.Secret)
		}
	}

	return secretsMap, nil
}
