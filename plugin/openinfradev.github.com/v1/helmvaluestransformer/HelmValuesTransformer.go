// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"sigs.k8s.io/kustomize/api/resid"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/yaml"
)

// Override values in HelmReleases
type plugin struct {
	h      *resmap.PluginHelpers
	Global map[string]string `json:"global,omitempty" yaml:"global,omitempty"`
	Charts []ReplacedChart   `json:"charts,omitempty" yaml:"charts,omitempty"`
	Logger *log.Logger
}

// ReplacedChart is including target information and chart values to override
type ReplacedChart struct {
	ChartName string                 `json:"chartName,omitempty" yaml:"chartName,omitempty"`
	ChartRef  string                 `json:"chartRef,omitempty" yaml:"chartRef,omitempty"`
	Override  map[string]interface{} `json:"override,omitempty" yaml:"override,omitempty"`
}

//nolint: golint
//noinspection GoUnusedGlobalVariable
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
		id := resid.NewResId(helmReleaseGvk, chart.ChartName)
		origin, err := m.GetById(id)
		if err != nil {
			return err
		}
		if origin == nil {
			p.Logger.Println("Can't find HelmRelease name: " + chart.ChartName)
			continue
		}
		if err := p.replaceChartRef(origin.Map(), chart.ChartRef); err != nil {
			return err
		}
		overrideResource := p.getResourceFromChart(chart)
		if err := origin.Patch(overrideResource.Copy()); err != nil {
			p.Logger.Println("patch error: " + err.Error())
			return err
		}
	}
	return nil
}

func (p *plugin) replaceChartRef(origin map[string]interface{}, chartRef string) (err error) {
	releaseSpec := origin["spec"].(map[string]interface{})
	chart := releaseSpec["chart"].(map[string]interface{})
	chart["ref"] = chartRef
	return nil
}

func (p *plugin) getResourceFromChart(replacedChart ReplacedChart) (r *resource.Resource) {
	patchMap := map[string]interface{}{}

	for inlinePath, val := range replacedChart.Override {
		p.createMapFromPaths(patchMap, strings.Split(inlinePath, "."), val)
	}

	resource := p.h.ResmapFactory().RF().FromMap(map[string]interface{}{
		"spec": map[string]interface{}{
			"values": patchMap,
		},
	})
	return resource
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
