package credstashsecret

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ouzi-dev/credstash-operator/pkg/mocks"
	"github.com/stretchr/testify/assert"
	errors2 "k8s.io/apimachinery/pkg/api/errors"

	credstashv1alpha1 "github.com/ouzi-dev/credstash-operator/pkg/apis/credstash/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	errorString    = "an error has occurred"
	name           = "credstashCR"
	namespace      = "credstash"
	credstashKey   = "key1"
	credstashValue = "value1"
	secretName     = "specialName"
)

type testReconcileItem struct {
	testName             string
	customResource       *credstashv1alpha1.CredstashSecret
	existsingSecret      *corev1.Secret
	credstashError       error
	expectedResultSecret *corev1.Secret
}

var (
	credstashGetterReturn = map[string][]byte{
		credstashKey: []byte(credstashValue),
	}
)

var tests = []testReconcileItem{
	{
		testName: "Credstash error",
		customResource: &credstashv1alpha1.CredstashSecret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: credstashv1alpha1.CredstashSecretSpec{
				Secrets: []credstashv1alpha1.CredstashSecretDef{
					{
						Key: credstashKey,
					},
				},
			},
		},
		existsingSecret:      nil,
		credstashError:       errors.New(errorString),
		expectedResultSecret: nil,
	},
	{
		testName: "No existing secret",
		customResource: &credstashv1alpha1.CredstashSecret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: credstashv1alpha1.CredstashSecretSpec{
				Secrets: []credstashv1alpha1.CredstashSecretDef{
					{
						Key: credstashKey,
					},
				},
			},
		},
		existsingSecret: nil,
		credstashError:  nil,
		expectedResultSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
			Type: corev1.SecretTypeOpaque,
		},
	},
	{
		testName: "Existing different data secret",
		customResource: &credstashv1alpha1.CredstashSecret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: credstashv1alpha1.CredstashSecretSpec{
				Secrets: []credstashv1alpha1.CredstashSecretDef{
					{
						Key: credstashKey,
					},
				},
			},
			Status: credstashv1alpha1.CredstashSecretStatus{
				SecretName: name,
			},
		},
		existsingSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"differentKey": []byte("differentValue"),
			},
			Type: corev1.SecretTypeOpaque,
		},
		credstashError: nil,
		expectedResultSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
			Type: corev1.SecretTypeOpaque,
		},
	},
	{
		testName: "Existing identical data secret",
		customResource: &credstashv1alpha1.CredstashSecret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: credstashv1alpha1.CredstashSecretSpec{
				Secrets: []credstashv1alpha1.CredstashSecretDef{
					{
						Key: credstashKey,
					},
				},
			},
			Status: credstashv1alpha1.CredstashSecretStatus{
				SecretName: name,
			},
		},
		existsingSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
		},
		credstashError: nil,
		expectedResultSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
			Type: corev1.SecretTypeOpaque,
		},
	},
	{
		testName: "Name specified and no existing secret",
		customResource: &credstashv1alpha1.CredstashSecret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: credstashv1alpha1.CredstashSecretSpec{
				SecretName: secretName,
				Secrets: []credstashv1alpha1.CredstashSecretDef{
					{
						Key: credstashKey,
					},
				},
			},
		},
		existsingSecret: nil,
		credstashError:  nil,
		expectedResultSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
			Type: corev1.SecretTypeOpaque,
		},
	},
	{
		testName: "Name specified and existing secret with same name and data",
		customResource: &credstashv1alpha1.CredstashSecret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: credstashv1alpha1.CredstashSecretSpec{
				SecretName: secretName,
				Secrets: []credstashv1alpha1.CredstashSecretDef{
					{
						Key: credstashKey,
					},
				},
			},
			Status: credstashv1alpha1.CredstashSecretStatus{
				SecretName: secretName,
			},
		},
		existsingSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
		},
		credstashError: nil,
		expectedResultSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
			Type: corev1.SecretTypeOpaque,
		},
	},
	{
		testName: "Name specified and existing secret with different name",
		customResource: &credstashv1alpha1.CredstashSecret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: credstashv1alpha1.CredstashSecretSpec{
				SecretName: secretName,
				Secrets: []credstashv1alpha1.CredstashSecretDef{
					{
						Key: credstashKey,
					},
				},
			},
			Status: credstashv1alpha1.CredstashSecretStatus{
				SecretName: secretName,
			},
		},
		existsingSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "different name",
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
		},
		credstashError: nil,
		expectedResultSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
			Type: corev1.SecretTypeOpaque,
		},
	},
	{
		testName: "Custom secret type",
		customResource: &credstashv1alpha1.CredstashSecret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: credstashv1alpha1.CredstashSecretSpec{
				SecretName: secretName,
				Secrets: []credstashv1alpha1.CredstashSecretDef{
					{
						Key: credstashKey,
					},
				},
				SecretType: corev1.SecretTypeDockerConfigJson,
			},
		},
		existsingSecret: nil,
		credstashError:  nil,
		expectedResultSecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
			},
			Data: credstashGetterReturn,
			Type: corev1.SecretTypeDockerConfigJson,
		},
	},
}

//nolint funlen
func TestReconcileCredstashSecret_Reconcile(t *testing.T) {
	for _, testData := range tests {
		t.Run(testData.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCredstashSecretGetter := mocks.NewMockSecretGetter(ctrl)

			// Register operator types with the runtime scheme.
			s := scheme.Scheme
			s.AddKnownTypes(credstashv1alpha1.SchemeGroupVersion, testData.customResource)

			// Objects to track in the fake client.
			objs := []runtime.Object{
				testData.customResource,
			}

			if testData.existsingSecret != nil {
				objs = append(objs, testData.existsingSecret)
			}

			// Create a fake client to mock API calls.
			cl := fake.NewFakeClientWithScheme(s, objs...)
			// Create a ReconcileCredstashSecret object with the scheme and fake client.
			r := &ReconcileCredstashSecret{client: cl, scheme: s, credstashSecretGetter: mockCredstashSecretGetter}

			// Mock request to simulate Reconcile() being called on an event for a
			// watched resource .
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      name,
					Namespace: namespace,
				},
			}

			mockCredstashSecretGetter.
				EXPECT().GetCredstashSecretsForCredstashSecretDefs(testData.customResource.Spec.Secrets).
				Return(credstashGetterReturn, testData.credstashError).
				Times(1)

			_, err := r.Reconcile(req)

			assert.Equal(t, testData.credstashError, err)

			var expectedSecretName string
			if testData.expectedResultSecret == nil {
				expectedSecretName = testData.customResource.Name
			} else {
				expectedSecretName = testData.expectedResultSecret.Name
			}

			// Check if Secret has been created and has the correct data
			secret := &corev1.Secret{}
			err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedSecretName, Namespace: namespace}, secret)

			if testData.expectedResultSecret == nil {
				assert.Error(t, err)
				assert.True(t, errors2.IsNotFound(err))
			} else {
				assert.Equal(t, testData.expectedResultSecret.Data, secret.Data)
				assert.Equal(t, testData.expectedResultSecret.Name, secret.Name)
				assert.Equal(t, testData.expectedResultSecret.Type, secret.Type)

				updatedCR := &credstashv1alpha1.CredstashSecret{}
				err = cl.Get(context.TODO(), req.NamespacedName, updatedCR)
				assert.NoError(t, err)
				assert.Equal(t, testData.expectedResultSecret.Name, updatedCR.Status.SecretName)
			}
		})
	}
}
