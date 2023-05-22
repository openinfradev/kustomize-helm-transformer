// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/api/filters/patchstrategicmerge"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/kustomize/kyaml/filtersutil"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/yaml"
)

// Override values in HelmReleases
type plugin struct {
	h      *resmap.PluginHelpers
	Global map[string]interface{} `json:"global,omitempty" yaml:"global,omitempty"`
	Charts []ReplacedChart        `json:"charts,omitempty" yaml:"charts,omitempty"`
	Logger *log.Logger
}

// ReplacedChart is including target information and chart values to override
type ReplacedChart struct {
	Name     string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Source   ChartSource            `json:"source,omitempty" yaml:"source,omitempty"`
	Override map[string]interface{} `json:"override,omitempty" yaml:"override,omitempty"`
}

// ChartSource defines the source of helm chart
type ChartSource struct {
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Version    string `json:"version,omitempty" yaml:"version,omitempty"`
	Type       string `json:"type,omitempty" yaml:"type,omitempty"`
}

// nolint: golint
// noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) Config(
	h *resmap.PluginHelpers, c []byte) (err error) {
	p.h = h
	p.Global = nil
	p.Charts = nil

	err = yaml.Unmarshal(c, p)
	if err != nil {
		return nil
	}
	if p.Charts == nil {
		return errors.New("helmValues is not expected to be nil")
	}
	p.Logger = log.New(os.Stdout, "[DEBUG] ", log.Lshortfile)
	return nil
}

func (p *plugin) Transform(m resmap.ResMap) (err error) {

	helmReleaseGvk := resid.Gvk{Group: "helm.fluxcd.io", Version: "v1", Kind: "HelmRelease"}
	for _, chart := range p.Charts {
		// replace references of HelmReleases
		id := resid.NewResId(helmReleaseGvk, chart.Name)
		origin, err := m.GetById(id)
		if err != nil {
			return err
		}
		if origin == nil {
			p.Logger.Println("Can't find HelmRelease name: " + chart.Name)
			continue
		}

		overrideChartResource, err := p.getChartResource(chart.Source)
		if err != nil {
			return err
		}
		err = p.applyPatch(origin, overrideChartResource)
		if err != nil {
			return err
		}

		overrideResource, err := p.getResourceFromChart(chart)
		if err != nil {
			return err
		}

		err = p.applyPatch(origin, overrideResource)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *plugin) applyPatch(resource, patch *resource.Resource) error {
	node, err := filtersutil.GetRNode(patch)
	if err != nil {
		return err
	}
	n, ns := resource.GetName(), resource.GetNamespace()
	err = filtersutil.ApplyToJSON(patchstrategicmerge.Filter{
		Patch: node,
	}, resource)
	if !resource.IsNilOrEmpty() {
		resource.SetName(n)
		resource.SetNamespace(ns)
	}
	return err
}

func (p *plugin) getChartResource(chartSource ChartSource) (r *resource.Resource, err error) {
	patchChartMap := map[string]interface{}{}
	if chartSource.Repository != "" {
		repository, err := p.replaceGlobalVar(chartSource.Repository)
		if err != nil {
			return nil, err
		}
		patchChartMap["repository"] = repository
	}

	if chartSource.Version != "" {
		version, err := p.replaceGlobalVar(chartSource.Version)
		if err != nil {
			return nil, err
		}
		patchChartMap["version"] = version
	}

	if chartSource.Name != "" {
		name, err := p.replaceGlobalVar(chartSource.Name)
		if err != nil {
			return nil, err
		}
		patchChartMap["name"] = name
	}

	if chartSource.Type != "" {
		chartType, err := p.replaceGlobalVar(chartSource.Type)
		if err != nil {
			return nil, err
		}
		patchChartMap["type"] = chartType
	}

	resource := p.h.ResmapFactory().RF().FromMap(map[string]interface{}{
		"spec": map[string]interface{}{
			"chart": patchChartMap,
		},
	})

	return resource, nil
}

func (p *plugin) getResourceFromChart(replacedChart ReplacedChart) (r *resource.Resource, err error) {
	patchMap := map[string]interface{}{}

	for inlinePath, val := range replacedChart.Override {
		newVal, err := p.replaceGlobalVar(val)
		if err != nil {
			return nil, err
		}
		paths := splitButIgnoreEscapedDot(inlinePath, "\uffff")
		p.createMapFromPaths(patchMap, paths, newVal)
	}

	resource := p.h.ResmapFactory().RF().FromMap(map[string]interface{}{
		"spec": map[string]interface{}{
			"values": patchMap,
		},
	})
	return resource, nil
}

// inlinePath is a path string using json dot notation
// i.e. "conf.ceph.admin_keyring"
func (p *plugin) createMapFromPaths(chart map[string]interface{}, paths []string, val interface{}) map[string]interface{} {
	currentPath := paths[0]
	if len(paths) == 1 {
		chart[currentPath] = val
		return chart
	}

	if chart[currentPath] == nil {
		chart[currentPath] = map[string]interface{}{}
	}
	chart[currentPath] = p.createMapFromPaths(chart[currentPath].(map[string]interface{}), paths[1:], val)
	return chart
}

func (p *plugin) replaceGlobalVar(original interface{}) (interface{}, error) {
	valueType := reflect.ValueOf(original).Kind()
	var inlineStr string
	// type checking of override value
	if valueType == reflect.Float64 || valueType == reflect.Float32 || valueType == reflect.Int {
		return original, nil
	} else if valueType == reflect.String {
		inlineStr = original.(string)
	} else if valueType == reflect.Slice || valueType == reflect.Map {
		val, _ := yaml.Marshal(original)
		inlineStr = string(val)
	}
	re := regexp.MustCompile(`\$\(([^\(\)])+\)`)
	isMatched := re.MatchString(inlineStr)

	// no global variable
	if !isMatched {
		return original, nil
	}

	for isMatched {
		findStr := re.FindString(inlineStr)
		globalVar := p.Global[findStr[2:len(findStr)-1]]

		// return error if global variable is not defined
		if globalVar == nil {
			return nil, errors.New("Can not found global variable named " + findStr)
		}

		if findStr == inlineStr {
			return globalVar, nil
		}

		inlineStr = strings.Replace(inlineStr, findStr, fmt.Sprintf("%v", globalVar), -1)
		isMatched = re.MatchString(inlineStr)
	}

	if valueType != reflect.String {
		err := yaml.Unmarshal([]byte(inlineStr), &original)
		return original, err
	}
	return inlineStr, nil
}

func splitButIgnoreEscapedDot(input, placeholder string) []string {
	temp := strings.ReplaceAll(input, "\\.", placeholder)
	parts := strings.Split(temp, ".")

	for i := range parts {
		parts[i] = strings.ReplaceAll(parts[i], placeholder, ".")
	}
	return parts
}
