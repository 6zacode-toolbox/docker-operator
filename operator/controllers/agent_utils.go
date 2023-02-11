package controllers

import (
	"reflect"

	v1 "github.com/6zacode-toolbox/docker-operator/operator/api/v1"
	v1batch "k8s.io/api/batch/v1"
	apiV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var NamespaceJobs = "default"
var ServiceAccountJob = "docker-agent-sa"
var ImageJob = "6zar/docker-agent:latest"

const (
	EventCreate string = "Create"
	EventUpdate string = "Update"
	EventDelete string = "Delete"
)

func CreateDockerHostCronJob(crd *v1.DockerHost) (*v1batch.CronJob, error) {
	podSpec, _ := CreateDockerHostPodSpec(crd)
	cronjobMinimal := InstantiateMinimalDockerHostCronJob(crd.Name, NamespaceJobs)

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

func CreateDockerComposeRunnerJob(crd *v1.DockerComposeRunner, action string) (*v1batch.Job, error) {
	podSpec, _ := CreateDockerComposeRunnerPodSpec(crd, action)
	cronjobMinimal := InstantiateMinimalDockerComposeRunnerJob(crd.Name, NamespaceJobs)

	cronjob := &v1batch.Job{
		ObjectMeta: cronjobMinimal.ObjectMeta,
		Spec: v1batch.JobSpec{
			Template: apiV1.PodTemplateSpec{
				Spec: podSpec,
			},
		},
	}
	return cronjob, nil
}

func CreateDockerHostPodSpec(crd *v1.DockerHost) (apiV1.PodSpec, error) {
	env := []apiV1.EnvVar{
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
	command := "agent"

	result := CreateDockerAgentPod(env, command, crd.GetCrdDefinition())
	return result, nil
}

func CreateDockerComposeRunnerPodSpec(crd *v1.DockerComposeRunner, action string) (apiV1.PodSpec, error) {
	/*
	   echo $COMPOSE_FILE
	   echo $REPO_ADDRESS
	   echo $EXECUTION_PATH
	   echo $GITHUB_TOKEN
	   # "up -d" or "down"
	   echo $ACTION
	*/

	env := []apiV1.EnvVar{
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
		{
			Name:  "COMPOSE_FILE",
			Value: crd.Spec.ComposeFile,
		},
		{
			Name:  "REPO_ADDRESS",
			Value: crd.Spec.RepoAddress,
		},
		{
			Name:  "EXECUTION_PATH",
			Value: crd.Spec.ExecutionPath,
		},
		//To get from a Secret froom vault
		{
			Name:  "GITHUB_TOKEN",
			Value: "TO get from a secret",
		},
		{
			Name:  "ACTION",
			Value: action,
		},
	}
	command := "agent"

	result := CreateDockerAgentPod(env, command, crd.GetCrdDefinition())
	return result, nil
}

func CreateDockerAgentPod(env []apiV1.EnvVar, command string, crd *v1.CrdDefinition) apiV1.PodSpec {
	container := apiV1.Container{
		Name:            "docker-agent",
		Image:           ImageJob,
		ImagePullPolicy: apiV1.PullAlways,
		Env:             env,
		Command:         []string{"/home/app/docker-agent"},
		Args:            []string{command, "--crd-api-version", crd.APIVersion, "--crd-namespace", crd.Namespace, "--crd-instance", crd.Name, "--crd-resource", crd.Resource},
		VolumeMounts: []apiV1.VolumeMount{{
			MountPath: "/certs",
			Name:      "docker-certs",
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
	jobName := name + "-dcr-job"
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

func CreateDockerComposeRunnerConfigMap(crd *v1.DockerComposeRunner) *apiV1.ConfigMap {
	configMapName := GenerateComposeRunnerConfigMapName(crd.Name)
	configMap := &apiV1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: NamespaceJobs,
		},
		Data: map[string]string{
			"hostip":        crd.Spec.HostIP,
			"composeFile":   crd.Spec.ComposeFile,
			"executionPath": crd.Spec.ExecutionPath,
			"repoAddress":   crd.Spec.RepoAddress,
		},
	}
	return configMap

}

func GenerateComposeRunnerConfigMapName(crdName string) string {
	return crdName + "-dcr-cm"
}

/*
  containers:
  - name: docker-agent
    image: 6zar/docker-agent:latest
    imagePullPolicy: Always
    env:
    - name: DOCKER_CERT_PATH
      value: "/certs"
    - name: DOCKER_HOST
      value: "tcp://192.168.2.162:2376"
    - name: DOCKER_TLS_VERIFY
      value: "1"
    command: ['/home/app/docker-agent', 'agent', '--crd-api-version', 'tool.6zacode-toolbox.github.io/v1', '--crd-namespace', 'default', '--crd-instance', 'dockerhost-sample', '--crd-resource', 'dockerhosts']
    volumeMounts:
    - mountPath: "/certs"
      name: docker-certs
      readOnly: true
  serviceAccountName: docker-agent-sa
  volumes:
  - name: docker-certs
    secret:
      secretName: docker-secret
*/
