/*
Copyright 2023 6zacode-toolbox.

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

package controllers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	v1 "github.com/6zacode-toolbox/docker-operator/operator/api/v1"
	"google.golang.org/grpc/status"
	v1batch "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cri-api/pkg/errors"
	errork8s "k8s.io/cri-api/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// DockerHostReconciler reconciles a DockerHost object
type DockerHostReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tool.6zacode-toolbox.github.io,resources=dockerhosts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tool.6zacode-toolbox.github.io,resources=dockerhosts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tool.6zacode-toolbox.github.io,resources=dockerhosts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DockerHost object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DockerHostReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	// Logic - Create a Cronjob to check status of host every X minutes
	// Only list for events of create/delete
	// Create: Creates Cronjob per host/CRD
	// Delete: Delete Cronjob per host/CRD
	// How to handle the sercret needed for the host
	// CronJobNamespace
	// Sercret need to be at the same namespace

	eventType := EventUpdate
	// TODO(user): your logic here
	//log.Log.Info(fmt.Sprintf("Called: %#v", req))
	log.Log.Info(fmt.Sprintf("Called: %#v", req.NamespacedName))

	//Desired State
	instance := &v1.DockerHost{}

	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		log.Log.Info(string(err.Error()))
		status, _ := status.FromError(err)
		log.Log.Info(status.Code().String())
		if errork8s.IsNotFound(err) || strings.Contains(string(err.Error()), "not found") {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			// CRD was removed (DELETE EVENT)
			// How to remove the objects?
			log.Log.Info("Delete Job")
			r.DeleteDockerHostCronJob(req.Name, NamespaceJobs)
			eventType = EventDelete
			log.Log.Info(fmt.Sprintf("Event: %s =>  %#v", eventType, req.NamespacedName))
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	// Check if Job was already created:
	if instance.Status.Instanced {
		// nothing to be done
		return reconcile.Result{}, nil
	}

	desiredJob, err := CreateDockerHostCronJob(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	currentStateJob, err := r.GetDockerHostCurrentState(req.Name)
	if err != nil {
		return reconcile.Result{}, err
	}
	checkCronJob(currentStateJob)

	//missing job, creating one
	if currentStateJob == nil {
		eventType = EventCreate
		log.Log.Info(fmt.Sprintf("Event: %s =>  %#v", eventType, req.NamespacedName))
	}
	//If is not delete(early event), if it didn't exist it must be deleted.
	if currentStateJob == nil {
		err = r.CreateDockerHostCronJob(desiredJob, instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	//Update detected, but no job exist
	if eventType == EventUpdate && currentStateJob != nil {
		if !reflect.DeepEqual(desiredJob.Spec, currentStateJob.Spec) {
			r.DeleteDockerHostCronJob(req.Name, NamespaceJobs)
			err = r.CreateDockerHostCronJob(desiredJob, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	} else {
		log.Log.Info(fmt.Sprintf("Unkown State: %#v", eventType))
	}
	return ctrl.Result{}, nil
}

func (r *DockerHostReconciler) CreateDockerHostCronJob(desiredJob *v1batch.CronJob, instance *v1.DockerHost) error {
	err := r.Create(context.TODO(), desiredJob)
	if err != nil {
		return nil
	}
	instance.Status.Instanced = true
	err = r.Status().Update(context.Background(), instance)
	return err
}

func (r *DockerHostReconciler) DeleteDockerHostCronJob(name string, namespace string) error {
	job := InstantiateMinimalDockerHostCronJob(name, namespace)
	err := r.Delete(context.TODO(), job)
	if err != nil {
		log.Log.Info(fmt.Sprintf("Error deleting Found job %s/%s\n", namespace, job))
		return err

	}
	return nil
}

func (r *DockerHostReconciler) GetDockerHostCurrentState(crdName string) (*v1batch.CronJob, error) {
	job := InstantiateMinimalDockerHostCronJob(crdName, NamespaceJobs)
	jobFound := &v1batch.CronJob{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: NamespaceJobs}, jobFound)
	if err != nil && errors.IsNotFound(err) {
		log.Log.Info(fmt.Sprintf("Not Found Job %s/%s\n", NamespaceJobs, job.Name))
		jobFound = nil
	}
	return jobFound, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DockerHostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.DockerHost{}).
		WithEventFilter(predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				oldGeneration := e.ObjectOld.GetGeneration()
				newGeneration := e.ObjectNew.GetGeneration()
				// Generation is only updated on spec changes (also on deletion),
				// not metadata or status
				// Filter out events where the generation hasn't changed to
				// avoid being triggered by status updates
				return oldGeneration != newGeneration
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				// The reconciler adds a finalizer so we perform clean-up
				// when the delete timestamp is added
				// Suppress Delete events to avoid filtering them out in the Reconcile function
				return true
			},
		}).
		Complete(r)
}
