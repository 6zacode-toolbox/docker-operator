package controllers_test

import (
	"strings"
	"testing"

	toolsv1 "github.com/6zacode-toolbox/docker-operator/operator/api/v1"
	"github.com/6zacode-toolbox/docker-operator/operator/controllers"
	apiV1 "k8s.io/api/core/v1"
)

func TestAddHostConnectionVars_SSH(t *testing.T) {

	crd := &toolsv1.DockerHost{
		Spec: toolsv1.DockerHostSpec{
			SSHConnection: toolsv1.SSHConnection{
				SSHUser: "user",
			},
			HostIP: "127.0.0.1",
		},
	}
	originalVarList := []apiV1.EnvVar{
		{
			Name:  "PROPERTY",
			Value: "VALUE",
		},
	}
	updatedVarList := controllers.AddHostConnectionVars(crd, &originalVarList)

	if len(*updatedVarList) != 2 {
		t.Errorf("Result of wrong size: %v ", len(*updatedVarList))
	}
	for _, v := range *updatedVarList {
		if v.Name == "DOCKER_HOST" && !strings.Contains(v.Value, "ssh://") {
			t.Errorf("DOCKER_HOST produce is not starting with ssh://: %v ", v.Value)
		}
	}

}

func TestAddHostConnectionVars_TLS(t *testing.T) {

	crd := &toolsv1.DockerHost{
		Spec: toolsv1.DockerHostSpec{
			HostIP: "127.0.0.1",
		},
	}
	originalVarList := []apiV1.EnvVar{
		{
			Name:  "PROPERTY",
			Value: "VALUE",
		},
	}
	updatedVarList := controllers.AddHostConnectionVars(crd, &originalVarList)

	if len(*updatedVarList) != 4 {
		t.Errorf("Result of wrong size: %v ", len(*updatedVarList))
	}
	for _, v := range *updatedVarList {
		if v.Name == "DOCKER_HOST" && !strings.Contains(v.Value, "tcp://") {
			t.Errorf("DOCKER_HOST produce is not starting with tcp://: %v ", v.Value)
		}
	}

}
