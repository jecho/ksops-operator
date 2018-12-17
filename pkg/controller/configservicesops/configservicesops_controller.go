/*
Copyright 2018 The Kubernetes Authors.

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

package configservicesops

import (
	"context"
	"gopkg.in/yaml.v2"
	"log"

	mygroupv1beta1 "github.com/jecho/ksops-test/pkg/apis/mygroup/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"go.mozilla.org/sops/decrypt"
	"k8s.io/client-go/kubernetes/scheme"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ConfigServiceSops Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this mygroup.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileConfigServiceSops{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("configservicesops-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to ConfigServiceSops
	err = c.Watch(&source.Kind{Type: &mygroupv1beta1.ConfigServiceSops{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mygroupv1beta1.ConfigServiceSops{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileConfigServiceSops{}

// ReconcileConfigServiceSops reconciles a ConfigServiceSops object
type ReconcileConfigServiceSops struct {
	client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileConfigServiceSops) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the ConfigServiceSops instance
	instance := &mygroupv1beta1.ConfigServiceSops{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	decodedManifest, err := decrypt.Data([]byte(instance.Spec.Manifest), "yaml")
	if err != nil {
		log.Println("Unable to decrypt payload.")
		log.Println(err)
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(decodedManifest), nil, nil)

	service := &corev1.Service{}
	if err != nil {
		log.Println(err)
	} else {
		// load the appropriate apiVersion
		service = obj.(*corev1.Service)
	}

	// check if namespace is nil and set it appropriately
	if service.Namespace == "" || len(service.Namespace) == 0 {
		service.Namespace = instance.Namespace
	}

	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// boilerplate
	// needs revisions
	found := &corev1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Printf("Creating Service %s/%s\n", service.Namespace, service.Name)
		err = r.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func createYaml(manifest string) map[interface{}]interface{} {
	var config interface{}
	err := yaml.Unmarshal([]byte(manifest), &config)
	if err != nil {
		log.Println(err)
	}

	data := config.(map[interface{}]interface{})

	return data
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
