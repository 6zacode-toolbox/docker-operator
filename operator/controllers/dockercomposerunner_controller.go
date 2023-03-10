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

	"google.golang.org/grpc/status"
	v1batch "k8s.io/api/batch/v1"
	apiV1 "k8s.io/api/core/v1"
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

	toolv1 "github.com/6zacode-toolbox/docker-operator/operator/api/v1"
	v1 "github.com/6zacode-toolbox/docker-operator/operator/api/v1"
)

// DockerComposeRunnerReconciler reconciles a DockerComposeRunner object
type DockerComposeRunnerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tool.6zacode-toolbox.github.io,resources=dockercomposerunners,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups=tool.6zacode-toolbox.github.io,resources=dockercomposerunners/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tool.6zacode-toolbox.github.io,resources=dockercomposerunners/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get
//+kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups=batch,resources=cronjobs/status,verbs=get
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete;deletecollection

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DockerComposeRunner object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile

func (r *DockerComposeRunnerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	/*
		Logic:
			If CRD deleted:
				- run down Job based on config map
				- run down Job
				- delete configMap and jobs
			IF not delete is a create/update....
			- Create ConfigMap/Up Job
			If Create:
				- Create configMap
				- Run up Job
	*/

	eventType := EventUpdate
	log.Log.Info(fmt.Sprintf("Called: %#v", req.NamespacedName))
	//Desired State
	instance := &toolv1.DockerComposeRunner{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		log.Log.Info(string(err.Error()))
		status, _ := status.FromError(err)
		log.Log.Info(status.Code().String())
		if errork8s.IsNotFound(err) || strings.Contains(string(err.Error()), "not found") {
			r.RunComposeDownJob(req.Name, NamespaceJobs)
			eventType = EventDelete
			log.Log.Info(fmt.Sprintf("Event: %s =>  %#v", eventType, req.NamespacedName))
			log.Log.Info("Event Type:" + eventType)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Log.Error(err, "Event Type(missing object):"+eventType)
		return reconcile.Result{}, err
	}
	// Check if Job was already created:
	if instance.Status.Instanced || instance.Status.Validated {
		// nothing to be done
		log.Log.Info("Event Type(object instantiated):" + eventType)
		return reconcile.Result{}, nil
	}
	err = r.CleanOldJob(instance.GetCrdDefinition())
	if err != nil {
		log.Log.Error(err, "Error cleaning previous jobs")
		return reconcile.Result{}, err
	}

	currentStateJob, err := r.GetDockerComposeRunnerCurrentState(req.Name)
	if err != nil {
		log.Log.Error(err, "Event Type(err getting currentState):"+eventType)
		return reconcile.Result{}, err
	}
	//missing job, creating one
	if currentStateJob == nil {
		eventType = EventCreate
		log.Log.Info(fmt.Sprintf("Event: %s =>  %#v", eventType, req.NamespacedName))
	}
	//If is not delete(early event), if it didn't exist it must be deleted.
	if currentStateJob == nil {
		eventType = EventCreate
		log.Log.Info(fmt.Sprintf("Event: %s =>  %#v", eventType, req.NamespacedName))
	}

	//Update detected, but no job exist
	err = r.RunComposeUpJob(instance)
	if err != nil {
		log.Log.Error(err, "Event Type(err running compose down/up):"+eventType)
		return reconcile.Result{}, err
	}
	log.Log.Info("Event Type(final return):" + eventType)
	return ctrl.Result{}, nil
}

func (r *DockerComposeRunnerReconciler) RunComposeUpJob(instance *toolv1.DockerComposeRunner) error {

	usingSSHConnection := false
	dockerHostCrd := &toolv1.DockerHost{}
	if instance.Spec.ResourceOwner != "" {
		err := r.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.ResourceOwner, Namespace: instance.Namespace}, dockerHostCrd)
		if err != nil {
			log.Log.Error(err, fmt.Sprintf("Getting connection details: %s/%s", instance.Namespace, instance.Name))
			return err
		} else {
			if dockerHostCrd.Spec.SSHConnection != (v1.SSHConnection{}) {
				usingSSHConnection = true
			}
		}
		log.Log.Info(fmt.Sprintf("Using SSH Connection: %t", usingSSHConnection))
	}

	log.Log.Info("RunComposeUpJob:" + instance.Name)

	definition := instance.GetCrdDefinition()
	cms := &apiV1.ConfigMap{}
	err := r.DeleteAllOf(context.TODO(), cms, client.InNamespace(definition.Namespace), client.MatchingLabels(GetLabels(definition)), client.GracePeriodSeconds(5))
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Error deleting configmaps from previous runs: %s/%s", definition.Namespace, definition.Name))
		return err
	}
	configMap := &apiV1.ConfigMap{}
	configMapName := GenerateComposeRunnerConfigMapName(instance.Name)
	if usingSSHConnection {
		configMap = CreateSSHDockerComposeRunnerConfigMap(instance, dockerHostCrd.Spec.SSHConnection.SSHUser, dockerHostCrd.Spec.HostIP, configMapName)
	} else {
		configMap = CreateTLSDockerComposeRunnerConfigMap(instance, configMapName)
	}
	err = r.Create(context.TODO(), configMap)
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Error creating configmap %s/%s", instance.Namespace, configMap))
		return err

	}
	err = r.DefineConfigMapName(instance, configMapName, *dockerHostCrd)
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Error defining configmap on CRDs %s/%s", instance.Namespace, configMap))
		return err

	}

	err = r.CreateJob(instance.GetCrdDefinition(), toolv1.COMPOSE_ACTION_UP, configMapName)
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Error creating job dcr job %s/%s", instance.Namespace, instance.Name))
		return err
	}
	return nil
}

