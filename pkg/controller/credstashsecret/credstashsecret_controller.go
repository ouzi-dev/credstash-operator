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
	"github.com/ouzi-dev/credstash-operator/pkg/aws"
	"github.com/ouzi-dev/credstash-operator/pkg/credstash"
	"github.com/ouzi-dev/credstash-operator/pkg/flags"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	credstashv1alpha1 "github.com/ouzi-dev/credstash-operator/pkg/apis/credstash/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const LabelNameForSelector = "controllerInstance"

var log = logf.Log.WithName("controller_credstashsecret")


// Add creates a new CredstashSecret Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCredstashSecret{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("credstashsecret-controller", mgr, controller.Options{Reconciler: r})
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

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner CredstashSecret
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &credstashv1alpha1.CredstashSecret{},
	})
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
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CredstashSecret object and makes changes based on the state read
// and what is in the CredstashSecret.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
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
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Secret object
	// TODO handle error
	secret, _ := r.newSecretForCR(instance)

	// Set CredstashSecret instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Secret already exists
	found := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.client.Create(context.TODO(), secret)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Secret already exists", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name)
	return reconcile.Result{}, nil
}

// newSecretForCR returns a busybox pod with the same name/namespace as the cr
func (r *ReconcileCredstashSecret) newSecretForCR(cr *credstashv1alpha1.CredstashSecret) (*corev1.Secret, error) {
	//TODO extract all this secret logic to its own method
	//TODO allow support for default controller level credentials as a catch-all
	awsAccessKeyIdSecret := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Spec.AWSConfig.Credentials.AWSAccessKeyId.Name, Namespace: cr.Namespace}, awsAccessKeyIdSecret)
	if err != nil {
		return nil, err
	}

	awsAccessKey := string(awsAccessKeyIdSecret.Data[cr.Spec.AWSConfig.Credentials.AWSAccessKeyId.Key])

	awsSecretAccessKeySecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Spec.AWSConfig.Credentials.AWSSecretAccessKey.Name, Namespace: cr.Namespace}, awsSecretAccessKeySecret)
	if err != nil {
		return nil, err
	}

	awsSecretAccessKey := string(awsAccessKeyIdSecret.Data[cr.Spec.AWSConfig.Credentials.AWSSecretAccessKey.Key])

	awsSession, err := aws.GetAwsSession(cr.Spec.AWSConfig.Region, awsAccessKey, awsSecretAccessKey)
	credstashSecretGetter := credstash.NewHelper(awsSession)

	credstashSecretsValueMap := credstashSecretGetter.GetCredstashSecretsForCredstashSecretDefs(cr.Spec.Secrets)

	if err != nil {
		return nil, err
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		StringData: credstashSecretsValueMap,
	}, nil
}

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