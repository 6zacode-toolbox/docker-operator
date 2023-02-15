package v1

type CrdDefinition struct {
	APIVersion string
	Namespace  string
	Name       string
	Resource   string
}

func (crd *DockerHost) GetCrdDefinition() *CrdDefinition {
	resource := "dockerhosts"
	return &CrdDefinition{
		APIVersion: crd.APIVersion,
		Namespace:  crd.Namespace,
		Name:       crd.Name,
		Resource:   resource,
	}
}

func (crd *DockerComposeRunner) GetCrdDefinition() *CrdDefinition {
	resource := "dockercomposerunners"
	return &CrdDefinition{
		APIVersion: crd.APIVersion,
		Namespace:  crd.Namespace,
		Name:       crd.Name,
		Resource:   resource,
	}
}

const COMPOSE_ACTION_UP = "up -d"
const COMPOSE_ACTION_DOWN = "down"
