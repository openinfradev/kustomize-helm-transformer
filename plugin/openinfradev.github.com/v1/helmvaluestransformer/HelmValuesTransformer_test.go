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
charts:
  - chartName: glance
    chartRef: master
    override:
      conf.ceph.admin_keyring: abcdefghijklmn
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
charts:
  - chartName: glance
    chartRef: master
    override:
      conf.ceph.admin_keyring: abcdefghijklmn
      conf.ceph.enabled: true
  - chartName: cinder
    chartRef: master
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
charts:
  - chartName: glance
    chartRef: taco-k8s-v20.07
    override:
      conf.ceph.admin_keyring: abcdefghijklmn
      conf.ceph.enabled: true
  - chartName: cinder
    chartRef: feature-a
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
    ref: taco-k8s-v20.07
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
