// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"log"
	"os"

	"sigs.k8s.io/kustomize/api/resid"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/yaml"
)

// Override values in HelmReleases
type plugin struct {
	Global     map[string]string `json:"global,omitempty" yaml:"global,omitempty"`
	HelmValues []helmValues      `json:"helmValues,omitempty" yaml:"helmValues,omitempty"`
	Logger     *log.Logger
}

type helmValues struct {
	ChartName string            `json:"chartName,omitempty" yaml:"chartName,omitempty"`
	ChartRef  string            `json:"chartRef,omitempty" yaml:"chartRef,omitempty"`
	Override  map[string]string `json:"override,omitempty" yaml:"override,omitempty"`
}

// HelmRelease is a CustomResource for Helm Operator.
type HelmRelease struct {
	APIVersion string            `json:"apiVersion" yaml:"apiVersion"`
	Kind       string            `json:"kind" yaml:"kind"`
	Metadata   map[string]string `json:"metadata" yaml:"metadata"`
	Spec       HelmReleaseSpec   `json:"spec" yaml:"spec"`
}

// HelmReleaseSpec is HelmRelease's spec
type HelmReleaseSpec struct {
	Chart           HelmReleaseChart `json:"chart" yaml:"chart"`
	ReleaseName     string           `json:"releaseName,omitempty" yaml:"releaseName,omitempty"`
	TargetNamespace string           `json:"targetNamespace,omitempty" yaml:"targetNamespace,omitempty"`
	Values          interface{}      `json:"values" yaml:"values"`
}

// HelmReleaseChart is HelmRelease's Chart definition.
type HelmReleaseChart struct {
	Git  string `json:"git" yaml:"git"`
	Path string `json:"path" yaml:"path"`
	Ref  string `json:"ref" yaml:"ref"`
}

//nolint: golint
//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) Config(
	_ *resmap.PluginHelpers, c []byte) (err error) {
	p.Global = nil
	p.HelmValues = nil

	err = yaml.Unmarshal(c, p)
	if err != nil {
		return nil
	}
	if p.HelmValues == nil {
		return errors.New("helmValues is not expected to be nil")
	}
	p.Logger = log.New(os.Stdout, "[DEBUG] ", log.Lshortfile)
	return nil
}

func (p *plugin) Transform(m resmap.ResMap) (err error) {

	helmReleaseGvk := resid.Gvk{Group: "helm.fluxcd.io", Version: "v1", Kind: "HelmRelease"}
	for _, hv := range p.HelmValues {
		// replace references of HelmReleases
		id := resid.NewResId(helmReleaseGvk, hv.ChartName)
		origin, err := m.GetById(id)
		if err != nil {
			return err
		}
		if origin == nil {
			continue
		}
		if err := p.replaceChartRef(origin.Map(), hv.ChartRef); err != nil {
			return err
		}
	}
	return nil
}

func (p *plugin) replaceChartRef(origin map[string]interface{}, chartRef string) (err error) {

	// TODO: Parse helmRelease into HelmRelease

	// Temporary implement
	releaseSpec := origin["spec"].(map[string]interface{})
	chart := releaseSpec["chart"].(map[string]interface{})
	chart["ref"] = chartRef
	p.Logger.Println("chart ref changed to: ", chart["ref"]) // DEBUG

	return nil
}
