package config

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	awsRegionField = "AWS_REGION"
	awsAccessKeyIDField = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKeyField = "AWS_SECRET_ACCESS_KEY"
	validationError = "secret %s in namespace %s does not contain valid config. Field %s is missing"
)

type k8sAwsConfigGetter struct {
	client client.Client
}

func NewAwsSecretGetter(client client.Client) AwsConfigGetter {
	return &k8sAwsConfigGetter{client}
}

func (k *k8sAwsConfigGetter) GetAwsConfig(secretName string, namespace string) (*AwsConfig, error) {
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      secretName,
	}
	// Fetch the Secret
	secret := &v1.Secret{}
	err := k.client.Get(context.TODO(), namespacedName, secret)
	if err != nil {
		return nil, err
	}

	err = validateAwsConfigSecret(secret)
	if err != nil {
		return nil, err
	}

	return &AwsConfig{
		Region: string(secret.Data[awsRegionField]),
		AwsAccessKeyID: string(secret.Data[awsAccessKeyIDField]),
		AwsSecretAccessKey: string(secret.Data[awsSecretAccessKeyField]),

	}, nil
}

func validateAwsConfigSecret(secret *v1.Secret) error {
	if string(secret.Data[awsRegionField]) == "" {
		return fmt.Errorf(validationError, secret.Name, secret.Namespace, awsRegionField)
	}

	if string(secret.Data[awsAccessKeyIDField]) == "" {
		return fmt.Errorf(validationError, secret.Name, secret.Namespace, awsAccessKeyIDField)
	}

	if string(secret.Data[awsSecretAccessKeyField]) == "" {
		return fmt.Errorf(validationError, secret.Name, secret.Namespace, awsSecretAccessKeyField)
	}

	return nil
}


