# kustomize-helm-transformer
HelmValuesTransformer is a Kustomize Plugin to transform values in `HelmRelease` CustomResource.
It helps to manage a lot of HelmRelease's value in single transformer file.  
Please take a look at the [example](https://github.com/openinfradev/kustomize-helm-transformer/tree/master/examples/helmvalues)

## Dependencies
### User
* kustomize v3.6.1 

### Contributor
* kustomize v3.6.1 
* go 1.13 

## Features
* Inline value path transform
* Chart Ref transform
<u>Source HelmRelease</u>
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
<u>Transformer Configuration</u>
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
<u>Result</u>
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

## How to install
```
git clone https://github.com/openinfradev/kustomize-helm-transformer.git
mkdir -p ~/.config/kustomize/plugin/openinfradev.github.com/v1/helmvaluestransformer
go build -buildmode plugin -o ~/.config/kustomize/plugin/openinfradev.github.com/v1/helmValuesTransformer/HelmValuesTransformer.so kustomize-helm-transformer/plugin/openinfradev.github.com/v1/helmvaluestransformer/HelmValuesTransformer.go
```
### How to use
```
kustomize build --enable_alpha_plugins kustomize-helm-transformer/examples/helmvalues/
```
