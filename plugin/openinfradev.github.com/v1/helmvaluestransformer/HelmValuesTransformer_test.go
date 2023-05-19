package main_test

import (
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestHelmValuesTransformerChartSource(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		BuildGoPlugin("openinfradev.github.com", "v1", "HelmValuesTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: openinfradev.github.com/v1
kind: HelmValuesTransformer
metadata:
  name: site
charts:
  - name: glance
    source:
      repository: http://repository:8879
      version: 1.0.0
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
    name: glance
    repository: TO_BE_FIXED
    version: 0.1.0
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
    name: glance
    repository: http://repository:8879
    version: 1.0.0
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
  - name: glance
    source:
      repository: http://repository:8879
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
    name: glance
    repository: TO_BE_FIXED
    version: 0.1.0
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
    name: glance
    repository: http://repository:8879
    version: 0.1.0
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
  repository: http://repository:8879
  admin_keyring: abcdefghijklmn
charts:
  - name: glance
    source:
      repository: $(repository)
    override:
      conf.ceph.admin_keyring: $(admin_keyring)
      conf.ceph.enabled: true
  - name: cinder
    source: 
      repository: $(repository)
    override:
      conf.ceph.admin_keyring: opqrstuvwxyz
`, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    name: glance
    repository: TO_BE_FIXED
    version: 0.1.0
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
    name: cinder
    repository: TO_BE_FIXED
    version: 0.1.0
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
    name: glance
    repository: http://repository:8879
    version: 0.1.0
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
    name: cinder
    repository: http://repository:8879
    version: 0.1.0
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
  storageClassName: ceph
  grafanaDomain: grafana.example.com
  keycloakDomain: keycloak.example.com
  realms: openinfradev
charts:
  - name: glance
    source:
      repository: http://repository-a:8879
      version: 1.0.1
    override:
      grafana\.ini: 
        server:
          domain: $(grafanaDomain)
          root_url: $(grafanaDomain)/grafana
        auth.generic_oauth:
          enabled: true
          auth_url: https://$(keycloakDomain)/auth/realms/$(realms)/protocol/openid-connect/auth
      conf.ceph.admin_keyring: $(glance_admin_keyring)
      conf.ceph.enabled: true
      images.tags.ks_user: $(docker_registry)/ubuntu-source-heat-engine-stein:$(image_tag)
      volumeClaimTemplates:
      - metadata:
          name: glance-data
        spec:
          storageClassName: $(storageClassName)
  - name: cinder
    source:
      repository: http://repository-b:8879
      version: 2.0.2
    override:
      conf.ceph.admin_keyring: $(cinder_admin_keyring)
`, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    name: glance
    repository: TO_BE_FIXED
    version: 0.1.0
  releaseName: glance
  targetNamespace: openstack
  values:
    grafana.ini:
    conf:
      ceph:
        admin_keyring: TO_BE_FIXED
        enabled: false
    images:
      tags:
        ks_user: TO_BE_FIXED
    volumeClaimTemplates:
    - TO_BE_FIXED
---
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: cinder
spec:
  chart:
    name: cinder
    repository: TO_BE_FIXED
    version: 0.1.0
  releaseName: cinder
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: TO_BE_FIXED
        enabled: false
`)
	th.AssertActualEqualsExpected(rm, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: glance
spec:
  chart:
    name: glance
    repository: http://repository-a:8879
    version: 1.0.1
  releaseName: glance
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: abcdefghijklmn
        enabled: true
    grafana.ini:
      auth.generic_oauth:
        auth_url: https://keycloak.example.com/auth/realms/openinfradev/protocol/openid-connect/auth
        enabled: true
      server:
        domain: grafana.example.com
        root_url: grafana.example.com/grafana
    images:
      tags:
        ks_user: sktdev/ubuntu-source-heat-engine-stein:taco-0.1.0
    volumeClaimTemplates:
    - metadata:
        name: glance-data
      spec:
        storageClassName: ceph
---
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: cinder
spec:
  chart:
    name: cinder
    repository: http://repository-b:8879
    version: 2.0.2
  releaseName: cinder
  targetNamespace: openstack
  values:
    conf:
      ceph:
        admin_keyring: opqrstuvwxyz
        enabled: false
`)
}

func TestGitChartSource(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		BuildGoPlugin("openinfradev.github.com", "v1", "HelmValuesTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: openinfradev.github.com/v1
kind: HelmValuesTransformer
metadata:
  name: site
charts:
  - name: kube-prometheus-stack
    source:
      repository: git@github.com:helm/charts
      version: master
      name: charts/stable/prometheus-operator
      type: git
`, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: kube-prometheus-stack
spec:
  chart:
    name: TO_BE_FIXED
    repository: TO_BE_FIXED
    version: TO_BE_FIXED
    type: TO_BE_FIXED
  releaseName: kube-prometheus-stack
  targetNamespace: lma
  values:
    conf:
      ceph:
        admin_keyring: abcde
        enabled: true
`)
	th.AssertActualEqualsExpected(rm, `
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: kube-prometheus-stack
spec:
  chart:
    name: charts/stable/prometheus-operator
    repository: git@github.com:helm/charts
    type: git
    version: master
  releaseName: kube-prometheus-stack
  targetNamespace: lma
  values:
    conf:
      ceph:
        admin_keyring: abcde
        enabled: true
`)
}
