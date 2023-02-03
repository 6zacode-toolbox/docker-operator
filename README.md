# docker-operator
A simple operator to help manage dokcer based hosts from a Kubernetes cluster to leverage gitops style automations to manage workloads on small machines like raspberry-pis



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
```