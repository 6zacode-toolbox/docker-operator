package controllers

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	v1 "github.com/6zacode-toolbox/docker-operator/operator/api/v1"
	v1batch "k8s.io/api/batch/v1"
	apiV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// var NamespaceJobs = "operator-system"
var NamespaceJobs = "default"
var ServiceAccountJob = "docker-agent-sa"
var ImageJob = "6zar/docker-agent:latest"

const (
	EventCreate string = "Create"
	EventUpdate string = "Update"
	EventDelete string = "Delete"
)

func ValidateDockerHost(crd *v1.DockerHost) (*v1.DockerHost, error) {
	if crd.Spec.HostIP == "" {
		err := fmt.Errorf("error on validation: missing HostIP")
		log.Log.Error(err, fmt.Sprintf("Missing HostIP infomration on CRD:%s/%s", crd.Namespace, crd.Name))
		return crd, err
	}
	return crd, nil
}

func CreateDockerHostCronJob(crd *v1.DockerHost) (*v1batch.CronJob, error) {
	podSpec, _ := CreateDockerHostPodSpec(crd)
	cronjobMinimal := InstantiateMinimalDockerHostCronJob(crd.Name, NamespaceJobs)

	cronjobMinimal.Labels = GetLabels(crd.GetCrdDefinition())
	cronjob := &v1batch.CronJob{
		ObjectMeta: cronjobMinimal.ObjectMeta,
		Spec: v1batch.CronJobSpec{
			Schedule:          "*/3 * * * *",
			ConcurrencyPolicy: v1batch.ForbidConcurrent,
			JobTemplate: v1batch.JobTemplateSpec{
				Spec: v1batch.JobSpec{
					Template: apiV1.PodTemplateSpec{
						Spec: podSpec,
					},
				},
			},
		},
	}
	return cronjob, nil
}

func CreateDockerComposeRunnerJob(crd *v1.CrdDefinition, action string) (*v1batch.Job, error) {
	podSpec, _ := CreateDockerComposeRunnerPodSpec(crd.Name, action)
	jobMinimal := InstantiateMinimalDockerComposeRunnerJob(crd.Name, NamespaceJobs)
	jobMinimal.Labels = GetLabels(crd)
	TTL := int32(30)
	job := &v1batch.Job{
		ObjectMeta: jobMinimal.ObjectMeta,
		Spec: v1batch.JobSpec{
			TTLSecondsAfterFinished: &TTL,
			Template: apiV1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: jobMinimal.Labels,
				},
				Spec: podSpec,
			},
		},
	}
	return job, nil
}

func CreateDockerHostPodSpec(crd *v1.DockerHost) (apiV1.PodSpec, error) {
	crdConfig := crd.GetCrdDefinition()
	env := []apiV1.EnvVar{
		{
			Name:  "CRD_API_VERSION",
			Value: crdConfig.APIVersion,
		},
		{
			Name:  "CRD_NAMESPACE",
			Value: crdConfig.Namespace,
		},
		{
			Name:  "CRD_NAME",
			Value: crdConfig.Name,
		},
		{
			Name:  "CRD_RESOURCE",
			Value: crdConfig.Resource,
		},
	}
	command := "agent"
	envEnhanced := AddHostConnectionVars(crd, &env)
	result := CreateDockerAgentPod(*envEnhanced, command)
	return result, nil
}
func AddHostConnectionVars(crd *v1.DockerHost, varList *[]apiV1.EnvVar) *[]apiV1.EnvVar {
	var hostVar []apiV1.EnvVar
	//for cert setups
	if crd.Spec.SSHConnection == (v1.SSHConnection{}) {
		hostVar = []apiV1.EnvVar{
			{
				Name:  "DOCKER_CERT_PATH",
				Value: "/certs",
			},
			{
				Name:  "DOCKER_HOST",
				Value: "tcp://" + crd.Spec.HostIP + ":2376",
			},
			{
				Name:  "DOCKER_TLS_VERIFY",
				Value: "1",
			},
		}
	} else {
		hostVar = []apiV1.EnvVar{
			{
				Name:  "DOCKER_HOST",
				Value: "ssh://" + crd.Spec.SSHConnection.SSHUser + "@" + crd.Spec.HostIP,
			},
		}
	}

	for _, v := range hostVar {
		*varList = append(*varList, v)
	}
	return varList
}

func CreateDockerComposeRunnerPodSpec(name, action string) (apiV1.PodSpec, error) {
	configMapNMame := GenerateComposeRunnerConfigMapName(name)
	env := []apiV1.EnvVar{
		{
			Name: "DOCKER_CERT_PATH",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "DOCKER_CERT_PATH",
				},
			},
		},
		{
			Name: "DOCKER_HOST",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "DOCKER_HOST",
				},
			},
		},
		{
			Name: "DOCKER_TLS_VERIFY",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "DOCKER_TLS_VERIFY",
				},
			},
		},
		{
			Name: "COMPOSE_FILE",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "COMPOSE_FILE",
				},
			},
		},
		{
			Name: "REPO_ADDRESS",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "REPO_ADDRESS",
				},
			},
		},
		{
			Name: "EXECUTION_PATH",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "EXECUTION_PATH",
				},
			},
		},
		{
			Name: "CRD_API_VERSION",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "CRD_API_VERSION",
				},
			},
		},
		{
			Name: "CRD_NAMESPACE",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "CRD_NAMESPACE",
				},
			},
		},
		{
			Name: "CRD_NAME",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "CRD_NAME",
				},
			},
		},
		{
			Name: "CRD_RESOURCE",
			ValueFrom: &apiV1.EnvVarSource{
				ConfigMapKeyRef: &apiV1.ConfigMapKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: configMapNMame,
					},
					Key: "CRD_RESOURCE",
				},
			},
		},
		//To get from a Secret froom vault
		{
			Name: "GITHUB_TOKEN",
			ValueFrom: &apiV1.EnvVarSource{
				SecretKeyRef: &apiV1.SecretKeySelector{
					LocalObjectReference: apiV1.LocalObjectReference{
						Name: "github",
					},
					Key: "GITHUB_TOKEN",
				},
			},
		},
		{
			Name:  "ACTION",
			Value: action,
		},
	}
	command := "compose-runner"

	result := CreateDockerAgentPod(env, command)
	return result, nil
}

