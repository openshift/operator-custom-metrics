module github.com/openshift/operator-custom-metrics

go 1.13

require (
	github.com/coreos/prometheus-operator v0.35.1
	github.com/openshift/api v3.9.1-0.20190424152011-77b8897ec79a+incompatible
	github.com/prometheus/client_golang v1.1.0
	github.com/sirupsen/logrus v1.4.2
	k8s.io/api v0.15.7
	k8s.io/apimachinery v0.15.7
	sigs.k8s.io/controller-runtime v0.3.0
)

// Pinned to kubernetes-1.15.7
replace k8s.io/client-go => k8s.io/client-go v0.15.7
