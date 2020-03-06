package credstash

import "github.com/ouzi-dev/credstash-operator/pkg/apis/credstash/v1alpha1"

type SecretGetter interface {
	GetCredstashSecretsForCredstashSecret(credstashSecret *v1alpha1.CredstashSecret) (map[string][]byte, error)
}
