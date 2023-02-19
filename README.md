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
    targetRevision: 0.16.0
    helm:
      values: |-
        useExternalSecret: true
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
if set `useExternalSecret`, it will try to use external secrets to populate mandatory secrets. 


## Install Secrets

### Manually 
#### Add secret on the default namespace 

```yaml 
- apiVersion: v1
  data:
    GITHUB_TOKEN: some_github_token
  kind: Secret
  metadata:
    name: github
    namespace: default
```
#### Add secret for docker certs

```bash 
kubectl create secret generic docker-secret \
  --from-file=cert.pem=$HOME/certs/docker_host/cert.pem \
  --from-file=key.pem=$HOME/certs/docker_host/key.pem \
  --from-file=ca.pem=$HOME/certs/docker_host/ca.pem
```

or
```yaml 
apiVersion: v1
kind: Secret
metadata:
  name: docker-secret
  namespace: default  
type: Opaque
data:
  ca.pem: ....
  cert.pem: ...
  key.pem: ...
  
```

Reference on how to generate this certs: [here](https://medium.com/p/c95e78817fa6)
> then use this certs to populate your docker host and your cluster secrets

### Using Secrets Operator

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: docker-ssh
  annotations:
    argocd.argoproj.io/sync-wave: "0"
spec:
  target:
    name: docker-ssh
  secretStoreRef:
    kind: ClusterSecretStore
    name: vault-secrets-backend
  refreshInterval: 10s
  dataFrom:
    - extract:      
        key: /docker-ssh
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: docker-secret
  annotations:
    argocd.argoproj.io/sync-wave: "0"
spec:
  target:
    name: docker-secret
  secretStoreRef:
    kind: ClusterSecretStore
    name: vault-secrets-backend
  refreshInterval: 10s
  dataFrom:
    - extract:      
        key: /docker-certs
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: github-secret
  annotations:
    argocd.argoproj.io/sync-wave: "0"
spec:
  target:
    name: github
  secretStoreRef:
    kind: ClusterSecretStore
    name: vault-secrets-backend
  refreshInterval: 10s
  data:
    - secretKey: GITHUB_TOKEN
      remoteRef:
        key: /ci-secrets        
        property: PERSONAL_ACCESS_TOKEN

```

Note that `/docker-certs` need to be create on vault with, if it doesn't exist already: 

```json
{
  "ca.pem": "-----BEGIN CERTIFICATE-----\n....\n-----END CERTIFICATE-----",
  "cert.pem": "-----BEGIN CERTIFICATE-----\n....\n-----END CERTIFICATE-----",
  "key.pem": "-----BEGIN RSA PRIVATE KEY-----\n....\n-----END RSA PRIVATE KEY-----"
}
```

## Sample of `config` file for docker

```bash
Host 192.168.2.*
   StrictHostKeyChecking no
   UserKnownHostsFile=/dev/null
```