package main_test

import (
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestHelmValuesTransformerChartRef(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		BuildGoPlugin("openinfradev.github.com", "v1", "HelmValuesTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: openinfradev.github.com/v1
kind: HelmValuesTransformer
metadata:
  name: site
charts:
  - chartName: glance
    chartRef: taco-k8s-v20.07
    override:
      conf.ceph.admin_keyring: abcde
      conf.ceph.enabled: true
`, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: glance
    ref: master
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: TACO_FIXME
        enabled: false
`)
	th.AssertActualEqualsExpected(rm, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: glance
    ref: taco-k8s-v20.07
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: abcde
        enabled: true
`)
}

func TestHelmValuesTransformerInlineOverride(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		BuildGoPlugin("openinfradev.github.com", "v1", "HelmValuesTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: openinfradev.github.com/v1
kind: HelmValuesTransformer
metadata:
  name: site
global:
  admin_keyring: abcdefghijklmn
charts:
  - chartName: glance
    chartRef: master
    override:
      conf.ceph.admin_keyring: $(admin_keyring)
      conf.ceph.enabled: true
`, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: glance
    ref: master
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: TACO_FIXME
        enabled: false
`)
	th.AssertActualEqualsExpected(rm, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: glance
    ref: master
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: abcdefghijklmn
        enabled: true
`)
}

func TestHelmValuesTransformerMultiCharts(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		BuildGoPlugin("openinfradev.github.com", "v1", "HelmValuesTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: openinfradev.github.com/v1
kind: HelmValuesTransformer
metadata:
  name: site
global:
  chartRef: master
  admin_keyring: abcdefghijklmn
charts:
  - chartName: glance
    chartRef: $(chartRef)
    override:
      conf.ceph.admin_keyring: $(admin_keyring)
      conf.ceph.enabled: true
  - chartName: cinder
    chartRef: $(chartRef)
    override:
      conf.ceph.admin_keyring: opqrstuvwxyz
`, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: glance
    ref: master
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: TACO_FIXME
        enabled: false
---
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: cinder
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: cinder
    ref: master
  releaseName: cinder
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: TACO_FIXME
        enabled: false
`)
	th.AssertActualEqualsExpected(rm, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: glance
    ref: master
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: abcdefghijklmn
        enabled: true
---
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: cinder
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: cinder
    ref: master
  releaseName: cinder
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: opqrstuvwxyz
        enabled: false
`)
}

func TestHelmValuesTransformerComplexValues(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		BuildGoPlugin("openinfradev.github.com", "v1", "HelmValuesTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: openinfradev.github.com/v1
kind: HelmValuesTransformer
metadata:
  name: site
global:
  glance_admin_keyring: abcdefghijklmn
  cinder_admin_keyring: opqrstuvwxyz
  docker_registry: sktdev
  image_tag: taco-0.1.0
charts:
  - chartName: glance
    chartRef: taco-k8s-v20.07
    override:
      conf.ceph.admin_keyring: $(glance_admin_keyring)
      conf.ceph.enabled: true
      images.tags.ks_user: $(docker_registry)/ubuntu-source-heat-engine-stein:$(image_tag)
  - chartName: cinder
    chartRef: feature-a
    override:
      conf.ceph.admin_keyring: $(cinder_admin_keyring)
`, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: glance
    ref: master
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: TACO_FIXME
        enabled: false
    images:
      tags:
        ks_user: TACO_FIXME
---
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: cinder
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: cinder
    ref: master
  releaseName: cinder
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: TACO_FIXME
        enabled: false
`)
	th.AssertActualEqualsExpected(rm, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: glance
    ref: taco-k8s-v20.07
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: abcdefghijklmn
        enabled: true
    images:
      tags:
        ks_user: sktdev/ubuntu-source-heat-engine-stein:taco-0.1.0
---
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: cinder
spec:
  chart:
    git: https://tde.sktelecom.com/stash/scm/openstack/openstack-helm.git
    path: cinder
    ref: feature-a
  releaseName: cinder
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: opqrstuvwxyz
        enabled: false
`)
}
