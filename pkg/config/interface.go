package config

type AwsConfig struct {
	Region             string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

type AwsConfigGetter interface {
	GetAwsConfig(secretName string, namespace string) (*AwsConfig, error)
}
