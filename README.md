# docker-operator
A simple operator to help manage dokcer based hosts from a Kubernetes cluster to leverage gitops style automations to manage workloads on small machines like raspberry-pis


# How to Use this operator

## Use helm to deploy the controller

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: 6zacode-docker-operator
  namespace: argocd
spec: 
  project: default
  source:
    repoURL: 'https://6zacode-toolbox.github.io/docker-operator'
    targetRevision: 0.10.0
    chart: docker-operator
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: operator-system
  syncPolicy:
      automated:
        prune: true
        selfHeal: true
      syncOptions:
        - CreateNamespace=true
      retry:
        limit: 5
        backoff:
          duration: 5s
          maxDuration: 5m0s
          factor: 2
```

## Add secret on the default namespace 

```yaml 
- apiVersion: v1
  data:
    GITHUB_TOKEN: some_github_token
  kind: Secret
  metadata:
    name: github
    namespace: default
```


## How this repo was created

```bash 
docker run --rm -it \
    --name kubebuilder \
    --platform linux/amd64 \
    -w /go/src \
    -v operator-sdk:/go/pkg \
    -v $(pwd):/go/src \
    -v /var/run/docker.sock:/var/run/docker.sock  \
    --privileged \
    6zar/kubebuilder

mkdir operator

cd operator

kubebuilder init --domain 6zacode-toolbox.github.io --license apache2 --owner "6zacode-toolbox" --repo github.com/6zacode-toolbox/docker-operator/operator

kubebuilder create api --group tool --version v1 --kind DockerHost --controller --resource

kubebuilder create api --group tool --version v1 --kind DockerComposeRunner --controller --resource

make manifests
# When doing complex objects start by parts, as serilziation may be tricky

make install

make run
```

## Error with bad object types
```bash 
Error: not all generators ran successfully
run `controller-gen rbac:roleName=manager-role crd webhook paths=./... output:crd:artifacts:config=config/crd/bases -w` to see all available markers, or `controller-gen rbac:roleName=manager-role crd webhook paths=./... output:crd:artifacts:config=config/crd/bases -h` for usage
make: *** [Makefile:43: manifests] Error 1
```

```bash 
https://6zacode-toolbox.github.io/docker-operator/index.yaml
```