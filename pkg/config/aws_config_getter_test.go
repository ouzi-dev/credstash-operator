package config

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

const (
	secretName = "credstash_config"
	secretNamespace = "credstash"
)

func TestK8sAwsConfigGetter_GetAwsConfig_ReturnsError_WhenSecretNotFound(t *testing.T) {
	s := scheme.Scheme

	// Objects to track in the fake client.
	objs := []runtime.Object{}

	// Create a fake client to mock API calls.
	client := fake.NewFakeClientWithScheme(s, objs...)

	subject := NewAwsSecretGetter(client)

	_, err := subject.GetAwsConfig(secretName, secretNamespace)

	assert.Error(t, err)
	assert.True(t,errors.IsNotFound(err))
}

func TestK8sAwsConfigGetter_GetAwsConfig_ReturnsError_WhenSecretIsInvalid(t *testing.T) {
	s := scheme.Scheme

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
			Namespace: secretNamespace,
		},
		Data: map[string][]byte{
				"AWS_ACCESS_KEY_ID": []byte("very"),
				"AWS_SECRET_ACCESS_KEY": []byte("secret"),
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		secret,
	}

	// Create a fake client to mock API calls.
	client := fake.NewFakeClientWithScheme(s, objs...)

	subject := NewAwsSecretGetter(client)

	_, err := subject.GetAwsConfig(secretName, secretNamespace)

	assert.Error(t, err)
	assert.False(t,errors.IsNotFound(err))
	assert.Equal(t, "secret credstash_config in namespace credstash does not contain valid config. Field AWS_REGION is missing", err.Error())
}

func TestK8sAwsConfigGetter_GetAwsConfig_ReturnsConfig_WhenSecretIsValid(t *testing.T) {
	s := scheme.Scheme

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
			Namespace: secretNamespace,
		},
		Data: map[string][]byte{
			"AWS_REGION": []byte("region"),
			"AWS_ACCESS_KEY_ID": []byte("very"),
			"AWS_SECRET_ACCESS_KEY": []byte("secret"),
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		secret,
	}

	// Create a fake client to mock API calls.
	client := fake.NewFakeClientWithScheme(s, objs...)

	subject := NewAwsSecretGetter(client)

	expectedConfig := &AwsConfig{awsSecretAccessKey: "secret", awsAccessKeyID: "very", region: "region"}
	config, err := subject.GetAwsConfig(secretName, secretNamespace)

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)
}