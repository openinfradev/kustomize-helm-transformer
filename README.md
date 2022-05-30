# kustomize-helm-transformer
HelmValuesTransformer is a Kustomize Plugin to transform values in `HelmRelease` CustomResource.
It helps to manage a lot of HelmRelease's value in single transformer file.  
Please take a look at the [example](https://github.com/openinfradev/kustomize-helm-transformer/tree/master/examples/helmvalues)

## Documents
* [Quick Start](docs/quickstart.md)
* [Contribution](docs/contribution.md)
* [CI Pipeline](docs/ci.md)


## Support 
* kustomize v4.2.0
* go 1.17.6

## Features
1. Replaced values of HelmRelease CustomResource using inline path
2. Replaced Chart Source of HelmRelease CustomResource

## Example
### Source HelmRelease
```
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: prometheus
spec:
  chart:
    repository: https://prometheus-community.github.io/helm-charts
    name: kube-prometheus-stack
    version: 14.5.0
    type: helmrepo
  releaseName: prometheus
  targetNamespace: lma
  values:
    conf:
      ceph:
        admin_keyring: TO_BE_FIXED
        enabled: false
```
### Transformer Configuration
```
apiVersion: openinfradev.github.com/v1
kind: HelmValuesTransformer
metadata:
  name: site
global:
  docker_registry: registry.cicd.stg.taco
charts:
  - name: prometheus
    source: 
      repository: git@github.com:helm/charts
      version: master
      name: charts/stable/prometheus-operator
      type: git
    override:
      conf.ceph.admin_keyring: abcde
      conf.ceph.enabled: true
```
### Result
```
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: prometheus
spec:
  chart:
    repository: git@github.com:helm/charts
    version: master
    name: charts/stable/prometheus-operator
    type: git
  releaseName: prometheus
  targetNamespace: lma
  values:
    conf:
      ceph:
        admin_keyring: abcde
        enabled: true
```