func CreateDockerAgentPod(env []apiV1.EnvVar, command string) apiV1.PodSpec {
	container := apiV1.Container{
		Name:            "docker-agent",
		Image:           ImageJob,
		ImagePullPolicy: apiV1.PullAlways,
		Env:             env,
		Command:         []string{"/home/app/docker-agent"},
		//Command: []string{"tail", "-f", "/dev/null"},
		Args: []string{command, "--crd-api-version", "$(CRD_API_VERSION)", "--crd-namespace", "$(CRD_NAMESPACE)", "--crd-instance", "$(CRD_NAME)", "--crd-resource", "$(CRD_RESOURCE)"},
		VolumeMounts: []apiV1.VolumeMount{
			{
				MountPath: "/certs",
				Name:      "docker-certs",
				ReadOnly:  true,
			},
			{
				MountPath: "/home/app/.ssh",
				Name:      "docker-ssh-volume",
				ReadOnly:  true,
			},
		},
	}
	result := apiV1.PodSpec{
		ServiceAccountName: ServiceAccountJob,
		RestartPolicy:      apiV1.RestartPolicyNever,
		Containers:         []apiV1.Container{container},
		Volumes: []apiV1.Volume{
			{
				Name: "docker-certs",
				VolumeSource: apiV1.VolumeSource{
					Secret: &apiV1.SecretVolumeSource{
						SecretName: "docker-secret",
					},
				},
			},
			{
				Name: "docker-ssh-volume",
				VolumeSource: apiV1.VolumeSource{
					Secret: &apiV1.SecretVolumeSource{
						SecretName: "docker-ssh",
					},
				},
			},
		},
	}
	return result
}

func InstantiateMinimalDockerHostCronJob(name string, namespace string) *v1batch.CronJob {
	jobName := name + "-dh-cron"
	job := &v1batch.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: v1batch.CronJobSpec{},
	}
	return job
}

func InstantiateMinimalDockerComposeRunnerJob(name string, namespace string) *v1batch.Job {
	jobName := name + "-dcr-job-" + RandStringRunes(4)
	job := &v1batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: v1batch.JobSpec{},
	}
	return job
}

func checkCronJob(job *v1batch.CronJob) {
	if job == nil {
		log.Log.Info("Job is nil")
	} else if reflect.DeepEqual(job, v1batch.CronJob{}) {
		log.Log.Info("Job is blank:" + job.Name)
	} else {
		log.Log.Info("Job exist:" + job.Name)
	}
}

func CreateTLSDockerComposeRunnerConfigMap(crd *v1.DockerComposeRunner) *apiV1.ConfigMap {
	configMapName := GenerateComposeRunnerConfigMapName(crd.Name)
	crdConfig := crd.GetCrdDefinition()
	labels := GetLabels(crdConfig)
	configMap := &apiV1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: NamespaceJobs,
			Labels:    labels,
		},
		Data: map[string]string{
			"DOCKER_CERT_PATH":  "/certs",
			"DOCKER_HOST":       "tcp://" + crd.Spec.HostIP + ":2376",
			"DOCKER_TLS_VERIFY": "1",
			"HOST_IP":           crd.Spec.HostIP,
			"COMPOSE_FILE":      crd.Spec.ComposeFile,
			"REPO_ADDRESS":      crd.Spec.RepoAddress,
			"EXECUTION_PATH":    crd.Spec.ExecutionPath,
			"CRD_API_VERSION":   crdConfig.APIVersion,
			"CRD_NAMESPACE":     crdConfig.Namespace,
			"CRD_NAME":          crdConfig.Name,
			"CRD_RESOURCE":      crdConfig.Resource,
		},
	}
	return configMap
}

func CreateSSHDockerComposeRunnerConfigMap(crd *v1.DockerComposeRunner, sshUser string, hostIp string) *apiV1.ConfigMap {
	configMapName := GenerateComposeRunnerConfigMapName(crd.Name)
	crdConfig := crd.GetCrdDefinition()
	labels := GetLabels(crdConfig)
	configMap := &apiV1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: NamespaceJobs,
			Labels:    labels,
		},
		Data: map[string]string{
			"DOCKER_CERT_PATH":  "",
			"DOCKER_HOST":       "ssh://" + sshUser + "@" + hostIp,
			"DOCKER_TLS_VERIFY": "",
			"HOST_IP":           crd.Spec.HostIP,
			"COMPOSE_FILE":      crd.Spec.ComposeFile,
			"REPO_ADDRESS":      crd.Spec.RepoAddress,
			"EXECUTION_PATH":    crd.Spec.ExecutionPath,
			"CRD_API_VERSION":   crdConfig.APIVersion,
			"CRD_NAMESPACE":     crdConfig.Namespace,
			"CRD_NAME":          crdConfig.Name,
			"CRD_RESOURCE":      crdConfig.Resource,
		},
	}
	return configMap
}

func GenerateComposeRunnerConfigMapName(crdName string) string {
	return crdName + "-dcr-cm"
}

func GetLabels(crdConfig *v1.CrdDefinition) map[string]string {
	result := map[string]string{
		"instance":      crdConfig.Name,
		"resouce-owner": crdConfig.Resource,
	}
	return result
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
