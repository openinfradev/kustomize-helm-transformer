// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"sigs.k8s.io/kustomize/api/resmap"
)

type plugin struct {
}

//nolint: golint
//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) Config(
	_ *resmap.PluginHelpers, c []byte) (err error) {
	return nil
}

func (p *plugin) Transform(m resmap.ResMap) (err error) {
	// TODO
	return nil
}