func (r *DockerComposeRunnerReconciler) DefineConfigMapName(instance *toolv1.DockerComposeRunner, configMapName string, dockerHost toolv1.DockerHost) error {
	patchInstance := &toolv1.DockerComposeRunner{
		ObjectMeta: instance.ObjectMeta,
		Status: v1.DockerComposeRunnerStatus{
			ConfigMapName: configMapName,
			Validated:     true,
			ResourceOwner: instance.Spec.ResourceOwner,
		},
	}

	err := r.Status().Update(context.TODO(), patchInstance)
	if err != nil {
		log.Log.Error(err, "trying to set configmap name:"+instance.Name)
		return err
	}
	//update dockerhost
	if !reflect.DeepEqual(dockerHost, toolv1.DockerHost{}) {

		if dockerHost.Status.AppConfigMap == nil {
			dockerHost.Status.AppConfigMap = make(map[string]string)
		}
		_, ok := dockerHost.Status.AppConfigMap[instance.Name]
		// If the key exists
		if ok {
			// Do something
			log.Log.Info("Status update - key already present on dockerhost")

		}
		dockerHost.Status.AppConfigMap[instance.Name] = configMapName

		err := r.Status().Update(context.TODO(), &dockerHost)
		if err != nil {
			log.Log.Error(err, "trying to set configmap name:"+instance.Name)
			return err
		}
	}

	log.Log.Info("Status update - config map name")
	return err
}
func (r *DockerComposeRunnerReconciler) RunComposeDownJob(name string, namespace string) error {
	definition := &toolv1.CrdDefinition{
		Name:       name,
		Namespace:  namespace,
		APIVersion: "tool.6zacode-toolbox.github.io/v1",
		Resource:   "dockercomposerunners",
	}

	configMapList := &apiV1.ConfigMapList{}
	err := r.List(context.TODO(), configMapList, client.InNamespace(definition.Namespace), client.MatchingLabels(GetLabels(definition)))
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Error creating job object %s/%s", definition.Namespace, definition.Name))
		return err

	}
	if len(configMapList.Items) < 1 {
		log.Log.Info(fmt.Sprintf("Nothing found to be deleted, missing configmap for: %s/%s", definition.Namespace, definition.Name))
	}
	for _, cm := range configMapList.Items {
		configMapName := cm.Name
		err = r.CreateJob(definition, toolv1.COMPOSE_ACTION_DOWN, configMapName)
		if err != nil {
			log.Log.Error(err, fmt.Sprintf("Error creating job object %s/%s", definition.Namespace, definition.Name))
			return err
		}
		break
	}

	return nil
}

