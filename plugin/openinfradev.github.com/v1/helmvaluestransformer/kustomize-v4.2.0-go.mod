module sigs.k8s.io/kustomize/plugin/openinfradev.github.com/v1/helmvaluestransformer

go 1.16

require (
    github.com/pkg/errors v0.9.1
	sigs.k8s.io/kustomize/api v0.8.9
	sigs.k8s.io/kustomize/kyaml v0.11.0
	sigs.k8s.io/yaml v1.2.0
)

replace sigs.k8s.io/kustomize/api => ../../../../api
