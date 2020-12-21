# Quick Start
## Using Docker
You can get the HelmValuesTransformer plugin installed kustomize container and can use it to build decapod yaml.
Mount your decapod-yaml and excute docker run command with site specific decapod-yaml directory like this:
```
$ docker run -it -v $(pwd)/examples:/decapod-yaml sktdev/decapod-kustomize:v1 kustomize build --enable_alpha_plugins /decapod-yaml/helmvalues -o /decapod-yaml/output.yml
```
And delete first line of output.yml. The line is debug message.  
## Go build
### Installation
> NOTE: go 1.14 must be installed on your environment.
```
git clone https://github.com/openinfradev/kustomize-helm-transformer.git
mkdir -p ~/.config/kustomize/plugin/openinfradev.github.com/v1/helmvaluestransformer
go build -buildmode plugin -o ~/.config/kustomize/plugin/openinfradev.github.com/v1/helmValuesTransformer/HelmValuesTransformer.so kustomize-helm-transformer/plugin/openinfradev.github.com/v1/helmvaluestransformer/HelmValuesTransformer.go
```
### Usage
```
kustomize build --enable_alpha_plugins kustomize-helm-transformer/examples/helmvalues/
```

### Run test
Run below command in the directory where HelmValuesTransformer.go exists.
```
go test
```