func (r *DockerComposeRunnerReconciler) CleanOldJob(definition *toolv1.CrdDefinition) error {
	//remove old jobs
	jobs := &v1batch.Job{}
	err := r.DeleteAllOf(context.TODO(), jobs, client.InNamespace(definition.Namespace), client.MatchingLabels(GetLabels(definition)), client.GracePeriodSeconds(5))
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Error creating job object %s/%s", definition.Namespace, definition.Name))
		return err

	}

	//remove old pods
	pods := &apiV1.Pod{}
	err = r.DeleteAllOf(context.TODO(), pods, client.InNamespace(definition.Namespace), client.MatchingLabels(GetLabels(definition)), client.GracePeriodSeconds(5))
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Error creating job object %s/%s", definition.Namespace, definition.Name))
		return err
	}
	return nil
}
func (r *DockerComposeRunnerReconciler) CreateJob(definition *toolv1.CrdDefinition, action string, configMapName string) error {

	//create new job
	job, err := CreateDockerComposeRunnerJob(definition, action, configMapName)
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Error creating DCR JOB: %s/%s", definition.Namespace, job.Name))
		return err

	}
	err = r.Create(context.TODO(), job)
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Error creating job: %s/%s", definition.Namespace, job.Name))
		return err

	}
	return nil
}

func (r *DockerComposeRunnerReconciler) CreateDockerComposeRunnerJob(desiredJob *v1batch.Job, instance *toolv1.DockerComposeRunner) error {
	err := r.Create(context.TODO(), desiredJob)
	if err != nil {
		return nil
	}
	patchInstance := &toolv1.DockerComposeRunner{
		ObjectMeta: instance.ObjectMeta,
		Status: v1.DockerComposeRunnerStatus{
			Instanced:     true,
			ResourceOwner: instance.Spec.ResourceOwner,
		},
	}
	//err = r.Status().Update(context.TODO(), patchInstance)
	//err = r.Status().Update(context.Background(), instance)
	log.Log.Info("Status update - instatiated")
	err = r.Status().Patch(context.TODO(), patchInstance, client.Merge)
	if err != nil {
		log.Log.Error(err, "Event Type(err Patch status).")
	}
	return err
}

func (r *DockerComposeRunnerReconciler) DeleteDockerComposeRunnerJob(name string, namespace string) error {
	// How to handle delete events on which you dont have enough info to fill an unkown deployment
	// Delete before a create
	// Last known state is needed to allow to know where to kill the compose
	job := InstantiateMinimalDockerComposeRunnerJob(name, namespace)
	err := r.Delete(context.TODO(), job)
	if err != nil {
		log.Log.Info(fmt.Sprintf("Error deleting Found job %s/%s", namespace, job))
		return err

	}
	return nil
}

func (r *DockerComposeRunnerReconciler) GetDockerComposeRunnerCurrentState(crdName string) (*v1batch.Job, error) {
	job := InstantiateMinimalDockerComposeRunnerJob(crdName, NamespaceJobs)
	jobFound := &v1batch.Job{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: NamespaceJobs}, jobFound)
	if err != nil && errors.IsNotFound(err) {
		log.Log.Info(fmt.Sprintf("Not Found Job %s/%s", NamespaceJobs, job.Name))
		jobFound = nil
	}
	return jobFound, nil
}

// SetupWithManager sets up the controller with the Manager and its filters.
func (r *DockerComposeRunnerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolv1.DockerComposeRunner{}).
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
