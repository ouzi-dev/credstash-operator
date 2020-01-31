package credstash

import "github.com/ouzi-dev/credstash-operator/pkg/apis/credstash/v1alpha1"

type SecretGetter interface {
	GetCredstashSecretsForCredstashSecretDefs(credstashDefs []v1alpha1.CredstashSecretDef) (map[string][]byte, error)
}