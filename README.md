# kustomize-helm-transformer
HelmValuesTransformer is a Kustomize Plugin to transform values in `HelmRelease` CustomResource.
It helps to manage a lot of HelmRelease's value in single transformer file.  
Please take a look at the [example](https://github.com/openinfradev/kustomize-helm-transformer/tree/master/examples/helmvalues)

## Documents
* [Quick Start](docs/quickstart.md)
* [Contribution](docs/contribution.md)
* [CI Pipeline](docs/ci.md)


## Support 
* kustomize v3.8.7
* go 1.14

## Features
1. Inline value path transform
2.  Chart Ref transform

## Example
### Source HelmRelease
```
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://github.com/openstack/openstack-helm.git
    path: glance
    ref: master
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: admin_keyring
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
  - chartName: glance
    chartRef: taco-k8s-v20.07
    override:
      conf.ceph.admin_keyring: abcde
      conf.ceph.enabled: true
```
### Result
```
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://github.com/openstack/openstack-helm.git
    path: glance
    ref: taco-k8s-v20.07
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: admin_keyring
        enabled: true
```
