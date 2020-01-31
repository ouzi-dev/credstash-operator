package credstashsecret

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ouzi-dev/credstash-operator/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"testing"

	credstashv1alpha1 "github.com/ouzi-dev/credstash-operator/pkg/apis/credstash/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	errorString = "an error has occured"
)

func TestCreatingSecretWhenCredstashError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredstashSecretGetter := mocks.NewMockSecretGetter(ctrl)
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	var (
		name            = "credstashCR"
		namespace       = "credstash"
	)

	credstashSecretDefs := []credstashv1alpha1.CredstashSecretDef{
		{
			Key: "I do not exist",
		},
	}

	// A credstashCR resource with metadata and spec.
	credstashCR := &credstashv1alpha1.CredstashSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: credstashv1alpha1.CredstashSecretSpec{
			Secrets: credstashSecretDefs,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		credstashCR,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(credstashv1alpha1.SchemeGroupVersion, credstashCR)
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
		EXPECT().GetCredstashSecretsForCredstashSecretDefs(credstashSecretDefs).
		Return(nil, errors.New(errorString)).
		Times(1)

	_, err := r.Reconcile(req)

	assert.Error(t, err)
	assert.Contains(t, errorString, err.Error())
}

func TestCreatingSecretWhenCredstashReturnsDataAndNoSecretExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredstashSecretGetter := mocks.NewMockSecretGetter(ctrl)
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	var (
		name            = "credstashCR"
		namespace       = "credstash"
		credstashKey	= "key1"
		credstashValue	= "value1"
		credstashGetterReturn = map[string][]byte{
			credstashKey: []byte(credstashValue),
		}
	)

	credstashSecretDefs := []credstashv1alpha1.CredstashSecretDef{
		{
			Key: "key1",
		},
	}

	// A credstashCR resource with metadata and spec.
	credstashCR := &credstashv1alpha1.CredstashSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: credstashv1alpha1.CredstashSecretSpec{
			Secrets: credstashSecretDefs,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		credstashCR,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(credstashv1alpha1.SchemeGroupVersion, credstashCR)
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
		EXPECT().GetCredstashSecretsForCredstashSecretDefs(credstashSecretDefs).
		Return(credstashGetterReturn, nil).
		Times(1)

	res, err := r.Reconcile(req)

	assert.NoError(t, err)
	assert.False(t,res.Requeue)

	// Check if Secret has been created and has the correct data
	secret := &corev1.Secret{}
	err = cl.Get(context.TODO(), req.NamespacedName, secret)

	assert.NoError(t, err)
	assert.Equal(t, credstashValue, string(secret.Data[credstashKey]))
}

func TestCreatingSecretWhenCredstashReturnsDataAndSecretExistsWithDifferentData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredstashSecretGetter := mocks.NewMockSecretGetter(ctrl)
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	var (
		name            = "credstashCR"
		namespace       = "credstash"
		credstashExistingKey	= "existing key"
		credstashExistingValue	= "existing value"
		credstashKey	= "key1"
		credstashValue	= "value1"
		credstashSecretExistingData = map[string][]byte{
			credstashExistingKey: []byte(credstashExistingValue),
		}
		credstashGetterReturn = map[string][]byte{
			credstashKey: []byte(credstashValue),
		}
	)

	credstashSecretDefs := []credstashv1alpha1.CredstashSecretDef{
		{
			Key: "key1",
		},
	}

	// A credstashCR resource with metadata and spec.
	credstashCR := &credstashv1alpha1.CredstashSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: credstashv1alpha1.CredstashSecretSpec{
			Secrets: credstashSecretDefs,
		},
	}

	// A existing secret corresponding to the credstash secret
	existingSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: credstashSecretExistingData,
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		credstashCR,
		existingSecret,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(credstashv1alpha1.SchemeGroupVersion, credstashCR)
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
		EXPECT().GetCredstashSecretsForCredstashSecretDefs(credstashSecretDefs).
		Return(credstashGetterReturn, nil).
		Times(1)

	res, err := r.Reconcile(req)

	assert.NoError(t, err)
	assert.False(t,res.Requeue)

	// Check if Secret has been created and has the correct data
	secret := &corev1.Secret{}
	err = cl.Get(context.TODO(), req.NamespacedName, secret)

	assert.NoError(t, err)
	assert.Equal(t, credstashValue, string(secret.Data[credstashKey]))
}

func TestCreatingSecretWhenCredstashReturnsDataAndSecretExistsWithSameData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredstashSecretGetter := mocks.NewMockSecretGetter(ctrl)
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	var (
		name            = "credstashCR"
		namespace       = "credstash"
		credstashKey	= "key1"
		credstashValue	= "value1"
		credstashGetterReturn = map[string][]byte{
			credstashKey: []byte(credstashValue),
		}
	)

	credstashSecretDefs := []credstashv1alpha1.CredstashSecretDef{
		{
			Key: "key1",
		},
	}

	// A credstashCR resource with metadata and spec.
	credstashCR := &credstashv1alpha1.CredstashSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: credstashv1alpha1.CredstashSecretSpec{
			Secrets: credstashSecretDefs,
		},
	}

	// A existing secret corresponding to the credstash secret
	existingSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: credstashGetterReturn,
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		credstashCR,
		existingSecret,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(credstashv1alpha1.SchemeGroupVersion, credstashCR)
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
		EXPECT().GetCredstashSecretsForCredstashSecretDefs(credstashSecretDefs).
		Return(credstashGetterReturn, nil).
		Times(1)

	res, err := r.Reconcile(req)

	assert.NoError(t, err)
	assert.False(t,res.Requeue)

	// Check if Secret has been created and has the correct data
	secret := &corev1.Secret{}
	err = cl.Get(context.TODO(), req.NamespacedName, secret)

	assert.NoError(t, err)
	assert.Equal(t, credstashValue, string(secret.Data[credstashKey]))
}