package credstash

//nolint
//go:generate mockgen -destination=../mocks/mock_credstash_secret_getter.go -package=mocks github.com/ouzi-dev/credstash-operator/pkg/credstash SecretGetter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ouzi-dev/credstash-operator/pkg/apis/credstash/v1alpha1"
	"github.com/ouzi-dev/credstash-operator/pkg/event"
	"github.com/versent/unicreds"
	"k8s.io/client-go/tools/record"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type secretGetter struct {
	eventRecorder record.EventRecorder
}

const (
	defaultCredstashTable  = "credential-store"
	credstashVersionLength = 19
)

var log = logf.Log.WithName("credstashsecret_getter")

func NewSecretGetter(awsSession *session.Session, eventRecorder record.EventRecorder) SecretGetter {
	unicreds.SetKMSConfig(awsSession.Config)
	unicreds.SetDynamoDBConfig(awsSession.Config)

	return &secretGetter{
		eventRecorder: eventRecorder,
	}
}

func (h *secretGetter) GetCredstashSecretsForCredstashSecret(
	credstashSecret *v1alpha1.CredstashSecret) (map[string][]byte, error) {
	ecryptionContext := unicreds.NewEncryptionContextValue()
	secretsMap := map[string][]byte{}

	for _, v := range credstashSecret.Spec.Secrets {
		table := v.Table
		if table == "" {
			table = defaultCredstashTable
		}

		mapKey := v.Name
		if mapKey == "" {
			mapKey = v.Key
		}

		if v.Version == "" {
			creds, err := unicreds.GetHighestVersionSecret(aws.String(table), v.Key, ecryptionContext)
			if err != nil {
				h.eventRecorder.Eventf(
					credstashSecret, event.TypeWarning,
					event.ReasonErrFetchingCredstashSecret, event.MessageFailedFetchingCredstashSecret,
					v.Key, "latest", table, err.Error())
				log.Error(err, "Failed fetching secret from credstash",
					"Secret.Key", v.Key, "Secret.Version", "latest", "Secret.Table", table)

				return nil, err
			}

			secretsMap[mapKey] = []byte(creds.Secret)
		} else {
			formattedVersion, err := formatCredstashVersion(v.Version)
			if err != nil {
				h.eventRecorder.Eventf(
					credstashSecret, event.TypeWarning,
					event.ReasonErrFetchingCredstashSecret, event.MessageFailedFetchingCredstashSecret,
					v.Key, v.Version, table, err.Error())
				log.Error(err, "Failed formatting secret version",
					"Secret.Key", v.Key, "Secret.Version", v.Version, "Secret.Table", table)
				return nil, err
			}

			creds, err := unicreds.GetSecret(aws.String(table), v.Key, formattedVersion, ecryptionContext)
			if err != nil {
				h.eventRecorder.Eventf(
					credstashSecret, event.TypeWarning,
					event.ReasonErrFetchingCredstashSecret, event.MessageFailedFetchingCredstashSecret,
					v.Key, formattedVersion, table, err.Error())
				log.Error(err, "Failed fetching secret from credstash",
					"Secret.Key", v.Key, "Secret.Version", formattedVersion, "Secret.Table", table)
				return nil, err
			}

			secretsMap[mapKey] = []byte(creds.Secret)
		}
	}

	return secretsMap, nil
}
