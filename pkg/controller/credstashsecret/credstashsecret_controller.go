/*
Copyright 2020 Ouzi Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package credstashsecret

import (
	"context"
	"reflect"

	"github.com/ouzi-dev/credstash-operator/pkg/aws"
	"github.com/ouzi-dev/credstash-operator/pkg/credstash"
	event_consts "github.com/ouzi-dev/credstash-operator/pkg/event"
	"github.com/ouzi-dev/credstash-operator/pkg/flags"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	credstashv1alpha1 "github.com/ouzi-dev/credstash-operator/pkg/apis/credstash/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	LabelNameForSelector = "operatorInstance"
	ControllerName       = "credstashsecret-controller"
)

var log = logf.Log.WithName("controller_credstashsecret")

// Add creates a new CredstashSecret Controller and adds it to the Manager.
// The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	reconciler, err := newReconciler(mgr)
	if err != nil {
		return err
	}

	return add(mgr, reconciler)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) (reconcile.Reconciler, error) {
	awsSession, err := aws.GetAwsSessionFromEnv()
	if err != nil {
		return nil, err
	}

	eventRecorder := mgr.GetEventRecorderFor(ControllerName)

	return &ReconcileCredstashSecret{
		client:                mgr.GetClient(),
		scheme:                mgr.GetScheme(),
		credstashSecretGetter: credstash.NewSecretGetter(awsSession, eventRecorder),
		eventRecorder:         eventRecorder,
	}, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// if using a label selector, ignore the CRs that we should not process
	pred := setupPredicateFuncs()

	// Watch for changes to primary resource CredstashSecret
	err = c.Watch(&source.Kind{Type: &credstashv1alpha1.CredstashSecret{}}, &handler.EnqueueRequestForObject{}, pred)
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner CredstashSecret
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &credstashv1alpha1.CredstashSecret{},
	}, pred)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCredstashSecret implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCredstashSecret{}

// ReconcileCredstashSecret reconciles a CredstashSecret object
type ReconcileCredstashSecret struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client                client.Client
	scheme                *runtime.Scheme
	credstashSecretGetter credstash.SecretGetter
	eventRecorder         record.EventRecorder
}

// Reconcile reads that state of the cluster for a CredstashSecret object and makes changes based on the state read
// and what is in the CredstashSecret.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
//nolint funlen
func (r *ReconcileCredstashSecret) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Namespace", request.Namespace, "Name", request.Name)
	reqLogger.Info("Reconciling CredstashSecret")

	// Fetch the CredstashSecret instance
	instance := &credstashv1alpha1.CredstashSecret{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		r.eventRecorder.Eventf(instance, event_consts.TypeWarning, event_consts.ReasonErrGeneric, err.Error())
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Secret object
	secret, err := r.secretForCR(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set CredstashSecret instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
		r.eventRecorder.Eventf(instance, event_consts.TypeWarning, event_consts.ReasonErrGeneric, err.Error())
		return reconcile.Result{}, err
	}

	// Check if this Secret already exists
	found := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, found)

	// Create new secret if it doesn't exist
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.client.Create(context.TODO(), secret)
		if err != nil {
			r.eventRecorder.Eventf(instance, event_consts.TypeWarning, event_consts.ReasonErrCreateSecret, event_consts.MessageFailedCreatingSecret,
				secret.Name, secret.Namespace, err.Error())
			return reconcile.Result{}, err
		}

		// Secret name has changed
		if instance.Status.SecretName != "" && secret.Name != instance.Status.SecretName {
			secretToDelete := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      instance.Status.SecretName,
					Namespace: instance.Namespace,
				},
			}

			reqLogger.Info(
				"Deleting old secret since name has changed",
				"Secret.Namespace",
				secretToDelete.Namespace,
				"Secret.Name",
				secretToDelete.Name)

			err = r.client.Delete(context.TODO(), secretToDelete)
			if err != nil {
				r.eventRecorder.Eventf(instance, event_consts.TypeWarning, event_consts.ReasonErrDeleteOldSecret, event_consts.MessageFailedDeletingOldSecret,
					secret.Name, secret.Namespace, err.Error())
				return reconcile.Result{}, err
			}
			r.eventRecorder.Eventf(instance, event_consts.TypeNormal, event_consts.ReasonSuccessDeleteOldSecret, event_consts.MessageSuccessDeletingOldSecret,
				secret.Name, secret.Namespace)

		}

		instance.Status.SecretName = secret.Name

		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			r.eventRecorder.Eventf(instance, event_consts.TypeWarning, event_consts.ReasonErrGeneric, err.Error())
			return reconcile.Result{}, err
		}

		// Secret created successfully - don't requeue
		r.eventRecorder.Eventf(instance, event_consts.TypeNormal, event_consts.ReasonSuccessCreateSecret, event_consts.MessageSuccessCreatingSecret,
			secret.Name, secret.Namespace)
		return reconcile.Result{}, nil
	} else if err != nil {
		r.eventRecorder.Eventf(instance, event_consts.TypeWarning, event_consts.ReasonErrGeneric, err.Error())
		return reconcile.Result{}, err
	}

	// Secret data or type is out of date with kalamaja secret
	if !reflect.DeepEqual(secret.Data, found.Data) || !reflect.DeepEqual(secret.Type, found.Type) {
		reqLogger.Info(
			"Updating Secret because contents have changed",
			"Secret.Namespace",
			secret.Namespace,
			"Secret.Name",
			secret.Name)

		err = r.client.Update(context.TODO(), secret)
		if err != nil {
			r.eventRecorder.Eventf(instance, event_consts.TypeWarning, event_consts.ReasonErrUpdateSecret, event_consts.MessageFailedUpdatingSecret,
				secret.Name, secret.Namespace, err.Error())
			return reconcile.Result{}, err
		}
		r.eventRecorder.Eventf(instance, event_consts.TypeNormal, event_consts.ReasonSuccessUpdateSecret, event_consts.MessageSuccessUpdatingSecret,
			secret.Name, secret.Namespace)
		return reconcile.Result{}, nil
	}

	// Secret already exists - don't requeue
	reqLogger.Info(
		"Skip reconcile: Secret already exists and is up to date",
		"Secret.Namespace",
		found.Namespace,
		"Secret.Name",
		found.Name)
	return reconcile.Result{}, nil
}

// secretForCR returns a secret the same name/namespace as the cr
func (r *ReconcileCredstashSecret) secretForCR(cr *credstashv1alpha1.CredstashSecret) (*corev1.Secret, error) {
	credstashSecretsValueMap, err := r.credstashSecretGetter.GetCredstashSecretsForCredstashSecret(cr)
	if err != nil {
		return nil, err
	}

	// default to custom resource name if name is not provided
	secretName := cr.Spec.SecretName
	if secretName == "" {
		secretName = cr.Name
	}

	// default to Opaque if not provided
	secretType := cr.Spec.SecretType
	if secretType == "" {
		secretType = corev1.SecretTypeOpaque
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: cr.Namespace,
			Labels:    cr.GetLabels(),
		},
		Data: credstashSecretsValueMap,
		Type: secretType,
	}

	return secret, nil
}

// setupPredicateFuncs makes sure that we only watch resources that match the correct label selector that we want
// nolint funlen
func setupPredicateFuncs() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			predLogger := log.WithValues("Namespace", e.Meta.GetNamespace(), "Name", e.Meta.GetName(),
				"LabelSelector", flags.SelectorLabelValue)

			if flags.SelectorLabelValue == "" {
				return true
			}

			val, ok := e.Meta.GetLabels()[LabelNameForSelector]

			shouldProcess := ok && val == flags.SelectorLabelValue
			if shouldProcess {
				predLogger.V(1).Info("Processing CR since it matches label selector")
			} else {
				predLogger.V(1).Info("Not processing CR since it doesn't match label selector")
			}

			return shouldProcess
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			predLogger := log.WithValues("Namespace", e.MetaNew.GetNamespace(), "Name", e.MetaNew.GetName(),
				"LabelSelector", flags.SelectorLabelValue)

			if flags.SelectorLabelValue == "" {
				return true
			}

			val, ok := e.MetaNew.GetLabels()[LabelNameForSelector]

			shouldProcess := ok && val == flags.SelectorLabelValue
			if shouldProcess {
				predLogger.V(1).Info("Processing CR since it matches label selector")
			} else {
				predLogger.V(1).Info("Not processing CR since it doesn't match label selector")
			}

			return shouldProcess
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			predLogger := log.WithValues("Namespace", e.Meta.GetNamespace(), "Name", e.Meta.GetName(),
				"LabelSelector", flags.SelectorLabelValue)

			if flags.SelectorLabelValue == "" {
				return true
			}

			val, ok := e.Meta.GetLabels()[LabelNameForSelector]

			shouldProcess := ok && val == flags.SelectorLabelValue
			if shouldProcess {
				predLogger.V(1).Info("Processing CR since it matches label selector")
			} else {
				predLogger.V(1).Info("Not processing CR since it doesn't match label selector")
			}

			return shouldProcess
		},
		GenericFunc: func(e event.GenericEvent) bool {
			predLogger := log.WithValues("Namespace", e.Meta.GetNamespace(), "Name", e.Meta.GetName(),
				"LabelSelector", flags.SelectorLabelValue)

			if flags.SelectorLabelValue == "" {
				return true
			}

			val, ok := e.Meta.GetLabels()[LabelNameForSelector]

			shouldProcess := ok && val == flags.SelectorLabelValue
			if shouldProcess {
				predLogger.V(1).Info("Processing CR since it matches label selector")
			} else {
				predLogger.V(1).Info("Not processing CR since it doesn't match label selector")
			}

			return shouldProcess
		},
	}
}
