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

## Building notes

```bash
go build -ldflags="-X 'controllers.NamespaceJobs=operator-system'"
```
